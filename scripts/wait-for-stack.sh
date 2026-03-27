#!/usr/bin/env sh
set -eu

timeout_seconds="${STACK_SMOKE_TIMEOUT:-${1:-120}}"
deadline=$(( $(date +%s) + timeout_seconds ))

wait_for_url() {
  name="$1"
  url="$2"

  if curl -fsS "$url" >/dev/null 2>&1; then
    echo "$name ready: $url"
    return 0
  fi
  return 1
}

while [ "$(date +%s)" -lt "${deadline}" ]; do
  if wait_for_url "control-plane" "http://127.0.0.1:8081/healthz" && \
     wait_for_url "clawmem" "http://127.0.0.1:8088/healthz" && \
     wait_for_url "trust-lab" "http://127.0.0.1:8090/readyz" && \
     wait_for_url "trust-lab-ui" "http://127.0.0.1:8091/"; then
    exit 0
  fi

  sleep 5
done

echo "core stack did not become ready within ${timeout_seconds}s" >&2
exit 1
