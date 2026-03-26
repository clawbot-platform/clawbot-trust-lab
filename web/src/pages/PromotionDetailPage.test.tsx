import { fireEvent, screen, waitFor } from "@testing-library/react";
import { vi } from "vitest";
import { api } from "../lib/api";
import { promotionDetail, updatedPromotionReview } from "../test/fixtures";
import { renderRoute } from "../test/renderRoute";

test("promotion detail renders linked context and submits a review action", async () => {
  vi.spyOn(api, "getPromotion").mockResolvedValue(promotionDetail);
  const reviewSpy = vi.spyOn(api, "reviewPromotion").mockResolvedValue(updatedPromotionReview);

  renderRoute(`/promotions/${promotionDetail.promotion.id}`);

  expect(await screen.findByRole("heading", { name: promotionDetail.promotion.scenario_id })).toBeInTheDocument();
  expect(screen.getByText(promotionDetail.promotion.rationale)).toBeInTheDocument();
  expect(screen.getByText("delegated_actor_present")).toBeInTheDocument();
  expect(screen.getByText("Memory refs: mem-trust-1, mem-replay-1")).toBeInTheDocument();
  expect(screen.getByText("Tier C used: true")).toBeInTheDocument();
  expect(screen.getByRole("link", { name: "View Round" })).toHaveAttribute("href", `/rounds/${promotionDetail.round_id}`);
  expect(screen.getByText("Historical review state: last updated 2026-03-25T12:05:00Z")).toBeInTheDocument();

  fireEvent.change(screen.getByLabelText("Status"), { target: { value: "accepted" } });
  fireEvent.change(screen.getByLabelText("Operator note"), {
    target: { value: updatedPromotionReview.note?.body ?? "" }
  });
  fireEvent.click(screen.getByRole("button", { name: "Save Review" }));

  await waitFor(() => {
    expect(reviewSpy).toHaveBeenCalledWith(
      promotionDetail.promotion.id,
      "accepted",
      updatedPromotionReview.note?.body ?? ""
    );
  });
});
