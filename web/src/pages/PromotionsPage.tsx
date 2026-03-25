import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { api } from "../lib/api";
import type { PromotionRecord, PromotionReviewStatus } from "../types/api";
import { SectionCard } from "../components/SectionCard";
import { StatusPill } from "../components/StatusPill";

const filters: Array<PromotionReviewStatus | ""> = ["", "accepted", "duplicate", "needs_follow_up", "false_signal"];

export function PromotionsPage() {
  const [status, setStatus] = useState<PromotionReviewStatus | "">("");
  const [promotions, setPromotions] = useState<PromotionRecord[]>([]);
  const [error, setError] = useState("");

  useEffect(() => {
    api.listPromotions(status || undefined).then(setPromotions).catch((err: Error) => setError(err.message));
  }, [status]);

  return (
    <div className="stack">
      <header className="page-header">
        <div>
          <p className="eyebrow">Operator Workflow</p>
          <h2>Promotions</h2>
        </div>
        <select value={status} onChange={(event) => setStatus(event.target.value as PromotionReviewStatus | "")}>
          {filters.map((item) => (
            <option key={item || "all"} value={item}>
              {item || "all statuses"}
            </option>
          ))}
        </select>
      </header>

      <SectionCard title="Promoted Challenger Cases">
        {error ? <p className="error-text">{error}</p> : null}
        <div className="list-grid">
          {promotions.map((item) => (
            <article className="list-item" key={item.promotion.id}>
              <div>
                <h3>{item.promotion.scenario_id}</h3>
                <p className="muted">{item.promotion.rationale}</p>
                <p className="meta-line">Round: {item.round_id}</p>
              </div>
              <div className="actions">
                <StatusPill value={item.review?.status ?? "needs_follow_up"} />
                <Link className="text-link" to={`/promotions/${item.promotion.id}`}>
                  Open
                </Link>
              </div>
            </article>
          ))}
        </div>
      </SectionCard>
    </div>
  );
}
