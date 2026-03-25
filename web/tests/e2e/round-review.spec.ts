import { expect, test } from "@playwright/test";
import { currentRound, previousRound } from "../../src/test/fixtures";
import { mockOperatorApi } from "./support";

test("operator can inspect a round, compare it, and open a report artifact", async ({ page }) => {
    await mockOperatorApi(page);

    await page.goto("/");
    await page.waitForLoadState("domcontentloaded");

    await expect(page.getByRole("heading", { name: "Benchmark Rounds" })).toBeVisible();
    await expect(page.getByText(currentRound.id)).toBeVisible();

    await page.getByRole("link", { name: /open/i }).first().click();

    await expect(page.getByRole("heading", { name: currentRound.id })).toBeVisible();
    await expect(page.getByText(/new blind spot discovered/i)).toBeVisible();
    await expect(page.getByText("0.67")).toBeVisible();

    await page.getByRole("combobox", { name: "Previous round" }).selectOption(previousRound.id);
    await expect(page.getByText("-0.33")).toBeVisible();

    await page.getByRole("link", { name: "Browse Reports" }).click();

    await expect(page.getByRole("heading", { name: currentRound.id })).toBeVisible();
    await expect(page.getByText("Executive Summary")).toBeVisible();
    await page.getByRole("button", { name: /round-summary\.json/i }).click();
    await expect(page.getByText(/"promotion_count": 1/)).toBeVisible();
});