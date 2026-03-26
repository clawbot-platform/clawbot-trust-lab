import { screen } from "@testing-library/react";
import { vi } from "vitest";
import { api } from "../lib/api";
import { currentRound, longRunSummary, previousRound, recommendations, rounds } from "../test/fixtures";
import { renderRoute } from "../test/renderRoute";

test("rounds page renders benchmark rounds with real-looking summary data", async () => {
  vi.spyOn(api, "listRounds").mockResolvedValue(rounds);
  vi.spyOn(api, "getTrendSummary").mockResolvedValue(longRunSummary);
  vi.spyOn(api, "listRecommendations").mockResolvedValue(recommendations);

  renderRoute("/");

  expect(await screen.findByText(currentRound.id)).toBeInTheDocument();
  expect(screen.getByText(previousRound.id)).toBeInTheDocument();
  expect(screen.getByText("new blind spot discovered")).toBeInTheDocument();
  expect(screen.getByText("0.67")).toBeInTheDocument();
  expect(screen.getByText("Rounds executed")).toBeInTheDocument();
  expect(screen.getAllByRole("link", { name: "Open" })).toHaveLength(2);
  expect(screen.getByRole("heading", { name: "Recommendation Snapshot" })).toBeInTheDocument();
  expect(screen.getByText(recommendations[0].type)).toBeInTheDocument();
  expect(screen.getByRole("link", { name: "View Recommendations" })).toBeInTheDocument();
});
