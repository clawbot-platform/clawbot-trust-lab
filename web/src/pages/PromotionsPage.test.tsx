import { screen } from "@testing-library/react";
import { vi } from "vitest";
import { api } from "../lib/api";
import { promotionRecord } from "../test/fixtures";
import { renderRoute } from "../test/renderRoute";

test("promotions page renders promoted challenger cases", async () => {
  vi.spyOn(api, "listPromotions").mockResolvedValue([promotionRecord]);

  renderRoute("/promotions");

  expect(await screen.findByText(promotionRecord.promotion.scenario_id)).toBeInTheDocument();
  expect(screen.getByText(promotionRecord.promotion.rationale)).toBeInTheDocument();
  expect(screen.getByText(`Round: ${promotionRecord.round_id}`)).toBeInTheDocument();
  expect(screen.getByText("needs_follow_up")).toBeInTheDocument();
});
