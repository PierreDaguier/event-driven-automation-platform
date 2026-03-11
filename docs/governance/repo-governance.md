# Repository Governance

## Labels Standard

### `type:*`
- `type:feature` - New product functionality
- `type:bug` - Defect or regression
- `type:chore` - Maintenance task
- `type:docs` - Documentation changes
- `type:test` - Test-related work
- `type:refactor` - Internal code restructuring

### `priority:*`
- `priority:P0` - Highest priority
- `priority:P1` - Important but not urgent
- `priority:P2` - Nice-to-have / later

### `area:*`
- `area:frontend` - UI and web frontend
- `area:backend` - API and backend services
- `area:infra` - Infrastructure and deployment
- `area:observability` - Metrics, logs, traces
- `area:ux` - User experience and design

## Milestones

- `MVP` - Baseline platform capability
- `Reliability` - Hardening and operational resilience
- `UX Polish` - Client-facing UX quality
- `v1.0.0` - Production-grade demo release milestone

## Branch Protection Rules (`main`)

- Pull request required before merge
- Required status checks: `api`, `web` (strict mode)
- Required approving review count: `1`
- Enforce admins: `true`
- Required linear history: `true`
- Required conversation resolution: `true`
- Force pushes: disabled
- Deletions: disabled

## Governance Workflow

1. Create or link an issue before substantial work.
2. Assign at least one `type:*`, one `priority:*`, and one `area:*` label.
3. Attach the issue to a milestone.
4. Open PR against `main` with testing notes and risk assessment.
5. Merge only after required checks and one approval are satisfied.
