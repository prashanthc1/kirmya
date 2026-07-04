# 13 · Migration & Rollback Strategy

This document defines how Kirmya evolves its PostgreSQL schema and, critically,
how we recover when a deploy that carries a migration needs to be rolled back.
It is the authoritative reference for the "stated rollback path" required by the
Definition of Done (`CLAUDE.md` §7) whenever a slice ships new infrastructure.

## 1. How migrations work today

Migrations are plain, forward-only SQL files in `backend/migrations/`, named
`NNN_description.sql` and applied in lexical order by the dependency-free runner
in `internal/platform/migrate`. The runner is idempotent: it records every
applied filename in a `schema_migrations` table and skips anything already
recorded, so it is safe to run on every boot. Each file executes inside a single
transaction, so a *failed* migration rolls itself back cleanly and leaves the
schema untouched — the binary then refuses to start (`main.go` calls
`log.Fatalf` on migration error), which is the desired fail-closed behavior.

There are deliberately **no down-migration files** and no automated reverse
command. This is a conscious choice, not an omission: down-migrations are
notoriously under-tested, frequently destroy data, and give false confidence.
Kirmya's rollback strategy is instead built on two pillars — **backward-
compatible schema changes (expand/contract)** and **database backups / point-in-
time recovery** — described below.

## 2. The golden rule: migrations must be backward compatible

The application is a long-running modular monolith that we deploy with a brief
overlap between old and new instances. A migration must therefore never break
the *currently running* version of the code. Concretely, a single migration may
freely do additive, non-breaking things, and must avoid destructive ones in the
same release as the code that stops using a column.

Safe in one step: adding a nullable column; adding a new table; adding an index
(use `CREATE INDEX CONCURRENTLY` for large tables, outside a transaction);
adding a column with a non-volatile default; widening a type; adding a new
enum value.

Unsafe in one step (must be split across releases): dropping or renaming a
column or table; adding a `NOT NULL` constraint to an existing column without a
default; narrowing a type; changing a column's meaning. Each of these will break
either the old code (during the deploy overlap) or the new code (if the deploy
is rolled back), which is exactly the situation this strategy exists to prevent.

## 3. Expand / contract — the rollback-safe change pattern

Any breaking change is performed as a sequence of releases, each individually
reversible by simply redeploying the previous binary (no schema change needed to
roll back):

**Expand.** Add the new shape alongside the old. Example for a rename
`headline → title`: migration adds a nullable `title` column; the new code
writes to both `title` and `headline` and reads `title` with a fallback to
`headline`. Ship and bake.

**Migrate data.** A separate migration (or a backfill job) copies
`headline` into `title`. Backfills over large tables run in batches, not one
giant `UPDATE`.

**Contract.** Once the new code has been stable in production long enough that
we will not roll back past it, a later release drops the now-unused `headline`
column. The drop is the *only* irreversible step, and it happens deliberately,
days after the behavior change, when the risk of needing it back is gone.

Because every intermediate release is backward compatible, **rolling back the
application is always just "redeploy the previous image"** — the database stays
ahead and keeps serving both versions.

## 4. Rolling back a release

When a deploy goes wrong, choose the lowest-impact option that resolves it:

1. **Code-only rollback (default, no DB action).** If the migration in the bad
   release followed §2 (backward compatible), redeploy the previous backend
   image. The new columns/tables simply sit unused; nothing else is required.
   This is the path the expand/contract discipline is designed to guarantee.

2. **Forward fix.** Often faster and safer than any reverse: write a new
   `NNN+1` migration that corrects the problem (e.g. drops a bad index, relaxes
   a constraint) and deploy it. Forward-only history stays linear and auditable.

3. **Manual reverse migration.** If a change must be undone at the schema level
   and it is safe to do so (no data loss), author a new forward migration that
   reverses it (e.g. `0NN_revert_xyz.sql` that drops the column added in
   `0MM`). Never delete or edit an already-applied migration file or hand-edit
   `schema_migrations`; the reversal is itself a new, recorded migration.

4. **Restore from backup / PITR (last resort, data-loss change).** If a
   destructive migration ran and corrupted or dropped data, the only true
   recovery is the backup. Restore the most recent backup or use point-in-time
   recovery to the moment just before the migration. This loses writes that
   happened after that point, so it is reserved for genuine emergencies — and is
   the reason destructive steps are isolated to their own release in §3.

## 5. Backups & point-in-time recovery (the real safety net)

Backups, not down-migrations, are what make destructive changes survivable.
Operationally we require, in the production environment:

- A nightly logical dump (`pg_dump`) retained per the environment's policy, plus
  WAL archiving enabled so point-in-time recovery to an arbitrary timestamp is
  possible. Managed Postgres (e.g. the Railway/Hostinger production targets in
  `docs/DEPLOY_*.md`) provides automated daily backups + PITR; confirm it is on.
- A **mandatory backup immediately before** any release containing a destructive
  (`DROP`/`ALTER ... DROP`/type-narrowing) migration. Capture the backup
  identifier/timestamp in the deploy record so the restore target is unambiguous.
- Periodic **restore drills** — a backup that has never been restored is a
  hypothesis, not a safety net. Verify a restore into a scratch database at least
  once per release train.

## 6. Pre-deploy checklist for any migration-bearing release

Before merging/shipping a slice that adds a migration:

- The migration is additive, or, if it is part of a breaking change, it is the
  expand step of an expand/contract sequence (§3) — not a one-shot
  drop/rename/`NOT NULL`.
- The migration file is the next `NNN` in sequence, runs in a transaction
  (the runner wraps it), and was tested by `make migrate` against a fresh DB and
  by re-running it (idempotency) against an already-migrated DB.
- Repository/integration tests pass against the new schema
  (`go test -tags=integration ./...`), since they apply all migrations.
- For destructive releases only: a fresh backup is taken and its identifier
  recorded, and the rollback decision (code-only vs. restore) is written into the
  deploy notes.
- The rollback path for *this* release is stated in the PR description, per the
  Definition of Done.

## 7. Quick reference

| Situation | Action |
|---|---|
| Bad release, migration was additive | Redeploy previous image (code-only) |
| Logic bug in a migration | Ship a new forward-fix migration |
| Need to undo a non-destructive schema change | New reverse forward-migration |
| Destructive migration lost data | Restore from backup / PITR |
| Renaming/dropping a column | Split into expand → backfill → contract releases |

> Forward-only is the rule. Reversibility comes from backward-compatible change
> design and tested backups — never from editing applied history or trusting an
> untested down-migration.
