import type { FormEvent, ReactNode } from "react";
import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { api } from "../lib/api";
import type { PromotionDetail, PromotionReviewStatus } from "../types/api";
import { SectionCard } from "../components/SectionCard";
import { StatusPill } from "../components/StatusPill";

const reviewStatuses: PromotionReviewStatus[] = ["accepted", "duplicate", "needs_follow_up", "false_signal"];

export function PromotionDetailPage() {
  const { promotionId = "" } = useParams();
  const [detail, setDetail] = useState<PromotionDetail | null>(null);
  const [status, setStatus] = useState<PromotionReviewStatus>("accepted");
  const [note, setNote] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    api.getPromotion(promotionId)
      .then((data) => {
        setDetail(data);
        if (data.review?.status) {
          setStatus(data.review.status);
        }
        if (data.review?.note?.body) {
          setNote(data.review.note.body);
        }
      })
      .catch((err: Error) => setError(err.message));
  }, [promotionId]);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    try {
      const review = await api.reviewPromotion(promotionId, status, note);
      setDetail((current) => (current ? { ...current, review } : current));
    } catch (err) {
      setError((err as Error).message);
    }
  }

  if (!detail) {
    return <p className="muted">{error || "Loading promotion..."}</p>;
  }

  return (
    <div className="stack">
      <header className="page-header">
        <div>
          <p className="eyebrow">Promotion Explorer</p>
          <h2>{detail.promotion.scenario_id}</h2>
        </div>
        <div className="toolbar">
          <Link className="button-link" to={`/rounds/${detail.round_id}`}>
            View Round
          </Link>
          <Link className="button-link" to={`/reports/${detail.round_id}`}>
            View Reports
          </Link>
        </div>
      </header>

      {error ? <p className="error-text">{error}</p> : null}

      <SectionCard title="Promotion Context">
        <div className="detail-grid">
          <Detail label="Round" value={detail.round_id} />
          <Detail label="Reason" value={<StatusPill value={detail.promotion.promotion_reason} />} />
          <Detail label="Detection status" value={<StatusPill value={detail.detection_result.status} />} />
          <Detail label="Recommendation" value={detail.detection_result.recommendation} />
        </div>
        <p className="narrative">{detail.promotion.rationale}</p>
      </SectionCard>

      <SectionCard title="Detection Explorer">
        <div className="chips">
          {detail.detection_result.reason_codes.map((reason) => (
            <span className="chip" key={reason}>
              {reason}
            </span>
          ))}
        </div>
        <p className="meta-line">Replay refs: {detail.detection_result.replay_case_refs.join(", ") || "none"}</p>
        <p className="meta-line">Trust refs: {detail.detection_result.trust_decision_refs.join(", ") || "none"}</p>
        <p className="meta-line">
          Tier C used: {String(((detail.detection_result.metadata?.tier_profile as { tier_c_used?: boolean } | undefined)?.tier_c_used) ?? false)}
        </p>
        {detail.scenario_result ? (
          <p className="meta-line">Memory refs: {detail.scenario_result.memory_record_refs.join(", ") || "none"}</p>
        ) : null}
      </SectionCard>

      <SectionCard title="Review Action">
        <p className="meta-line">
          Historical review state: {detail.review ? `last updated ${detail.review.updated_at}` : "no persisted operator review on this promotion yet"}
        </p>
        <form className="review-form" onSubmit={onSubmit}>
          <label>
            Status
            <select value={status} onChange={(event) => setStatus(event.target.value as PromotionReviewStatus)}>
              {reviewStatuses.map((item) => (
                <option key={item} value={item}>
                  {item}
                </option>
              ))}
            </select>
          </label>
          <label>
            Operator note
            <textarea value={note} onChange={(event) => setNote(event.target.value)} rows={5} />
          </label>
          <button className="primary-button" type="submit">
            Save Review
          </button>
        </form>
      </SectionCard>
    </div>
  );
}

function Detail({ label, value }: { label: string; value: string | number | ReactNode }) {
  return (
    <div>
      <p className="metric-label">{label}</p>
      <div className="metric-value">{value}</div>
    </div>
  );
}
