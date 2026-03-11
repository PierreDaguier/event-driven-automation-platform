#!/usr/bin/env bash
set -euo pipefail

REPO="${1:-}"
if [[ -z "${REPO}" ]]; then
  echo "Usage: $0 <owner/repo>"
  exit 1
fi

labels=(
  "type:feature|0E8A16"
  "type:bug|D73A4A"
  "type:chore|6E7681"
  "type:docs|1D76DB"
  "type:test|BFD4F2"
  "type:refactor|5319E7"
  "priority:P0|B60205"
  "priority:P1|D93F0B"
  "priority:P2|FBCA04"
  "area:frontend|0052CC"
  "area:backend|5319E7"
  "area:infra|0E8A16"
  "area:observability|C2E0C6"
  "area:ux|FBCA04"
)

for entry in "${labels[@]}"; do
  name="${entry%%|*}"
  color="${entry##*|}"
  gh label create "${name}" --repo "${REPO}" --color "${color}" --force
  echo "label ensured: ${name}"
done

gh api --method PUT "repos/${REPO}/branches/main/protection" \
  -H "Accept: application/vnd.github+json" \
  -f required_linear_history=true \
  -f enforce_admins=true \
  -f required_pull_request_reviews.dismiss_stale_reviews=true \
  -f required_pull_request_reviews.required_approving_review_count=1 \
  -F required_status_checks.strict=true \
  -F required_status_checks.contexts[]='CI / api' \
  -F required_status_checks.contexts[]='CI / web'

echo "GitHub branch protection configured for ${REPO}" 
