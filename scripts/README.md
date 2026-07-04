# scripts

Repo automation that runs in CI and locally.

## `check-api-contract.mjs` — API contract drift gate

Stands up the **api-contract-guardian** as an automated CI check
(`docs/12_AGENTS_AND_SCALING.md` §7, item 2) so contract drift is caught on
every PR instead of only when the agent is run on demand.

### What it does

It compares two sets of `(METHOD, PATH)` pairs and fails if they disagree:

1. **Code** — every route registered in `backend/internal/*/api/routes.go`
   (Go 1.22+ `ServeMux` patterns: `mux.HandleFunc`, `mux.Handle`, and the local
   `reg(...)` / `recruiter(...)` / `admin(...)` helper closures).
2. **Docs** — every route documented in the markdown tables of
   `docs/06_API_CONTRACTS.md`.

It reports, and exits non-zero on, two kinds of drift:

- **Documented but missing in code** — a stale or aspirational doc entry.
- **In code but not documented** — an endpoint someone added without updating
  the contract.

Transport/infra routes (SSE `*/stream`, websocket `/ws`, `/health`, the internal
`/users/search`) are listed as informational only and never fail the build; see
`CODE_ONLY_IGNORE` in the script. Path-param notation (`:id` vs `{id}`), the
`/api/v1` prefix, query strings, and compressed CRUD rows (`POST/PUT/DELETE …
[/:id]`) are normalized — the assumptions are documented in the script header.

### Run it locally

```sh
node scripts/check-api-contract.mjs
```

Exit codes: `0` = in sync, `1` = drift detected (report printed), `2` = hard
error (e.g. a required file is missing). No dependencies — plain Node ≥ 18.

### CI

Wired into `.github/workflows/contract-check.yml`, which runs on every
`pull_request` and on pushes to `master`/`main` via `actions/setup-node`. It is
a **real merge gate**: a failing contract check blocks the PR.

### Fixing drift

The check is a lint, not a fixer. When it fails, resolve the drift at the
source:

- Endpoint added/changed in code → update the matching table in
  `docs/06_API_CONTRACTS.md`.
- Doc lists an endpoint that no longer exists → remove/correct the doc row, or
  restore the route.
- Genuinely-internal infra route → add its path pattern to `CODE_ONLY_IGNORE`
  in `scripts/check-api-contract.mjs` (with a comment saying why).

### Disabling / rollback (escape hatch)

> **Do not disable CI gates to ship faster.** This is the documented escape
> hatch for a true emergency (e.g. the checker itself is broken and blocking an
> unrelated hotfix), not a routine workaround. Prefer fixing the drift.

If you must bypass it temporarily, in order of least- to most-disruptive:

1. **Skip a single hotfix PR** — temporarily mark the `API contract drift` job
   non-required in the branch-protection settings (GitHub → Settings → Branches),
   merge the hotfix, then immediately re-mark it required. Leaves the gate intact
   for everyone else.
2. **Quiet a known-internal route** — if the failure is a legitimately
   undocumented infra route, add it to `CODE_ONLY_IGNORE` rather than disabling
   the whole job.
3. **Disable the workflow** — last resort. Either delete/rename
   `.github/workflows/contract-check.yml`, or wrap the run step with an `if:`
   guard. Open a ticket to re-enable it in the same sprint.

**Re-enable** by reverting the change (restore the workflow file and/or re-mark
the job required) and confirming `node scripts/check-api-contract.mjs` passes.
