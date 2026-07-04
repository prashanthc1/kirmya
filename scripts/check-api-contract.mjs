#!/usr/bin/env node
// check-api-contract.mjs
//
// CI guard that mirrors a slice of what the `api-contract-guardian` agent does:
// it keeps `docs/06_API_CONTRACTS.md` in sync with the Go HTTP route
// registrations, so contract drift is caught on every PR instead of only when
// the agent is run on demand.
//
// It compares two sets of (METHOD, PATH) pairs:
//   1. CODE  - every route registered in `backend/internal/*/api/routes.go`.
//   2. DOC   - every route documented in the tables of `docs/06_API_CONTRACTS.md`.
//
// It reports:
//   - routes in CODE but missing from the DOC (undocumented endpoints), and
//   - routes in the DOC with no matching CODE route (stale / aspirational docs).
//
// Exit code: 0 when the two sets agree, 1 when drift is found, 2 on a hard
// error (e.g. files missing). The report it prints is meant to be readable in a
// GitHub Actions log.
//
// -----------------------------------------------------------------------------
// Pragmatic assumptions (documented on purpose - this is a lint, not a parser):
//
//  * Route strings in Go use the Go 1.22+ ServeMux pattern syntax,
//    `"METHOD /api/v1/path/{param}"`, registered via `mux.HandleFunc(...)`,
//    `mux.Handle(...)`, or a thin local helper closure such as
//    `reg("GET /api/v1/...", h.Foo)` / `recruiter(...)` / `admin(...)`. We do
//    not evaluate Go; we just scan for string literals whose first token is an
//    HTTP method. That is robust to which helper wraps them.
//
//  * Path params are normalized so the two notations compare equal:
//    Go `{id}` / `{provider}` / `{slug}`  <->  doc `:id` / `:provider` / `:slug`.
//    Every param is collapsed to a single placeholder token `{}` so the *name*
//    of the param never matters, only its position.
//
//  * The `/api/v1` prefix is stripped from both sides (the doc tables mostly
//    omit it and state it once in the header; the code always includes it).
//
//  * Query strings in the doc (`/search?q=&type=`) are dropped - they are not
//    part of the route's method+path identity.
//
//  * The doc compresses some rows. We expand them:
//      - "POST/PUT/DELETE" in a Method cell  -> separate methods.
//      - An optional trailing segment written as "[/:id]" -> the collection
//        verb (POST) targets the base path and the item verbs (PUT/PATCH/DELETE)
//        target base + segment (covers e.g. `/profiles/me/experiences[/:id]`).
//      - "WS"/"SSE" are treated as pseudo-methods; the code uses a plain GET for
//        SSE/WS upgrades, so those rows are reported as informational only and
//        never fail the build (see WS_METHODS below).
//
//  * Routes that exist in code but are intentionally outside the documented
//    contract surface (health checks, SSE `/stream` channels, websocket `/ws`,
//    internal `/users/search`) are listed as info, never a hard failure (see
//    CODE_ONLY_IGNORE). Everything else that is in code but undocumented - or
//    documented but not in code - fails the build, because that is real drift.
// -----------------------------------------------------------------------------

import { readFileSync, readdirSync, existsSync } from "node:fs";
import { join, dirname } from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = dirname(fileURLToPath(import.meta.url));
const REPO_ROOT = join(__dirname, "..");
const BACKEND_GLOB_DIR = join(REPO_ROOT, "backend", "internal");
const CONTRACT_DOC = join(REPO_ROOT, "docs", "06_API_CONTRACTS.md");

const HTTP_METHODS = ["GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"];
// Pseudo-methods documented in the contract that do not map 1:1 to a Go method
// (the SSE/WS upgrade handlers register as GET). Never fail the build on these.
const WS_METHODS = ["WS", "SSE"];

// Code-only path suffixes that are deliberately undocumented infra/transport
// channels (server-sent events, websockets, health, internal search). Reported
// as info, never a hard failure.
const CODE_ONLY_IGNORE = [/\/stream$/, /\/ws$/, /^\/health$/, /\/users\/search$/];

function fail(msg) {
  console.error(`check-api-contract: ${msg}`);
  process.exit(2);
}

// --- path normalization ------------------------------------------------------

function normalizePath(rawPath) {
  let p = rawPath.trim();
  // Drop query string.
  p = p.split("?")[0];
  // Strip the common /api/v1 prefix so doc (often omits it) and code align.
  p = p.replace(/^\/api\/v1/, "");
  // Collapse both `{param}` and `:param` to a single positional placeholder.
  p = p.replace(/\{[^/}]+\}/g, "{}").replace(/:[^/]+/g, "{}");
  // Strip a trailing slash (but keep root "/").
  if (p.length > 1) p = p.replace(/\/+$/, "");
  if (p === "") p = "/";
  return p;
}

function key(method, path) {
  return `${method.toUpperCase()} ${normalizePath(path)}`;
}

// --- CODE side: parse Go route registrations ---------------------------------

function findRouteFiles() {
  const files = [];
  if (!existsSync(BACKEND_GLOB_DIR)) {
    fail(`backend module dir not found: ${BACKEND_GLOB_DIR}`);
  }
  for (const module of readdirSync(BACKEND_GLOB_DIR, { withFileTypes: true })) {
    if (!module.isDirectory()) continue;
    const candidate = join(BACKEND_GLOB_DIR, module.name, "api", "routes.go");
    if (existsSync(candidate)) files.push(candidate);
  }
  return files.sort();
}

// Match any string literal whose content begins with an HTTP method token and a
// space, e.g.  "GET /api/v1/jobs/{id}". This covers mux.HandleFunc / mux.Handle
// and every helper-closure form, because they all pass the pattern as the first
// string argument.
const ROUTE_LITERAL = new RegExp(
  `"\\s*(${HTTP_METHODS.join("|")})\\s+(/[^"]*)"`,
  "g"
);

function parseCodeRoutes() {
  const routes = new Map(); // key -> { method, path, module, file }
  for (const file of findRouteFiles()) {
    const module = file.split(/[\\/]/).slice(-3, -2)[0]; // .../<module>/api/routes.go
    const src = readFileSync(file, "utf8");
    let m;
    ROUTE_LITERAL.lastIndex = 0;
    while ((m = ROUTE_LITERAL.exec(src)) !== null) {
      const method = m[1].toUpperCase();
      const path = m[2];
      routes.set(key(method, path), {
        method,
        path: normalizePath(path),
        module,
        file,
      });
    }
  }
  return routes;
}

// --- DOC side: parse markdown tables -----------------------------------------

// Expand a doc method cell that may list several methods, e.g. "POST/PUT/DELETE".
function expandMethods(cell) {
  return cell
    .split("/")
    .map((s) => s.trim().toUpperCase())
    .filter(Boolean);
}

// Given a doc method cell and path cell, produce the concrete (method, path)
// pairs the row represents.
//
// Two shapes occur in the contract doc:
//   (a) A single method + a plain path                -> one pair.
//   (b) A compressed CRUD row: several methods +
//       a path with an optional trailing "[/:id]"     -> one pair per method,
//       where the *collection* verb (POST) targets the base path and the
//       *item* verbs (PUT/PATCH/DELETE) target the base + segment. This mirrors
//       the project's REST convention (POST /x ; PUT|DELETE /x/{id}) and avoids
//       inventing cartesian-product combinations the doc never intends (e.g.
//       there is no DELETE on the bare collection path).
//
// Multi-method rows on a path *without* an optional segment just fan out to the
// same path for each method.
function expandRow(methodCell, pathCell) {
  const methods = expandMethods(methodCell).filter(
    (t) => HTTP_METHODS.includes(t) || WS_METHODS.includes(t)
  );
  const p = pathCell.trim().replace(/`/g, "").trim();

  const optional = p.match(/^(.*)\[(\/[^\]]+)\]$/);
  const pairs = [];

  if (optional && methods.length > 1) {
    // Compressed CRUD row.
    const base = optional[1];
    const item = base + optional[2];
    for (const method of methods) {
      const path = method === "POST" ? base : item;
      pairs.push({ method, path });
    }
    return pairs;
  }

  // Either a single method, or a multi-method row with no optional segment.
  // If an optional segment is present with a single method, emit both variants.
  const paths =
    optional && methods.length <= 1
      ? [optional[1], optional[1] + optional[2]]
      : [p];
  for (const method of methods) {
    for (const path of paths) pairs.push({ method, path });
  }
  return pairs;
}

function parseDocRoutes() {
  if (!existsSync(CONTRACT_DOC)) fail(`contract doc not found: ${CONTRACT_DOC}`);
  const src = readFileSync(CONTRACT_DOC, "utf8");
  const routes = new Map(); // key -> { method, path, raw }
  const wsRoutes = [];

  for (const line of src.split(/\r?\n/)) {
    const trimmed = line.trim();
    if (!trimmed.startsWith("|")) continue;
    const cells = trimmed
      .split("|")
      .slice(1, -1)
      .map((c) => c.trim());
    if (cells.length < 2) continue;

    const methodCell = cells[0];
    const pathCell = cells[1];

    // Skip header / separator rows.
    if (/^-+$/.test(methodCell.replace(/\s/g, ""))) continue;
    if (/^method$/i.test(methodCell)) continue;

    // A path cell must start with a `/` (after optional backtick) to be a route.
    const looksLikePath = /^`?\//.test(pathCell);
    if (!looksLikePath) continue;

    const pairs = expandRow(methodCell, pathCell);
    if (pairs.length === 0) continue;

    for (const { method, path } of pairs) {
      if (WS_METHODS.includes(method)) {
        wsRoutes.push({ method, path: normalizePath(path) });
        continue;
      }
      routes.set(key(method, path), {
        method,
        path: normalizePath(path),
        raw: `${method} ${pathCell.replace(/`/g, "")}`,
      });
    }
  }
  return { routes, wsRoutes };
}

// --- compare -----------------------------------------------------------------

function isIgnoredCodeOnly(path) {
  return CODE_ONLY_IGNORE.some((re) => re.test(path));
}

function main() {
  const code = parseCodeRoutes();
  const { routes: doc, wsRoutes } = parseDocRoutes();

  const codeKeys = new Set(code.keys());
  const docKeys = new Set(doc.keys());

  // Documented but no matching code route -> hard drift (stale/aspirational doc).
  const docOnly = [...docKeys].filter((k) => !codeKeys.has(k));
  // In code but not documented.
  const codeOnlyAll = [...codeKeys].filter((k) => !docKeys.has(k));
  const codeOnlyIgnored = codeOnlyAll.filter((k) =>
    isIgnoredCodeOnly(code.get(k).path)
  );
  const codeOnly = codeOnlyAll.filter((k) => !isIgnoredCodeOnly(code.get(k).path));

  console.log("API contract drift check");
  console.log("========================");
  console.log(`code routes:       ${codeKeys.size}  (${findRouteFiles().length} routes.go files)`);
  console.log(`documented routes: ${docKeys.size}  (docs/06_API_CONTRACTS.md)`);
  if (wsRoutes.length) {
    console.log(
      `ws/sse rows (info, never gated): ${wsRoutes
        .map((r) => `${r.method} ${r.path}`)
        .join(", ")}`
    );
  }
  console.log("");

  let drift = false;

  if (docOnly.length) {
    drift = true;
    console.log(`DOCUMENTED BUT MISSING IN CODE (${docOnly.length}):`);
    for (const k of docOnly.sort()) {
      const r = doc.get(k);
      console.log(`  - ${r.method} ${r.path}`);
    }
    console.log("");
  }

  if (codeOnly.length) {
    // Code-only routes are a real drift signal (someone added an endpoint and
    // forgot the doc). Gate on them too - that is the whole point of the check.
    drift = true;
    console.log(`IN CODE BUT NOT DOCUMENTED (${codeOnly.length}):`);
    for (const k of codeOnly.sort()) {
      const r = code.get(k);
      console.log(`  - ${r.method} ${r.path}   [${r.module}]`);
    }
    console.log("");
  }

  if (codeOnlyIgnored.length) {
    console.log(
      `IN CODE, INTENTIONALLY UNDOCUMENTED (info - transport/infra, ${codeOnlyIgnored.length}):`
    );
    for (const k of codeOnlyIgnored.sort()) {
      const r = code.get(k);
      console.log(`  - ${r.method} ${r.path}   [${r.module}]`);
    }
    console.log("");
  }

  if (!drift) {
    console.log("OK - documented contract and code routes are in sync.");
    process.exit(0);
  }

  console.log(
    "DRIFT DETECTED. Update docs/06_API_CONTRACTS.md (or the routes) so they agree."
  );
  console.log(
    "If a code route is intentionally undocumented infra, add it to CODE_ONLY_IGNORE in this script."
  );
  process.exit(1);
}

main();
