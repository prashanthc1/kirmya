# 14 · Secrets Management

Kirmya reads secrets through `os.Getenv` at the point of use (e.g. `JWT_SECRET`,
`MFA_ENC_KEY`, `DATABASE_URL`, `ANTHROPIC_API_KEY`, SMTP credentials). To keep
those call sites unchanged while still integrating a real secret manager in
production, all secret *sourcing* is funnelled through one seam:
`internal/platform/secrets`, invoked once at startup in `main.go` before
anything reads a secret.

## Backends (`SECRETS_BACKEND`)

**`env` (default).** Use the ambient environment / `.env`. This is the local-dev
and "I inject env vars myself" path. `secrets.Load` is a no-op.

**`file`.** Load a secret bundle from `SECRETS_FILE` (default
`/run/secrets/kirmya.env`) and export every key into the process environment.
This single mechanism integrates essentially every production secret manager,
because they all deliver secrets to a workload as a mounted file or directory:

- **Kubernetes Secrets** mounted as a volume.
- **HashiCorp Vault** via the Vault Agent sidecar / injector (renders secrets to
  a file).
- **AWS Secrets Manager** (and SSM Parameter Store) via the Secrets Store CSI
  driver, which mounts the secret as a file.
- **External Secrets Operator**, which syncs any of the above into a K8s Secret.
- **Docker / Swarm secrets**, mounted under `/run/secrets`.

The bundle may be a JSON object (`{"JWT_SECRET":"…","MFA_ENC_KEY":"…"}`) or
dotenv-style `KEY=VALUE` lines (with `#` comments and optional quotes).

### Precedence (`SECRETS_OVERRIDE`)

With the `file` backend, the secret store is treated as the source of truth and
**overrides** any pre-existing environment value (`SECRETS_OVERRIDE=true`, the
default). Set `SECRETS_OVERRIDE=false` to make the file only *fill in* keys that
are not already set in the environment.

## Fail-fast in production

When `APP_ENV=production`, `main.go` calls `secrets.Require("JWT_SECRET",
"MFA_ENC_KEY")` after loading, so a missing critical secret aborts the boot with
a clear message instead of silently falling back to an insecure dev key. Secret
*values* are never logged — only the count loaded and the source path.

## Wiring (docker-compose.prod.yml)

The production compose file passes `SECRETS_BACKEND` and `SECRETS_FILE` through
from the environment, defaulting to `env`. To switch a deployment to file-based
secrets, set `SECRETS_BACKEND=file` and mount the bundle at `SECRETS_FILE`
(e.g. a Docker secret or a CSI-driver volume) — no image or code change needed.

## Adding another backend

To integrate a secret manager via its API directly (rather than a mounted file),
add a `case` to `secrets.Load` in `internal/platform/secrets/secrets.go` that
fetches the bundle and returns a `map[string]string` through the same export
path. Keep the package dependency-light: prefer the file/CSI integration above
unless a direct API pull is genuinely required, since the file approach already
covers Kubernetes, Vault, AWS, and Docker without adding SDK dependencies.
