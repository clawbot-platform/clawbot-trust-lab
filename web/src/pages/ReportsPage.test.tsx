import { fireEvent, screen } from "@testing-library/react";
import { vi } from "vitest";
import { api } from "../lib/api";
import {
  currentRound,
  executiveSummaryReport,
  reportDescriptors,
  roundSummaryReport
} from "../test/fixtures";
import { renderRoute } from "../test/renderRoute";

test("reports page renders artifact content and allows switching artifacts", async () => {
  vi.spyOn(api, "getReports").mockResolvedValue(reportDescriptors);
  const artifactSpy = vi
    .spyOn(api, "getReportArtifact")
    .mockResolvedValueOnce(executiveSummaryReport)
    .mockResolvedValueOnce(roundSummaryReport);

  renderRoute(`/reports/${currentRound.id}`);

  expect(await screen.findByText("Executive Summary")).toBeInTheDocument();

  fireEvent.click(screen.getByRole("button", { name: /round-summary\.json/i }));

  expect(await screen.findByText(/"promotion_count": 1/)).toBeInTheDocument();
  expect(artifactSpy).toHaveBeenNthCalledWith(1, currentRound.id, "executive-summary.md");
  expect(artifactSpy).toHaveBeenNthCalledWith(2, currentRound.id, "round-summary.json");
});
