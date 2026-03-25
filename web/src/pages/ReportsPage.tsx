import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { marked } from "marked";
import { api } from "../lib/api";
import type { ReportContent, ReportDescriptor } from "../types/api";
import { SectionCard } from "../components/SectionCard";

export function ReportsPage() {
  const { roundId = "" } = useParams();
  const [reports, setReports] = useState<ReportDescriptor[]>([]);
  const [active, setActive] = useState<ReportContent | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    api.getReports(roundId)
      .then((items) => {
        setReports(items);
        if (items[0]) {
          return api.getReportArtifact(roundId, items[0].artifact_name);
        }
        return null;
      })
      .then((firstItem) => {
        if (firstItem) {
          setActive(firstItem);
        }
      })
      .catch((err: Error) => setError(err.message));
  }, [roundId]);

  async function selectReport(item: ReportDescriptor) {
    try {
      const content = await api.getReportArtifact(roundId, item.artifact_name);
      setActive(content);
    } catch (err) {
      setError((err as Error).message);
    }
  }

  return (
    <div className="stack">
      <header className="page-header">
        <div>
          <p className="eyebrow">Report Browser</p>
          <h2>{roundId}</h2>
        </div>
        <Link className="button-link" to={`/rounds/${roundId}`}>
          Back To Round
        </Link>
      </header>

      {error ? <p className="error-text">{error}</p> : null}

      <div className="report-layout">
        <SectionCard title="Artifacts">
          <div className="report-list">
            {reports.map((item) => (
              <button className="report-link" key={item.artifact_name} onClick={() => selectReport(item)} type="button">
                <strong>{item.artifact_name}</strong>
                <span>{item.kind}</span>
              </button>
            ))}
          </div>
        </SectionCard>

        <SectionCard title={active?.descriptor.artifact_name ?? "Report Content"}>
          {!active ? <p className="muted">Select a report artifact.</p> : null}
          {active?.descriptor.kind === "markdown" ? (
            <article
              className="markdown-body"
              dangerouslySetInnerHTML={{ __html: marked.parse(active.content) as string }}
            />
          ) : null}
          {active?.descriptor.kind !== "markdown" && active ? <pre className="report-pre">{active.content}</pre> : null}
        </SectionCard>
      </div>
    </div>
  );
}
