#!/usr/bin/env sh
set -eu

SOURCE_ENV="${1:-.env.example}"
TARGET_ENV="${2:-.env}"

if [ ! -f "$SOURCE_ENV" ]; then
  echo "Missing source env file: $SOURCE_ENV"
  exit 1
fi

cp "$SOURCE_ENV" "$TARGET_ENV"

project_name="${COMPOSE_PROJECT_NAME_OVERRIDE:-clawbot-trust-lab-ci}"
if [ -n "${GITHUB_RUN_ID:-}" ]; then
  project_name="clawbot-trust-lab-ci-${GITHUB_RUN_ID}"
fi

tmp_file="$(mktemp)"
awk -v project_name="$project_name" '
BEGIN { replaced = 0 }
/^COMPOSE_PROJECT_NAME=/ {
  print "COMPOSE_PROJECT_NAME=" project_name
  replaced = 1
  next
}
{ print }
END {
  if (replaced == 0) {
    print "COMPOSE_PROJECT_NAME=" project_name
  }
}
' "$TARGET_ENV" > "$tmp_file"
mv "$tmp_file" "$TARGET_ENV"

mkdir -p reports var/replay-archive var/docker/clawmem version1-validation-output

echo "Prepared CI env file at $TARGET_ENV"
