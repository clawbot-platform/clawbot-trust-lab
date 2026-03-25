import { screen } from "@testing-library/react";
import { vi } from "vitest";
import { api } from "../lib/api";
import { currentRound, previousRound, rounds } from "../test/fixtures";
import { renderRoute } from "../test/renderRoute";

test("rounds page renders benchmark rounds with real-looking summary data", async () => {
  vi.spyOn(api, "listRounds").mockResolvedValue(rounds);

  renderRoute("/");

  expect(await screen.findByText(currentRound.id)).toBeInTheDocument();
  expect(screen.getByText(previousRound.id)).toBeInTheDocument();
  expect(screen.getByText("new blind spot discovered")).toBeInTheDocument();
  expect(screen.getByText("0.67")).toBeInTheDocument();
  expect(screen.getAllByRole("link", { name: "Open" })).toHaveLength(2);
});
