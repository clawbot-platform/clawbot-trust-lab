#!/usr/bin/env python3
from __future__ import annotations

import argparse
import datetime as dt
import html
import json
import os
import re
import shlex
import subprocess
import sys
import textwrap
import urllib.error
import urllib.request
from dataclasses import dataclass, field
from pathlib import Path
from typing import Any


@dataclass
class CheckResult:
    name: str
    kind: str  # command, api, file, doc
    passed: bool
    summary: str
    details: str = ""
    command: str | None = None
    url: str | None = None
    metadata: dict[str, Any] = field(default_factory=dict)


def now_iso() -> str:
    return dt.datetime.now().astimezone().isoformat(timespec="seconds")


def run_cmd(name: str, cmd: list[str], cwd: Path, timeout: int = 900, env: dict[str, str] | None = None) -> CheckResult:
    try:
        proc = subprocess.run(
            cmd,
            cwd=str(cwd),
            env=env,
            capture_output=True,
            text=True,
            timeout=timeout,
        )
        passed = proc.returncode == 0
        out = (proc.stdout or "") + ("\n" if proc.stdout and proc.stderr else "") + (proc.stderr or "")
        summary = f"exit={proc.returncode}"
        return CheckResult(
            name=name,
            kind="command",
            passed=passed,
            summary=summary,
            details=out.strip(),
            command=" ".join(shlex.quote(c) for c in cmd),
            metadata={"returncode": proc.returncode},
        )
    except FileNotFoundError as e:
        return CheckResult(
            name=name,
            kind="command",
            passed=False,
            summary="command not found",
            details=str(e),
            command=" ".join(shlex.quote(c) for c in cmd),
        )
    except subprocess.TimeoutExpired as e:
        details = (e.stdout or "") + ("\n" if e.stdout and e.stderr else "") + (e.stderr or "")
        return CheckResult(
            name=name,
            kind="command",
            passed=False,
            summary=f"timeout after {timeout}s",
            details=details.strip(),
            command=" ".join(shlex.quote(c) for c in cmd),
        )


def run_cmd_text(name: str, cmd: list[str], cwd: Path, timeout: int = 900, env: dict[str, str] | None = None) -> CheckResult:
    result = run_cmd(name, cmd, cwd, timeout=timeout, env=env)
    result.kind = "command"
    return result


def http_json(name: str, method: str, url: str, payload: dict[str, Any] | None = None, timeout: int = 30) -> CheckResult:
    data = None
    headers = {"Accept": "application/json"}
    if payload is not None:
        data = json.dumps(payload).encode("utf-8")
        headers["Content-Type"] = "application/json"
    req = urllib.request.Request(url, method=method.upper(), data=data, headers=headers)
    try:
        with urllib.request.urlopen(req, timeout=timeout) as resp:
            body = resp.read().decode("utf-8", errors="replace")
            try:
                parsed = json.loads(body)
                pretty = json.dumps(parsed, indent=2, ensure_ascii=False)
            except json.JSONDecodeError:
                parsed = None
                pretty = body
            return CheckResult(
                name=name,
                kind="api",
                passed=200 <= resp.status < 300,
                summary=f"http {resp.status}",
                details=pretty,
                url=url,
                metadata={"status": resp.status, "json": parsed},
            )
    except urllib.error.HTTPError as e:
        body = e.read().decode("utf-8", errors="replace")
        return CheckResult(
            name=name,
            kind="api",
            passed=False,
            summary=f"http {e.code}",
            details=body.strip(),
            url=url,
            metadata={"status": e.code},
        )
    except Exception as e:
        return CheckResult(
            name=name,
            kind="api",
            passed=False,
            summary="request failed",
            details=str(e),
            url=url,
        )


def http_text(name: str, url: str, timeout: int = 30) -> CheckResult:
    req = urllib.request.Request(url, method="GET", headers={"Accept": "*/*"})
    try:
        with urllib.request.urlopen(req, timeout=timeout) as resp:
            body = resp.read().decode("utf-8", errors="replace")
            return CheckResult(
                name=name,
                kind="api",
                passed=200 <= resp.status < 300,
                summary=f"http {resp.status}",
                details=body[:12000],
                url=url,
                metadata={"status": resp.status},
            )
    except urllib.error.HTTPError as e:
        body = e.read().decode("utf-8", errors="replace")
        return CheckResult(
            name=name,
            kind="api",
            passed=False,
            summary=f"http {e.code}",
            details=body.strip(),
            url=url,
            metadata={"status": e.code},
        )
    except Exception as e:
        return CheckResult(
            name=name,
            kind="api",
            passed=False,
            summary="request failed",
            details=str(e),
            url=url,
        )


def file_exists(name: str, path: Path) -> CheckResult:
    exists = path.exists()
    return CheckResult(
        name=name,
        kind="file",
        passed=exists,
        summary="exists" if exists else "missing",
        details=str(path),
        metadata={"path": str(path)},
    )


def text_contains(name: str, path: Path, patterns: list[str]) -> CheckResult:
    if not path.exists():
        return CheckResult(name=name, kind="doc", passed=False, summary="missing file", details=str(path))
    text = path.read_text(encoding="utf-8", errors="replace")
    missing = [p for p in patterns if p not in text]
    return CheckResult(
        name=name,
        kind="doc",
        passed=not missing,
        summary="all patterns found" if not missing else f"missing {len(missing)} pattern(s)",
        details=("Missing:\n- " + "\n- ".join(missing)) if missing else str(path),
        metadata={"path": str(path), "missing": missing},
    )


def summarize_scenarios(result: CheckResult) -> CheckResult:
    parsed = result.metadata.get("json") if result.metadata else None
    if not isinstance(parsed, dict):
        result.passed = False
        result.summary += "; invalid json"
        return result
    data = parsed.get("data")
    if not isinstance(data, list):
        result.passed = False
        result.summary += "; missing data list"
        return result
    ids = [str(item.get("id", "")) for item in data if isinstance(item, dict)]
    categories = {
        "human": [i for i in ids if re.search(r"(^|-)h[12](-|$)|human", i)],
        "benign_agentic": [i for i in ids if re.search(r"(^|-)a[123](-|$)|agent-assisted|delegated", i)],
        "suspicious": [i for i in ids if re.search(r"(^|-)s[1-5](-|$)|weak|refund|scope-drift|repeated", i)],
        "variants": [i for i in ids if re.search(r"(^|-)v[1-7](-|$)|variant|challenger", i)],
    }
    missing_groups = [k for k, v in categories.items() if not v]
    result.passed = result.passed and not missing_groups
    result.summary += f"; scenarios={len(ids)}; missing_groups={','.join(missing_groups) or 'none'}"
    result.metadata["scenario_ids"] = ids
    result.metadata["scenario_groups"] = categories
    return result


def summarize_rounds(result: CheckResult) -> CheckResult:
    parsed = result.metadata.get("json") if result.metadata else None
    if not isinstance(parsed, dict):
        result.passed = False
        return result
    data = parsed.get("data")
    if not isinstance(data, list):
        result.passed = False
        result.summary += "; missing data list"
        return result
    round_ids = [str(item.get("id", "")) for item in data if isinstance(item, dict)]
    rec_counts = []
    prod_bridge_ok = 0
    tier_c_seen = 0
    for item in data:
        if not isinstance(item, dict):
            continue
        summary = item.get("summary") or {}
        if isinstance(summary, dict):
            rec_counts.append(summary.get("recommendations", 0) or 0)
            if summary.get("evaluation_mode") and summary.get("blocking_mode"):
                prod_bridge_ok += 1
            if summary.get("tier_c_usage_count", 0):
                tier_c_seen += 1
    result.summary += f"; rounds={len(round_ids)}; rounds_with_prod_bridge={prod_bridge_ok}; rounds_with_tier_c_usage={tier_c_seen}; recommendation_totals={sum(rec_counts)}"
    result.metadata["round_ids"] = round_ids
    return result


def summarize_promotions(result: CheckResult) -> CheckResult:
    parsed = result.metadata.get("json") if result.metadata else None
    if not isinstance(parsed, dict):
        result.passed = False
        return result
    data = parsed.get("data")
    if not isinstance(data, list):
        result.passed = False
        result.summary += "; missing data list"
        return result
    historical = sorted({str(item.get("round_id", "")) for item in data if isinstance(item, dict) and item.get("round_id")})
    result.summary += f"; promotions={len(data)}; distinct_rounds={len(historical)}"
    result.metadata["round_ids"] = historical
    return result


def summarize_recommendations(result: CheckResult) -> CheckResult:
    parsed = result.metadata.get("json") if result.metadata else None
    if not isinstance(parsed, dict):
        result.passed = False
        return result
    data = parsed.get("data")
    if not isinstance(data, list):
        result.passed = False
        result.summary += "; missing data list"
        return result
    types = sorted({str(item.get("type", "")) for item in data if isinstance(item, dict) and item.get("type")})
    result.passed = result.passed and bool(types)
    result.summary += f"; recommendations={len(data)}; types={', '.join(types) if types else 'none'}"
    return result


def summarize_trends(result: CheckResult) -> CheckResult:
    parsed = result.metadata.get("json") if result.metadata else None
    if not isinstance(parsed, dict):
        result.passed = False
        return result
    data = parsed.get("data")
    if not isinstance(data, dict):
        result.passed = False
        result.summary += "; missing data object"
        return result
    keys = [
        "rounds_executed",
        "promotions_over_time",
        "replay_pass_rate_over_time",
        "new_blind_spots_discovered",
        "recommendation_counts_by_type",
    ]
    present = [k for k in keys if k in data]
    missing = [k for k in keys if k not in data]
    result.passed = result.passed and len(present) >= 3
    result.summary += f"; present_keys={len(present)}; missing={','.join(missing) or 'none'}"
    return result


def summarize_scheduler(result: CheckResult) -> CheckResult:
    parsed = result.metadata.get("json") if result.metadata else None
    if not isinstance(parsed, dict):
        result.passed = False
        return result
    data = parsed.get("data")
    if not isinstance(data, dict):
        result.passed = False
        result.summary += "; missing data object"
        return result
    mode = data.get("mode") or data.get("status") or "unknown"
    interval = data.get("interval_seconds") or data.get("interval")
    result.summary += f"; scheduler_mode={mode}; interval={interval}"
    return result


def summarize_report_files(name: str, reports_dir: Path) -> CheckResult:
    required = {
        "round-summary.json",
        "round-summary.md",
        "detection-delta.json",
        "promotion-report.json",
        "executive-summary.md",
        "recommendation-report.json",
    }
    if not reports_dir.exists():
        return CheckResult(name=name, kind="file", passed=False, summary="reports dir missing", details=str(reports_dir))
    round_dirs = [p for p in reports_dir.iterdir() if p.is_dir() and p.name.startswith("round-")]
    missing_by_round: dict[str, list[str]] = {}
    legacy_reconstructible: dict[str, str] = {}
    for rd in round_dirs:
        present = {p.name for p in rd.iterdir() if p.is_file()}
        missing = sorted(required - present)
        if missing:
            if len(missing) == 1 and missing[0] == "recommendation-report.json":
                reconstructible, reason = is_legacy_recommendation_gap(rd)
                if reconstructible:
                    legacy_reconstructible[rd.name] = reason
                    continue
            missing_by_round[rd.name] = missing
    passed = bool(round_dirs) and not missing_by_round
    details = ""
    if missing_by_round or legacy_reconstructible:
        parts = []
        for rid, missing in sorted(missing_by_round.items()):
            parts.append(f"{rid}: missing {', '.join(missing)}")
        for rid, reason in sorted(legacy_reconstructible.items()):
            parts.append(f"{rid}: legacy reconstructible gap for recommendation-report.json ({reason})")
        details = "\n".join(parts)
    else:
        details = f"checked {len(round_dirs)} round director{'y' if len(round_dirs)==1 else 'ies'}"
    return CheckResult(
        name=name,
        kind="file",
        passed=passed,
        summary=f"round_dirs={len(round_dirs)}; missing_rounds={len(missing_by_round)}; legacy_reconstructible={len(legacy_reconstructible)}",
        details=details,
    )


def is_legacy_recommendation_gap(round_dir: Path) -> tuple[bool, str]:
    summary_path = round_dir / "round-summary.json"
    if not summary_path.exists():
        return False, "round-summary.json missing"

    try:
        data = json.loads(summary_path.read_text(encoding="utf-8"))
    except Exception:
        return False, "round-summary.json unreadable"

    embedded_recommendations = data.get("recommendations")
    if isinstance(embedded_recommendations, list) and embedded_recommendations:
        return True, "embedded recommendations present in round-summary.json"

    scenario_results = data.get("scenario_results")
    if isinstance(scenario_results, list) and scenario_results:
        return True, "scenario_results present for deterministic bootstrap reconstruction"

    if (round_dir / "promotion-report.json").exists() or (round_dir / "detection-delta.json").exists():
        return True, "promotion/delta artifacts present for recommendation bootstrap enrichment"

    return False, "insufficient persisted artifacts"


def md_escape_code(text: str) -> str:
    return text.replace("```", "``\\`")


def to_markdown(results: list[CheckResult], meta: dict[str, Any]) -> str:
    total = len(results)
    passed = sum(1 for r in results if r.passed)
    failed = total - passed
    lines = []
    lines.append("# Clawbot Trust Lab Version 1 Validation Report")
    lines.append("")
    lines.append(f"Generated: {meta['generated_at']}")
    lines.append("")
    lines.append(f"Repo root: `{meta['repo_root']}`")
    lines.append("")
    lines.append(f"Summary: **{passed} passed / {failed} failed / {total} total**")
    lines.append("")
    lines.append("## Checks")
    lines.append("")
    for r in results:
        status = "PASS" if r.passed else "FAIL"
        lines.append(f"### [{status}] {r.name}")
        lines.append("")
        lines.append(f"- Kind: `{r.kind}`")
        lines.append(f"- Summary: {r.summary}")
        if r.command:
            lines.append(f"- Command: `{r.command}`")
        if r.url:
            lines.append(f"- URL: `{r.url}`")
        if r.details:
            lines.append("")
            lines.append("```text")
            lines.append(md_escape_code(r.details[:12000]))
            lines.append("```")
        lines.append("")
    return "\n".join(lines)


def to_html(results: list[CheckResult], meta: dict[str, Any]) -> str:
    total = len(results)
    passed = sum(1 for r in results if r.passed)
    failed = total - passed
    blocks = []
    for r in results:
        css = "pass" if r.passed else "fail"
        details = html.escape(r.details[:12000]) if r.details else ""
        cmd = f"<div><b>Command:</b> <code>{html.escape(r.command)}</code></div>" if r.command else ""
        url = f"<div><b>URL:</b> <code>{html.escape(r.url)}</code></div>" if r.url else ""
        blocks.append(f"""
        <section class=\"card {css}\">
          <h2>[{'PASS' if r.passed else 'FAIL'}] {html.escape(r.name)}</h2>
          <div><b>Kind:</b> {html.escape(r.kind)}</div>
          <div><b>Summary:</b> {html.escape(r.summary)}</div>
          {cmd}
          {url}
          {'<pre>' + details + '</pre>' if details else ''}
        </section>
        """)
    return f"""<!doctype html>
<html lang=\"en\"><head><meta charset=\"utf-8\"><title>Clawbot Trust Lab Version 1 Validation Report</title>
<style>
body {{ font-family: Arial, sans-serif; margin: 24px; background: #fafafa; color: #222; }}
summary, h1, h2 {{ color: #111; }}
.stats {{ display:flex; gap:16px; margin: 16px 0 24px; }}
.badge {{ padding: 8px 12px; border-radius: 999px; background: #eee; }}
.pass .badge, .pass h2 {{ color: #0a5; }}
.fail .badge, .fail h2 {{ color: #b00; }}
.card {{ background:#fff; border:1px solid #ddd; border-left: 6px solid #ccc; padding:16px; margin: 16px 0; border-radius: 8px; }}
.card.pass {{ border-left-color: #0a5; }}
.card.fail {{ border-left-color: #b00; }}
pre {{ background:#111; color:#eee; padding:12px; overflow:auto; border-radius: 6px; white-space: pre-wrap; }}
code {{ background:#f0f0f0; padding: 2px 4px; border-radius: 4px; }}
</style></head>
<body>
<h1>Clawbot Trust Lab Version 1 Validation Report</h1>
<div>Generated: {html.escape(meta['generated_at'])}</div>
<div>Repo root: <code>{html.escape(meta['repo_root'])}</code></div>
<div class=\"stats\">
  <div class=\"badge\">Passed: {passed}</div>
  <div class=\"badge\">Failed: {failed}</div>
  <div class=\"badge\">Total: {total}</div>
</div>
{''.join(blocks)}
</body></html>"""


def main() -> int:
    ap = argparse.ArgumentParser(description="Run Clawbot Trust Lab Version 1 validation checks and generate Markdown + HTML reports.")
    ap.add_argument("--repo-root", default=".", help="Path to clawbot-trust-lab repo root")
    ap.add_argument("--api-base", default="http://127.0.0.1:8090", help="Base URL for trust-lab API")
    ap.add_argument("--ui-base", default="http://127.0.0.1:8091", help="Base URL for the optional operator UI")
    ap.add_argument("--deployment-mode", choices=["local", "docker"], default="local", help="Validate a local source run or the Docker-based Version 1 stack")
    ap.add_argument("--compose-file", default="docker-compose.v1.yml", help="Compose file used for Version 1 Docker deployment checks")
    ap.add_argument("--compose-env-file", default="docker-compose.v1.env", help="Compose env file used for Version 1 Docker deployment checks")
    ap.add_argument("--skip-compose-checks", action="store_true", help="Skip docker compose checks even when deployment mode is docker")
    ap.add_argument("--skip-backend", action="store_true")
    ap.add_argument("--skip-web", action="store_true")
    ap.add_argument("--skip-api", action="store_true")
    ap.add_argument("--run-round", action="store_true", help="POST a fresh benchmark round during validation")
    ap.add_argument("--output-dir", default="version1-validation-output")
    args = ap.parse_args()

    repo_root = Path(args.repo_root).resolve()
    out_dir = Path(args.output_dir).resolve()
    out_dir.mkdir(parents=True, exist_ok=True)

    results: list[CheckResult] = []

    # Documentation / contract checks
    results.append(file_exists("README exists", repo_root / "README.md"))
    results.append(file_exists("api.md exists", repo_root / "docs" / "api.md"))
    results.append(file_exists("Version 1 deployment guide exists", repo_root / "docs" / "deploying-clawbot-trust-lab-v1.md"))
    results.append(file_exists("scenario catalog exists", repo_root / "docs" / "phase-9-scenario-catalog.md"))
    results.append(text_contains(
        "README presents Version 1 and planned Version 2 clearly",
        repo_root / "README.md",
        ["Version 1", "Docker", "phase9_validation_report.py", "Planned Version 2"],
    ))
    results.append(text_contains(
        "api.md includes Recommendation / Trend / Scheduler schemas",
        repo_root / "docs" / "api.md",
        ["Recommendation", "TrendSummary", "SchedulerStatus", "SchedulerRunResponse", "legacy aliases", "recommendation-report.json"],
    ))
    if args.deployment_mode == "docker":
        results.append(file_exists("Version 1 compose file exists", repo_root / args.compose_file))
        results.append(file_exists("Version 1 compose env file exists", repo_root / args.compose_env_file))
        results.append(file_exists("trust-lab Dockerfile exists", repo_root / "deploy" / "docker" / "clawbot-trust-lab.Dockerfile"))
        results.append(file_exists("operator UI Dockerfile exists", repo_root / "deploy" / "docker" / "operator-ui" / "Dockerfile"))

    if args.deployment_mode == "docker" and not args.skip_compose_checks:
        compose_env = os.environ.copy()
        compose_cmd = [
            "docker",
            "compose",
            "--env-file",
            args.compose_env_file,
            "-f",
            args.compose_file,
            "ps",
        ]
        results.append(run_cmd_text("docker compose ps", compose_cmd, repo_root, timeout=120, env=compose_env))

    if not args.skip_backend:
        backend_cmds = [
            ("go test ./...", ["go", "test", "./..."], 1800),
            ("go vet ./...", ["go", "vet", "./..."], 1200),
            ("golangci-lint run ./...", ["golangci-lint", "run", "./..."], 1800),
            ("gosec ./...", ["gosec", "./..."], 1800),
            ("govulncheck ./...", ["govulncheck", "./..."], 1800),
        ]
        for name, cmd, timeout in backend_cmds:
            results.append(run_cmd(name, cmd, repo_root, timeout=timeout, env=os.environ.copy()))

    web_dir = repo_root / "web"
    if not args.skip_web and web_dir.exists():
        web_cmds = [
            ("npm run lint", ["npm", "run", "lint"], 1200),
            ("npm run test", ["npm", "run", "test"], 1800),
            ("npm run build", ["npm", "run", "build"], 1800),
            ("npm run test:e2e", ["npm", "run", "test:e2e"], 2400),
        ]
        for name, cmd, timeout in web_cmds:
            results.append(run_cmd(name, cmd, web_dir, timeout=timeout, env=os.environ.copy()))

    reports_dir = repo_root / "reports"
    results.append(summarize_report_files("reports directory contains expected artifacts", reports_dir))

    if not args.skip_api:
        api = args.api_base.rstrip("/")
        api_checks = [
            http_json("GET /healthz", "GET", f"{api}/healthz"),
            http_json("GET /readyz", "GET", f"{api}/readyz"),
            summarize_scenarios(http_json("GET /api/v1/scenarios", "GET", f"{api}/api/v1/scenarios")),
        ]
        if args.run_round:
            api_checks.append(http_json(
                "POST /api/v1/benchmark/rounds/run",
                "POST",
                f"{api}/api/v1/benchmark/rounds/run",
                payload={"scenario_family": "commerce"},
            ))
        api_checks.extend([
            summarize_rounds(http_json("GET /api/v1/benchmark/rounds", "GET", f"{api}/api/v1/benchmark/rounds")),
            summarize_rounds(http_json("GET /api/v1/operator/rounds", "GET", f"{api}/api/v1/operator/rounds")),
            summarize_promotions(http_json("GET /api/v1/operator/promotions", "GET", f"{api}/api/v1/operator/promotions")),
            summarize_recommendations(http_json("GET /api/v1/benchmark/recommendations", "GET", f"{api}/api/v1/benchmark/recommendations")),
            summarize_recommendations(http_json("GET /api/v1/operator/recommendations", "GET", f"{api}/api/v1/operator/recommendations")),
            summarize_trends(http_json("GET /api/v1/benchmark/trends/summary", "GET", f"{api}/api/v1/benchmark/trends/summary")),
            summarize_trends(http_json("GET /api/v1/operator/trends/summary", "GET", f"{api}/api/v1/operator/trends/summary")),
            summarize_scheduler(http_json("GET /api/v1/benchmark/scheduler/status", "GET", f"{api}/api/v1/benchmark/scheduler/status")),
        ])
        if args.deployment_mode == "docker":
            api_checks.append(http_text("GET operator UI root", args.ui_base.rstrip("/") + "/"))
        results.extend(api_checks)

    meta = {
        "generated_at": now_iso(),
        "repo_root": str(repo_root),
    }

    md = to_markdown(results, meta)
    html_doc = to_html(results, meta)
    md_path = out_dir / "version1-validation-report.md"
    html_path = out_dir / "version1-validation-report.html"
    md_path.write_text(md, encoding="utf-8")
    html_path.write_text(html_doc, encoding="utf-8")
    (out_dir / "version1-validation-report.md").write_text(md, encoding="utf-8")
    (out_dir / "version1-validation-report.html").write_text(html_doc, encoding="utf-8")

    print(f"Wrote Markdown report: {md_path}")
    print(f"Wrote HTML report: {html_path}")
    failed = sum(1 for r in results if not r.passed)
    return 1 if failed else 0


if __name__ == "__main__":
    sys.exit(main())
