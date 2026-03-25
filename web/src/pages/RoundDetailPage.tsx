import type { ReactNode } from "react";
import { useEffect, useMemo, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { api } from "../lib/api";
import type { BenchmarkRound, RoundComparison } from "../types/api";
import { SectionCard } from "../components/SectionCard";
import { StatusPill } from "../components/StatusPill";

export function RoundDetailPage() {
  const { roundId = "" } = useParams();
  const [round, setRound] = useState<BenchmarkRound | null>(null);
  const [rounds, setRounds] = useState<BenchmarkRound[]>([]);
  const [comparison, setComparison] = useState<RoundComparison | null>(null);
  const [selectedPrevious, setSelectedPrevious] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    Promise.all([api.getRound(roundId), api.listRounds()])
      .then(([roundData, allRounds]) => {
        setRound(roundData);
        setRounds(allRounds);
      })
      .catch((err: Error) => setError(err.message));
  }, [roundId]);

  useEffect(() => {
    if (!selectedPrevious) {
      setComparison(null);
      return;
    }
    api.compareRounds(roundId, selectedPrevious).then(setComparison).catch((err: Error) => setError(err.message));
  }, [roundId, selectedPrevious]);

  const previousOptions = useMemo(
    () => rounds.filter((item) => item.id !== roundId),
    [roundId, rounds]
  );

  if (!round) {
    return <p className="muted">{error || "Loading round..."}</p>;
  }

  return (
    <div className="stack">
      <header className="page-header">
        <div>
          <p className="eyebrow">Round Detail</p>
          <h2>{round.id}</h2>
        </div>
        <Link className="button-link" to={`/reports/${round.id}`}>
          Browse Reports
        </Link>
      </header>

      <div className="metrics-grid">
        <Metric label="Stable set" value={`${round.stable_set.passed_count}/${round.stable_set.total_count}`} />
        <Metric label="Living set" value={`${round.living_set.caught_count}/${round.living_set.total_count}`} />
        <Metric label="Replay pass rate" value={round.summary.replay_pass_rate.toFixed(2)} />
        <Metric label="Robustness" value={<StatusPill value={round.summary.robustness_outcome} />} />
      </div>

      <SectionCard title="Promoted Cases">
        {round.promotion_results.length === 0 ? <p className="muted">No promotions in this round.</p> : null}
        <div className="list-grid">
          {round.promotion_results.map((item) => (
            <article className="list-item" key={item.id}>
              <div>
                <h3>{item.scenario_id}</h3>
                <p className="muted">{item.rationale}</p>
              </div>
              <div className="actions">
                <StatusPill value={item.promotion_reason} />
                <Link className="text-link" to={`/promotions/${item.id}`}>
                  Review
                </Link>
              </div>
            </article>
          ))}
        </div>
      </SectionCard>

      <SectionCard title="Replay Regression">
        <div className="table">
          <div className="table-row table-head">
            <span>Scenario</span>
            <span>Set</span>
            <span>Status</span>
            <span>Recommendation</span>
          </div>
          {round.scenario_results
            .filter((item) => item.set_kind === "replay_regression")
            .map((item) => (
              <div className="table-row" key={item.id}>
                <span>{item.scenario_id}</span>
                <span>{item.set_kind}</span>
                <span>
                  <StatusPill value={item.final_detection_status} />
                </span>
                <span>{item.final_recommendation}</span>
              </div>
            ))}
        </div>
      </SectionCard>

      <SectionCard
        title="Round Comparison"
        action={
          <select
            aria-label="Previous round"
            value={selectedPrevious}
            onChange={(event) => setSelectedPrevious(event.target.value)}
          >
            <option value="">Select previous round</option>
            {previousOptions.map((item) => (
              <option key={item.id} value={item.id}>
                {item.id}
              </option>
            ))}
          </select>
        }
      >
        {!comparison ? <p className="muted">Select a previous round to compare.</p> : null}
        {comparison ? (
          <div className="comparison-grid">
            <Metric label="Promotion delta" value={comparison.promotions_count_delta} />
            <Metric label="Replay pass delta" value={comparison.replay_pass_rate_delta.toFixed(2)} />
            <Metric label="Challenger delta" value={comparison.challenger_count_delta} />
            <Metric label="Detection delta count" value={comparison.detection_delta_count} />
          </div>
        ) : null}
      </SectionCard>
    </div>
  );
}

function Metric({ label, value }: { label: string; value: number | string | ReactNode }) {
  return (
    <div className="metric-card">
      <p className="metric-label">{label}</p>
      <div className="metric-value">{value}</div>
    </div>
  );
}
