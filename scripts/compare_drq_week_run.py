#!/usr/bin/env python3
import argparse
import json
from pathlib import Path
from collections import Counter
from statistics import mean
from typing import Any, Optional

TARGET_CASES = [
    "commerce-v2-expired-inactive-mandate",
    "commerce-v3-approval-removed",
    "commerce-s3-approval-removed-after-authorization",
]


def load_json(path: Path):
    with path.open("r", encoding="utf-8") as f:
        return json.load(f)


def resolve_root(root: Path) -> Path:
    direct = sorted(root.glob("round-*"))
    if direct:
        return root

    nested = root / "reports"
    if nested.exists():
        return nested

    return root


def round_dirs(root: Path):
    return sorted([p for p in root.glob("round-*") if p.is_dir()])


def walk_nodes(obj: Any):
    yield obj
    if isinstance(obj, dict):
        for value in obj.values():
            yield from walk_nodes(value)
    elif isinstance(obj, list):
        for item in obj:
            yield from walk_nodes(item)


def is_round_id(value: Any) -> bool:
    return isinstance(value, str) and value.startswith("round-")


def score_round_candidate(d: dict) -> int:
    score = 0

    if is_round_id(d.get("id")):
        score += 10
    if isinstance(d.get("summary"), dict):
        score += 8
    if isinstance(d.get("scenario_results"), list):
        score += 8
    if isinstance(d.get("detector_version"), str):
        score += 4
    if isinstance(d.get("started_at"), str):
        score += 2
    if isinstance(d.get("completed_at"), str):
        score += 2
    if isinstance(d.get("round_status"), str):
        score += 2
    if isinstance(d.get("reports"), dict):
        score += 1

    return score


def find_best_round_object(obj: Any) -> Optional[dict]:
    best = None
    best_score = -1

    for node in walk_nodes(obj):
        if not isinstance(node, dict):
            continue
        score = score_round_candidate(node)
        if score > best_score:
            best = node
            best_score = score

    return best if best_score > 0 else None


def is_summary_dict(d: dict) -> bool:
    required = {"round_id", "promotion_count", "replay_pass_rate", "robustness_outcome"}
    return required.issubset(set(d.keys()))


def find_summary_dict(obj: Any) -> Optional[dict]:
    for node in walk_nodes(obj):
        if isinstance(node, dict) and is_summary_dict(node):
            return node
    return None


def score_scenario_results_candidate(lst: list) -> int:
    if not isinstance(lst, list) or not lst:
        return -1
    if not all(isinstance(item, dict) for item in lst):
        return -1

    hits = 0
    for item in lst:
        if {"scenario_id", "final_detection_status", "passed"}.issubset(set(item.keys())):
            hits += 1

    if hits == 0:
        return -1

    return hits


def find_scenario_results_list(obj: Any) -> list:
    best = []
    best_score = -1

    for node in walk_nodes(obj):
        if isinstance(node, list):
            score = score_scenario_results_candidate(node)
            if score > best_score:
                best = node
                best_score = score

    return best if best_score > 0 else []


def score_recommendation_list_candidate(lst: list) -> int:
    if not isinstance(lst, list) or not lst:
        return -1
    if not all(isinstance(item, dict) for item in lst):
        return -1

    hits = 0
    for item in lst:
        if "type" in item:
            hits += 1

    return hits if hits > 0 else -1


def find_recommendation_list(obj: Any) -> list:
    best = []
    best_score = -1

    for node in walk_nodes(obj):
        if isinstance(node, list):
            score = score_recommendation_list_candidate(node)
            if score > best_score:
                best = node
                best_score = score

    return best if best_score > 0 else []


def classify_phase(round_id: str, phase_b_start_round: str) -> str:
    return "phase_b" if round_id >= phase_b_start_round else "phase_a"


def summarize_rounds(rounds: list[dict]) -> dict:
    if not rounds:
        return {
            "count": 0,
            "first_round_id": None,
            "last_round_id": None,
            "detector_versions": [],
            "total_promotions": 0,
            "avg_promotions": 0.0,
            "avg_replay_pass_rate": 0.0,
            "zero_promotion_rounds": 0,
            "perfect_replay_rounds": 0,
            "robustness_outcomes": {},
            "recommendation_counts": {},
            "latest_target_case_statuses": {},
        }

    rounds_sorted = sorted(rounds, key=lambda r: r["round_id"])

    promotion_values = []
    replay_values = []
    robustness = Counter()
    rec_counts = Counter()
    detector_versions = sorted({r["detector_version"] for r in rounds_sorted if r["detector_version"]})

    for r in rounds_sorted:
        promotion_values.append(r["promotion_count"])
        replay_values.append(r["replay_pass_rate"])
        robustness[r["robustness_outcome"]] += 1
        rec_counts.update(r["recommendation_types"])

    latest = rounds_sorted[-1]
    latest_statuses = latest["target_case_statuses"]

    return {
        "count": len(rounds_sorted),
        "first_round_id": rounds_sorted[0]["round_id"],
        "last_round_id": rounds_sorted[-1]["round_id"],
        "detector_versions": detector_versions,
        "total_promotions": sum(promotion_values),
        "avg_promotions": mean(promotion_values) if promotion_values else 0.0,
        "avg_replay_pass_rate": mean(replay_values) if replay_values else 0.0,
        "zero_promotion_rounds": sum(1 for v in promotion_values if v == 0),
        "perfect_replay_rounds": sum(1 for v in replay_values if v == 1),
        "robustness_outcomes": dict(robustness),
        "recommendation_counts": dict(rec_counts),
        "latest_target_case_statuses": latest_statuses,
    }


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--root", required=True, help="Path to local reports root")
    parser.add_argument(
        "--phase-b-start-round",
        required=True,
        help="Round ID where Phase B starts, e.g. round-20260401070127",
    )
    parser.add_argument("--output", required=True, help="Output markdown path")
    args = parser.parse_args()

    root = resolve_root(Path(args.root).expanduser().resolve())
    dirs = round_dirs(root)
    if not dirs:
        raise SystemExit(f"No round-* directories found under {root}")

    rows = []
    skipped = []

    for d in dirs:
        summary_path = d / "round-summary.json"
        report_path = d / "round-report.json"
        recommendation_path = d / "recommendation-report.json"

        summary_obj = None
        report_obj = None
        recommendation_obj = None

        if summary_path.exists():
            try:
                summary_obj = load_json(summary_path)
            except Exception as e:
                skipped.append((d.name, f"failed to read round-summary.json: {e}"))
                continue

        if report_path.exists():
            try:
                report_obj = load_json(report_path)
            except Exception:
                report_obj = None

        if recommendation_path.exists():
            try:
                recommendation_obj = load_json(recommendation_path)
            except Exception:
                recommendation_obj = None

        # Find best round object first, usually from round-report.json, otherwise round-summary.json
        round_obj = find_best_round_object(report_obj) if report_obj is not None else None
        if round_obj is None and summary_obj is not None:
            round_obj = find_best_round_object(summary_obj)

        # Find nested summary dict
        nested_summary = find_summary_dict(summary_obj) if summary_obj is not None else None
        if nested_summary is None and round_obj is not None and isinstance(round_obj.get("summary"), dict):
            nested_summary = round_obj["summary"]
        if nested_summary is None and report_obj is not None:
            nested_summary = find_summary_dict(report_obj)

        if nested_summary is None:
            skipped.append((d.name, "could not find summary dict"))
            continue

        round_id = (
            nested_summary.get("round_id")
            or (round_obj.get("id") if isinstance(round_obj, dict) else None)
            or d.name
        )

        detector_version = ""
        if isinstance(round_obj, dict) and isinstance(round_obj.get("detector_version"), str):
            detector_version = round_obj["detector_version"]

        promotion_count = int(nested_summary.get("promotion_count", 0) or 0)
        replay_pass_rate = float(nested_summary.get("replay_pass_rate", 0) or 0)
        robustness_outcome = str(nested_summary.get("robustness_outcome", "unknown"))

        # Find scenario results from round-report.json first, then fallback
        scenario_results = find_scenario_results_list(report_obj) if report_obj is not None else []
        if not scenario_results and isinstance(round_obj, dict) and isinstance(round_obj.get("scenario_results"), list):
            scenario_results = round_obj["scenario_results"]
        if not scenario_results and summary_obj is not None:
            scenario_results = find_scenario_results_list(summary_obj)

        by_id = {
            item.get("scenario_id"): item
            for item in scenario_results
            if isinstance(item, dict) and item.get("scenario_id")
        }

        target_case_statuses = {}
        for case in TARGET_CASES:
            item = by_id.get(case, {})
            target_case_statuses[case] = {
                "final_detection_status": item.get("final_detection_status"),
                "passed": item.get("passed"),
                "triggered_rule_ids": item.get("triggered_rule_ids"),
            }

        recommendation_types = []
        if recommendation_obj is not None:
            recs = find_recommendation_list(recommendation_obj)
            recommendation_types = [
                item["type"]
                for item in recs
                if isinstance(item, dict) and isinstance(item.get("type"), str)
            ]
        elif report_obj is not None:
            recs = find_recommendation_list(report_obj)
            recommendation_types = [
                item["type"]
                for item in recs
                if isinstance(item, dict) and isinstance(item.get("type"), str)
            ]

        rows.append(
            {
                "round_id": round_id,
                "phase": classify_phase(round_id, args.phase_b_start_round),
                "detector_version": detector_version,
                "promotion_count": promotion_count,
                "replay_pass_rate": replay_pass_rate,
                "robustness_outcome": robustness_outcome,
                "recommendation_types": recommendation_types,
                "target_case_statuses": target_case_statuses,
            }
        )

    phase_a_rows = [r for r in rows if r["phase"] == "phase_a"]
    phase_b_rows = [r for r in rows if r["phase"] == "phase_b"]

    if not phase_a_rows and not phase_b_rows:
        raise SystemExit(
            "No rounds classified into either phase. "
            "Check --phase-b-start-round and the local report files."
        )

    phase_a = summarize_rounds(phase_a_rows)
    phase_b = summarize_rounds(phase_b_rows)

    md = []
    md.append("# DRQ Week-Run Summary")
    md.append("")
    md.append("## Run segmentation")
    md.append(f"- Root used: `{root}`")
    md.append(f"- Phase B starts at round: `{args.phase_b_start_round}`")
    md.append("")
    md.append("## Parse diagnostics")
    md.append(f"- Round directories discovered: {len(dirs)}")
    md.append(f"- Parsed rounds: {len(rows)}")
    md.append(f"- Phase A rounds classified: {len(phase_a_rows)}")
    md.append(f"- Phase B rounds classified: {len(phase_b_rows)}")
    md.append(f"- Skipped round directories: {len(skipped)}")
    if skipped:
        md.append("- Skipped samples:")
        for name, reason in skipped[:5]:
            md.append(f"  - `{name}`: {reason}")
    md.append("")
    md.append("## Headline comparison")
    md.append("")
    md.append("| Metric | Phase A | Phase B |")
    md.append("|---|---:|---:|")
    md.append(f"| Rounds | {phase_a['count']} | {phase_b['count']} |")
    md.append(f"| First round | `{phase_a['first_round_id']}` | `{phase_b['first_round_id']}` |")
    md.append(f"| Last round | `{phase_a['last_round_id']}` | `{phase_b['last_round_id']}` |")
    md.append(f"| Total promotions | {phase_a['total_promotions']} | {phase_b['total_promotions']} |")
    md.append(f"| Avg promotions / round | {phase_a['avg_promotions']:.2f} | {phase_b['avg_promotions']:.2f} |")
    md.append(f"| Avg replay pass rate | {phase_a['avg_replay_pass_rate']:.2f} | {phase_b['avg_replay_pass_rate']:.2f} |")
    md.append(f"| Zero-promotion rounds | {phase_a['zero_promotion_rounds']} | {phase_b['zero_promotion_rounds']} |")
    md.append(f"| Perfect replay rounds | {phase_a['perfect_replay_rounds']} | {phase_b['perfect_replay_rounds']} |")
    md.append("")
    md.append("## Detector versions seen")
    md.append(f"- Phase A: {', '.join(phase_a['detector_versions']) if phase_a['detector_versions'] else 'None'}")
    md.append(f"- Phase B: {', '.join(phase_b['detector_versions']) if phase_b['detector_versions'] else 'None'}")
    md.append("")
    md.append("## Robustness outcomes")
    md.append(f"- Phase A: {phase_a['robustness_outcomes']}")
    md.append(f"- Phase B: {phase_b['robustness_outcomes']}")
    md.append("")
    md.append("## Targeted weak-case status in latest round of each phase")
    md.append("")
    for scenario_id in TARGET_CASES:
        a = phase_a["latest_target_case_statuses"].get(scenario_id, {})
        b = phase_b["latest_target_case_statuses"].get(scenario_id, {})
        md.append(f"### {scenario_id}")
        md.append(
            f"- Phase A latest: status=`{a.get('final_detection_status')}`, "
            f"passed=`{a.get('passed')}`, rules=`{a.get('triggered_rule_ids')}`"
        )
        md.append(
            f"- Phase B latest: status=`{b.get('final_detection_status')}`, "
            f"passed=`{b.get('passed')}`, rules=`{b.get('triggered_rule_ids')}`"
        )
        md.append("")
    md.append("## Recommendation type totals")
    md.append(f"- Phase A: {phase_a['recommendation_counts']}")
    md.append(f"- Phase B: {phase_b['recommendation_counts']}")
    md.append("")
    md.append("## Conclusion")
    if (
        phase_a["count"] > 0
        and phase_b["count"] > 0
        and phase_b["avg_promotions"] < phase_a["avg_promotions"]
        and phase_b["avg_replay_pass_rate"] > phase_a["avg_replay_pass_rate"]
    ):
        md.append("- Phase B improved detector performance versus Phase A.")
        md.append("- Promotions fell materially in the tuned window.")
        md.append("- Replay pass rate improved in the tuned window.")
    else:
        md.append("- Phase B did not clearly improve over Phase A; review round details manually.")

    output = Path(args.output).expanduser().resolve()
    output.write_text("\n".join(md) + "\n", encoding="utf-8")
    print(f"Wrote {output}")


if __name__ == "__main__":
    main()