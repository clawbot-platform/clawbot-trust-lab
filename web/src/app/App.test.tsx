import "@testing-library/jest-dom";
import { render, screen } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { Layout } from "../components/Layout";

test("renders operator navigation", () => {
  render(
    <MemoryRouter future={{ v7_relativeSplatPath: true, v7_startTransition: true }}>
      <Routes>
        <Route element={<Layout />} path="/" />
      </Routes>
    </MemoryRouter>
  );

  expect(screen.getByText("Trust Lab Operator")).toBeInTheDocument();
  expect(screen.getByText("Rounds")).toBeInTheDocument();
  expect(screen.getByText("Promotions")).toBeInTheDocument();
  expect(screen.getByText("Recommendations")).toBeInTheDocument();
});
