import { fireEvent, screen } from "@testing-library/react";
import { vi } from "vitest";
import { api } from "../lib/api";
import { promotionRecords } from "../test/fixtures";
import { renderRoute } from "../test/renderRoute";

test("promotions page renders all promotions with paging and historical review visibility", async () => {
  vi.spyOn(api, "listPromotions").mockResolvedValue(promotionRecords);

  renderRoute("/promotions");

  expect(await screen.findByText(promotionRecords[0].promotion.scenario_id)).toBeInTheDocument();
  expect(screen.getByText(promotionRecords[0].promotion.rationale)).toBeInTheDocument();
  expect(screen.getAllByText(`Round: ${promotionRecords[0].round_id}`)).toHaveLength(3);
  expect(screen.getByText("needs follow up")).toBeInTheDocument();
  expect(screen.getByText("Page 1 of 2")).toBeInTheDocument();
  expect(screen.queryByText(promotionRecords[3].promotion.scenario_id)).not.toBeInTheDocument();

  fireEvent.click(screen.getByRole("button", { name: "Next" }));

  expect(screen.getByText("Page 2 of 2")).toBeInTheDocument();
  expect(screen.getByText(promotionRecords[3].promotion.scenario_id)).toBeInTheDocument();
  expect(screen.getByText(`Round: ${promotionRecords[3].round_id}`)).toBeInTheDocument();
  expect(screen.getByRole("option", { name: "duplicate" })).toBeInTheDocument();
  expect(screen.getByText(/Review state: current operator review available/)).toBeInTheDocument();
});
