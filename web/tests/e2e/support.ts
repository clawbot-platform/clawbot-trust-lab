import type { Page, Route } from "@playwright/test";
import {
  currentRound,
  executiveSummaryReport,
  previousRound,
  promotionDetail,
  reportDescriptors,
  roundComparison,
  roundSummaryReport
} from "../../src/test/fixtures";

type ReviewState = {
  status: "accepted" | "duplicate" | "needs_follow_up" | "false_signal";
  note: string;
  updatedAt: string;
};

function fulfill(route: Route, data: unknown, status = 200) {
  return route.fulfill({
    status,
    contentType: "application/json",
    body: JSON.stringify({ data })
  });
}

export async function mockOperatorApi(page: Page) {
  const reviewState: ReviewState = {
    status: promotionDetail.review?.status ?? "needs_follow_up",
    note: promotionDetail.review?.note?.body ?? "",
    updatedAt: promotionDetail.review?.updated_at ?? "2026-03-25T12:05:00Z"
  };

  await page.route("**/api/v1/operator/**", async (route) => {
    const request = route.request();
    const url = new URL(request.url());
    const pathname = url.pathname;
    const parts = pathname.split("/").filter(Boolean);

    if (pathname === "/api/v1/operator/rounds" && request.method() === "GET") {
      return fulfill(route, [currentRound, previousRound]);
    }

    if (pathname === `/api/v1/operator/rounds/${currentRound.id}` && request.method() === "GET") {
      return fulfill(route, currentRound);
    }

    if (pathname === `/api/v1/operator/rounds/${currentRound.id}/compare` && request.method() === "GET") {
      return fulfill(route, roundComparison);
    }

    if (pathname === "/api/v1/operator/promotions" && request.method() === "GET") {
      return fulfill(route, [
        {
          round_id: promotionDetail.round_id,
          promotion: promotionDetail.promotion,
          review: {
            promotion_id: promotionDetail.promotion.id,
            status: reviewState.status,
            note: reviewState.note
              ? {
                  id: "note-live",
                  body: reviewState.note,
                  created_at: reviewState.updatedAt
                }
              : undefined,
            updated_at: reviewState.updatedAt
          }
        }
      ]);
    }

    if (pathname === `/api/v1/operator/promotions/${promotionDetail.promotion.id}` && request.method() === "GET") {
      return fulfill(route, {
        ...promotionDetail,
        review: {
          promotion_id: promotionDetail.promotion.id,
          status: reviewState.status,
          note: reviewState.note
            ? {
                id: "note-live",
                body: reviewState.note,
                created_at: reviewState.updatedAt
              }
            : undefined,
          updated_at: reviewState.updatedAt
        }
      });
    }

    if (pathname === `/api/v1/operator/promotions/${promotionDetail.promotion.id}/review` && request.method() === "POST") {
      const payload = request.postDataJSON() as { status: ReviewState["status"]; note?: string };
      reviewState.status = payload.status;
      reviewState.note = payload.note ?? "";
      reviewState.updatedAt = "2026-03-25T12:10:00Z";
      return fulfill(route, {
        promotion_id: promotionDetail.promotion.id,
        status: reviewState.status,
        note: reviewState.note
          ? {
              id: "note-live",
              body: reviewState.note,
              created_at: reviewState.updatedAt
            }
          : undefined,
        updated_at: reviewState.updatedAt
      });
    }

    if (pathname === `/api/v1/operator/reports/${currentRound.id}` && request.method() === "GET") {
      return fulfill(route, reportDescriptors);
    }

    if (parts[0] === "api" && parts[1] === "v1" && parts[2] === "operator" && parts[3] === "reports" && parts[4] === currentRound.id && request.method() === "GET") {
      const artifactName = decodeURIComponent(parts.slice(5).join("/"));
      if (artifactName === "executive-summary.md") {
        return fulfill(route, executiveSummaryReport);
      }
      if (artifactName === "round-summary.json") {
        return fulfill(route, roundSummaryReport);
      }
    }

    return route.fulfill({
      status: 404,
      contentType: "application/json",
      body: JSON.stringify({ error: { message: `unhandled mock route: ${request.method()} ${pathname}` } })
    });
  });
}
