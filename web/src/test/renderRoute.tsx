import type { ReactElement } from "react";
import { render } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { Layout } from "../components/Layout";
import { PromotionDetailPage } from "../pages/PromotionDetailPage";
import { PromotionsPage } from "../pages/PromotionsPage";
import { RecommendationsPage } from "../pages/RecommendationsPage";
import { ReportsPage } from "../pages/ReportsPage";
import { RoundDetailPage } from "../pages/RoundDetailPage";
import { RoundsPage } from "../pages/RoundsPage";

export function renderRoute(initialEntry: string): ReturnType<typeof render> {
  return render(
    <MemoryRouter
      future={{ v7_relativeSplatPath: true, v7_startTransition: true }}
      initialEntries={[initialEntry]}
    >
      <Routes>
        <Route element={<Layout />} path="/">
          <Route element={<RoundsPage />} index />
          <Route element={<RoundDetailPage />} path="rounds/:roundId" />
          <Route element={<PromotionsPage />} path="promotions" />
          <Route element={<PromotionDetailPage />} path="promotions/:promotionId" />
          <Route element={<RecommendationsPage />} path="recommendations" />
          <Route element={<ReportsPage />} path="reports/:roundId" />
        </Route>
      </Routes>
    </MemoryRouter>
  );
}

export function wrapWithRouter(element: ReactElement, initialEntry = "/"): ReturnType<typeof render> {
  return render(
    <MemoryRouter
      future={{ v7_relativeSplatPath: true, v7_startTransition: true }}
      initialEntries={[initialEntry]}
    >
      {element}
    </MemoryRouter>
  );
}
