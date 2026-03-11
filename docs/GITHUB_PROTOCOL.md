# GitHub Protocol (Team-style)

## Branch strategy

- `main`: protected, release only.
- `develop`: integration branch.
- `feat/*`, `fix/*`, `chore/*`, `docs/*` for scoped work.

Examples:
- `feat/workflow-engine-core`
- `feat/ui-dashboard-runs`
- `fix/webhook-signature-validation`

## Conventional commits

- `feat:`, `fix:`, `chore:`, `docs:`, `test:`, `refactor:`, `ci:`

Examples:
- `feat(api): add webhook signature verification middleware`
- `feat(worker): implement retry with exponential backoff`
- `fix(ui): handle empty run history state`

## Pull Request checklist

- Why
- What changed
- Screenshots/GIFs
- Testing notes
- Risks & rollback
- Linked issues (`Closes #...`)

## Labels

Managed by `scripts/setup-github.sh`:
- `type:feature`, `type:bug`, `type:chore`, `type:docs`, `type:test`, `type:refactor`
- `priority:P0`, `priority:P1`, `priority:P2`
- `area:frontend`, `area:backend`, `area:infra`, `area:observability`, `area:ux`

## Repo hardening script

```bash
./scripts/setup-github.sh owner/event-driven-automation-platform
```

This script configures:
- labels,
- branch protection on `main` (PR required, status checks, linear history).

## Suggested release cadence

- `v0.1.0` MVP
- `v0.2.0` observability + UX polish
- `v1.0.0` client-facing release
