#!/bin/sh
# Vérifie le nom de branche et les messages de commit d'une merge request (voir CONTRIBUTING.md)
set -eu

BRANCH="${CI_MERGE_REQUEST_SOURCE_BRANCH_NAME:-$(git rev-parse --abbrev-ref HEAD)}"
BASE="${CI_MERGE_REQUEST_DIFF_BASE_SHA:-origin/main}"

BRANCH_RE='^(feat|fix|docs|refactor|test|ci|build|chore|hotfix)/[a-z0-9._-]+$'
COMMIT_RE='^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\([a-z0-9-]+\))?!?: .+'

status=0

if echo "$BRANCH" | grep -Eq "$BRANCH_RE"; then
  echo "OK   branche : $BRANCH"
else
  echo "FAIL branche : '$BRANCH' ne respecte pas <type>/<description> (ex : feat/api-products)"
  status=1
fi

for sha in $(git rev-list --no-merges "$BASE..HEAD"); do
  subject=$(git log -1 --format=%s "$sha")
  short=$(git rev-parse --short "$sha")
  if echo "$subject" | grep -Eq "$COMMIT_RE"; then
    echo "OK   $short $subject"
  else
    echo "FAIL $short $subject"
    status=1
  fi
done

if [ "$status" -ne 0 ]; then
  echo
  echo "Format attendu : <type>(<scope>): <description>"
  echo "Types : feat fix docs style refactor perf test build ci chore revert"
fi
exit "$status"
