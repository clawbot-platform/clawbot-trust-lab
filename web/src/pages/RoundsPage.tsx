import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { api } from "../lib/api";
import type { BenchmarkRound } from "../types/api";
import { SectionCard } from "../components/SectionCard";
import { StatusPill } from "../components/StatusPill";

export function RoundsPage() {
  const [rounds, setRounds] = useState<BenchmarkRound[]>([]);
  const [error, setError] = useState<string>("");

  useEffect(() => {
    api.listRounds().then(setRounds).catch((err: Error) => setError(err.message));
  }, []);

  return (
    <div className="stack">
      <header className="page-header">
        <div>
          <p className="eyebrow">Operator Workflow</p>
          <h2>Rounds</h2>
        </div>
        <p className="muted">Inspect benchmark rounds, replay pass rate, and promoted challenger cases.</p>
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
            <span />
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
              <span>
                <Link className="text-link" to={`/rounds/${round.id}`}>
                  Open
                </Link>
              </span>
            </div>
          ))}
        </div>
      </SectionCard>
    </div>
  );
}
