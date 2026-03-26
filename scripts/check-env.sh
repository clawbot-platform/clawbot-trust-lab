#!/usr/bin/env sh
set -eu

ENV_FILE="${1:-.env}"
VALIDATE_OPTIONAL_STACK="${VALIDATE_OPTIONAL_STACK:-0}"

if [ ! -f "$ENV_FILE" ]; then
  echo "Missing $ENV_FILE. Copy .env.example to $ENV_FILE first."
  exit 1
fi

core_required_vars="
COMPOSE_PROJECT_NAME
POSTGRES_IMAGE
POSTGRES_HOST
POSTGRES_PORT
POSTGRES_DB
POSTGRES_USER
POSTGRES_PASSWORD
CONTROL_PLANE_HOST
CONTROL_PLANE_PORT
CLAWMEM_HOST
CLAWMEM_PORT
TRUST_LAB_HOST
TRUST_LAB_PORT
TRUST_LAB_UI_HOST
TRUST_LAB_UI_PORT
STACK_SMOKE_TIMEOUT
APP_ENV
LOG_LEVEL
SHUTDOWN_TIMEOUT
CONTROL_PLANE_BASE_URL
CONTROL_PLANE_TIMEOUT
CLAWMEM_BASE_URL
CLAWMEM_TIMEOUT
SCENARIO_PACKS_DIR
REPLAY_ARCHIVE_DIR
REPORTS_DIR
BENCHMARK_SCHEDULER_ENABLED
BENCHMARK_SCHEDULER_SCENARIO_FAMILY
BENCHMARK_SCHEDULER_INTERVAL
BENCHMARK_SCHEDULER_MAX_RUNS
BENCHMARK_SCHEDULER_DRY_RUN
"

optional_required_vars=""

# shellcheck disable=SC1090
set -a
case "$ENV_FILE" in
  */*) . "$ENV_FILE" ;;
  *) . "./$ENV_FILE" ;;
esac
set +a

validate_vars() {
  required_vars="$1"

  for var in $required_vars; do
    eval "value=\${$var:-}"
    if [ -z "$value" ]; then
      echo "ERROR: Required variable '$var' is missing or empty."
      exit 1
    fi

    case "$value" in
      *REQUIRED_SECRET*|*replace_me*)
        echo "ERROR: Variable '$var' still contains a placeholder value."
        exit 1
        ;;
    esac
  done
}

validate_paths() {
  for path in \
    "../clawbot-server" \
    "../clawmem" \
    "./configs/scenario-packs"; do
    if [ ! -e "$path" ]; then
      echo "ERROR: Required path '$path' is missing. Version 1 Docker builds expect sibling checkouts of clawbot-server and clawmem."
      exit 1
    fi
  done
}

validate_vars "$core_required_vars"
validate_paths

if [ "$VALIDATE_OPTIONAL_STACK" = "1" ] && [ -n "$optional_required_vars" ]; then
  validate_vars "$optional_required_vars"
  echo "Environment file validation passed for core + optional stack: $ENV_FILE"
else
  echo "Environment file validation passed for core stack: $ENV_FILE"
fi
