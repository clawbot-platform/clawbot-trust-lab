import { BrowserRouter, Route, Routes } from "react-router-dom";
import { Layout } from "../components/Layout";
import { PromotionDetailPage } from "../pages/PromotionDetailPage";
import { PromotionsPage } from "../pages/PromotionsPage";
import { ReportsPage } from "../pages/ReportsPage";
import { RoundDetailPage } from "../pages/RoundDetailPage";
import { RoundsPage } from "../pages/RoundsPage";

export function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<Layout />} path="/">
          <Route element={<RoundsPage />} index />
          <Route element={<RoundDetailPage />} path="rounds/:roundId" />
          <Route element={<PromotionsPage />} path="promotions" />
          <Route element={<PromotionDetailPage />} path="promotions/:promotionId" />
          <Route element={<ReportsPage />} path="reports/:roundId" />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}
