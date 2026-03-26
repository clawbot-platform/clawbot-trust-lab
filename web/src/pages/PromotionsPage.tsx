import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { api } from "../lib/api";
import type { PromotionRecord, PromotionReviewStatus } from "../types/api";
import { SectionCard } from "../components/SectionCard";
import { StatusPill } from "../components/StatusPill";

const filters: Array<PromotionReviewStatus | ""> = ["", "accepted", "duplicate", "needs_follow_up", "false_signal"];
const pageSize = 3;

export function PromotionsPage() {
  const [status, setStatus] = useState<PromotionReviewStatus | "">("");
  const [promotions, setPromotions] = useState<PromotionRecord[]>([]);
  const [page, setPage] = useState(1);
  const [error, setError] = useState("");

  useEffect(() => {
    api.listPromotions(status || undefined).then(setPromotions).catch((err: Error) => setError(err.message));
  }, [status]);

  useEffect(() => {
    setPage(1);
  }, [status]);

  const totalPages = Math.max(1, Math.ceil(promotions.length / pageSize));
  const visiblePromotions = promotions.slice((page - 1) * pageSize, page * pageSize);

  return (
    <div className="stack">
      <header className="page-header">
        <div>
          <p className="eyebrow">Operator Workflow</p>
          <h2>Promotions</h2>
        </div>
        <div className="toolbar">
          <Link className="button-link" to="/recommendations">
            View Recommendations
          </Link>
          <select value={status} onChange={(event) => setStatus(event.target.value as PromotionReviewStatus | "")}>
            {filters.map((item) => (
              <option key={item || "all"} value={item}>
                {item || "all statuses"}
              </option>
            ))}
          </select>
        </div>
      </header>

      <SectionCard title="Promoted Challenger Cases">
        {error ? <p className="error-text">{error}</p> : null}
        <div className="list-grid">
          {visiblePromotions.map((item) => (
            <article className="list-item" key={item.promotion.id}>
              <div className="list-item-stack">
                <div className="item-head">
                  <h3>{item.promotion.scenario_id}</h3>
                  <div className="actions">
                    <StatusPill value={item.promotion.promotion_reason} />
                    <StatusPill value={item.review?.status ?? "unreviewed"} />
                  </div>
                </div>
                <p className="muted">{item.promotion.rationale}</p>
                <p className="meta-line">Round: {item.round_id}</p>
                <p className="meta-line">
                  Review state: {item.review ? "current operator review available" : "historical promotion with no persisted review"}
                </p>
              </div>
              <div className="actions">
                <Link className="text-link" to={`/promotions/${item.promotion.id}`}>
                  Open
                </Link>
              </div>
            </article>
          ))}
        </div>
        <div className="pager">
          <button className="pager-button" disabled={page === 1} onClick={() => setPage((current) => Math.max(1, current - 1))} type="button">
            Previous
          </button>
          <span className="meta-line">
            Page {page} of {totalPages}
          </span>
          <button
            className="pager-button"
            disabled={page === totalPages}
            onClick={() => setPage((current) => Math.min(totalPages, current + 1))}
            type="button"
          >
            Next
          </button>
        </div>
      </SectionCard>
    </div>
  );
}
