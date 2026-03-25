import { fireEvent, screen, waitFor } from "@testing-library/react";
import { vi } from "vitest";
import { api } from "../lib/api";
import { currentRound, previousRound, roundComparison, rounds } from "../test/fixtures";
import { renderRoute } from "../test/renderRoute";

test("round detail renders summary metrics and comparison deltas", async () => {
  vi.spyOn(api, "getRound").mockResolvedValue(currentRound);
  vi.spyOn(api, "listRounds").mockResolvedValue(rounds);
  const compareSpy = vi.spyOn(api, "compareRounds").mockResolvedValue(roundComparison);

  renderRoute(`/rounds/${currentRound.id}`);

  expect(await screen.findByRole("heading", { name: currentRound.id })).toBeInTheDocument();
  expect(screen.getByText("2/2")).toBeInTheDocument();
  expect(screen.getByText("2/3")).toBeInTheDocument();
  expect(screen.getByText("0.67")).toBeInTheDocument();
  expect(screen.getByText("commerce-suspicious-refund-attempt")).toBeInTheDocument();

  fireEvent.change(screen.getByRole("combobox", { name: "Previous round" }), {
    target: { value: previousRound.id }
  });

  await waitFor(() => {
    expect(compareSpy).toHaveBeenCalledWith(currentRound.id, previousRound.id);
  });
  expect(await screen.findByText("-0.33")).toBeInTheDocument();
  expect(screen.getByText("Detection delta count")).toBeInTheDocument();
});
