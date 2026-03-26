import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import { api } from "../lib/api";
import type { BenchmarkRecommendation } from "../types/api";
import { SectionCard } from "../components/SectionCard";
import { StatusPill } from "../components/StatusPill";

const priorityFilters = ["all", "high", "medium", "low"] as const;

export function RecommendationsPage() {
  const [recommendations, setRecommendations] = useState<BenchmarkRecommendation[]>([]);
  const [priority, setPriority] = useState<(typeof priorityFilters)[number]>("all");
  const [error, setError] = useState("");

  useEffect(() => {
    api.listRecommendations().then(setRecommendations).catch((err: Error) => setError(err.message));
  }, []);

  const filtered = useMemo(() => {
    if (priority === "all") {
      return recommendations;
    }
    return recommendations.filter((item) => item.priority === priority);
  }, [priority, recommendations]);

  return (
    <div className="stack">
      <header className="page-header">
        <div>
          <p className="eyebrow">Operator Workflow</p>
          <h2>Recommendations</h2>
        </div>
        <div className="toolbar">
          <Link className="button-link" to="/">
            View Rounds
          </Link>
          <select aria-label="Recommendation priority" value={priority} onChange={(event) => setPriority(event.target.value as (typeof priorityFilters)[number])}>
            {priorityFilters.map((item) => (
              <option key={item} value={item}>
                {item === "all" ? "all priorities" : item}
              </option>
            ))}
          </select>
        </div>
      </header>

      <SectionCard title="Fraud-Team Guidance">
        {error ? <p className="error-text">{error}</p> : null}
        {filtered.length === 0 ? <p className="muted">No recommendations were generated yet.</p> : null}
        <div className="list-grid">
          {filtered.map((item) => (
            <article className="list-item" key={item.id}>
              <div className="list-item-stack">
                <div className="item-head">
                  <h3>{item.type}</h3>
                  <div className="actions">
                    <StatusPill value={item.priority} />
                  </div>
                </div>
                <p className="muted">{item.rationale}</p>
                <p className="meta-line">Suggested action: {item.suggested_action}</p>
                <p className="meta-line">Round: {item.linked_round_id}</p>
                <p className="meta-line">Scenarios: {item.linked_scenario_ids.join(", ")}</p>
                {item.existing_control_integration_note ? (
                  <p className="meta-line">Sidecar note: {item.existing_control_integration_note}</p>
                ) : null}
                {item.supporting_rule_ids?.length ? (
                  <div className="chips">
                    {item.supporting_rule_ids.map((rule) => (
                      <span className="chip" key={rule}>
                        {rule}
                      </span>
                    ))}
                  </div>
                ) : null}
              </div>
              <div className="actions">
                <Link className="text-link" to={`/rounds/${item.linked_round_id}`}>
                  View Round
                </Link>
              </div>
            </article>
          ))}
        </div>
      </SectionCard>
    </div>
  );
}
