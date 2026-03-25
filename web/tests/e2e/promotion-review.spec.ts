import { expect, test } from "@playwright/test";
import { promotionDetail } from "../../src/test/fixtures";
import { mockOperatorApi } from "./support";

test("operator can inspect a promotion and save a review action", async ({ page }) => {
    await mockOperatorApi(page);

    await page.goto("/promotions");
    await page.waitForLoadState("domcontentloaded");

    await expect(page.getByRole("heading", { name: "Promotions", exact: true })).toBeVisible();
    await expect(page.getByText(promotionDetail.promotion.scenario_id)).toBeVisible();

    await page.getByRole("link", { name: "Open" }).click();

    await expect(page.getByRole("heading", { name: promotionDetail.promotion.scenario_id })).toBeVisible();
    await expect(page.getByText(promotionDetail.promotion.rationale)).toBeVisible();
    await expect(page.getByText("delegated_actor_present")).toBeVisible();

    await page.getByLabel("Status").selectOption("accepted");
    await page.getByLabel("Operator note").fill("Promote this challenger into replay coverage.");
    await page.getByRole("button", { name: "Save Review" }).click();

    await page.goto("/promotions");
    await expect(page.locator("span.pill").filter({ hasText: "accepted" })).toBeVisible();
});