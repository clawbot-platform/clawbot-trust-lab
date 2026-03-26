import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { api } from "../lib/api";
import type { BenchmarkRecommendation, BenchmarkRound, LongRunSummary } from "../types/api";
import { SectionCard } from "../components/SectionCard";
import { StatusPill } from "../components/StatusPill";

export function RoundsPage() {
  const [rounds, setRounds] = useState<BenchmarkRound[]>([]);
  const [summary, setSummary] = useState<LongRunSummary | null>(null);
  const [recommendations, setRecommendations] = useState<BenchmarkRecommendation[]>([]);
  const [error, setError] = useState<string>("");

  useEffect(() => {
    Promise.all([api.listRounds(), api.getTrendSummary(), api.listRecommendations()])
      .then(([roundItems, trendSummary, recommendationItems]) => {
        setRounds(roundItems);
        setSummary(trendSummary);
        setRecommendations(recommendationItems);
      })
      .catch((err: Error) => setError(err.message));
  }, []);

  return (
    <div className="stack">
      <header className="page-header">
        <div>
          <p className="eyebrow">Operator Workflow</p>
          <h2>Rounds</h2>
        </div>
        <div className="toolbar">
          <Link className="button-link" to="/promotions">
            Review Promotions
          </Link>
          <Link className="button-link" to="/recommendations">
            View Recommendations
          </Link>
        </div>
      </header>

      <SectionCard title="Benchmark Rounds">
        {error ? <p className="error-text">{error}</p> : null}
        <div className="table">
          <div className="table-row table-head">
            <span>Round</span>
            <span>Family</span>
            <span>Promotions</span>
            <span>Replay Pass Rate</span>
            <span>Outcome</span>
            <span>Actions</span>
          </div>
          {rounds.map((round) => (
            <div className="table-row" key={round.id}>
              <span>{round.id}</span>
              <span>{round.scenario_family}</span>
              <span>{round.summary.promotion_count}</span>
              <span>{round.summary.replay_pass_rate.toFixed(2)}</span>
              <span>
                <StatusPill value={round.summary.robustness_outcome} />
              </span>
              <span className="table-actions">
                <Link className="text-link" to={`/rounds/${round.id}`}>
                  Open
                </Link>
                <Link className="text-link" to={`/reports/${round.id}`}>
                  Reports
                </Link>
              </span>
            </div>
          ))}
        </div>
      </SectionCard>

      <SectionCard title="Long-Run Trend Summary">
        {!summary ? <p className="muted">Loading trend summary...</p> : null}
        {summary ? (
          <div className="metrics-grid">
            <Metric label="Rounds executed" value={summary.rounds_executed} />
            <Metric label="New blind spots" value={summary.new_blind_spots_discovered} />
            <Metric label="Regressions observed" value={summary.regressions_observed} />
            <Metric label="Recurring patterns" value={summary.top_recurring_evasion_patterns.length} />
          </div>
        ) : null}
      </SectionCard>

      <SectionCard title="Recommendation Snapshot">
        {recommendations.length === 0 ? <p className="muted">No recommendations yet.</p> : null}
        <div className="list-grid">
          {recommendations.slice(0, 3).map((item) => (
            <article className="list-item" key={item.id}>
              <div className="list-item-stack">
                <div className="item-head">
                  <h3>{item.type}</h3>
                  <StatusPill value={item.priority} />
                </div>
                <p className="muted">{item.rationale}</p>
                <p className="meta-line">Suggested action: {item.suggested_action}</p>
              </div>
              <div className="actions">
                <Link className="text-link" to={`/rounds/${item.linked_round_id}`}>
                  Round
                </Link>
              </div>
            </article>
          ))}
        </div>
      </SectionCard>
    </div>
  );
}

function Metric({ label, value }: { label: string; value: number }) {
  return (
    <div className="metric-card">
      <p className="metric-label">{label}</p>
      <div className="metric-value">{value}</div>
    </div>
  );
}
