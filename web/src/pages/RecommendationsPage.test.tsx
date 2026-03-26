import { screen } from "@testing-library/react";
import { vi } from "vitest";
import { api } from "../lib/api";
import { recommendations } from "../test/fixtures";
import { renderRoute } from "../test/renderRoute";

test("recommendations page renders actionable round guidance", async () => {
  vi.spyOn(api, "listRecommendations").mockResolvedValue(recommendations);

  renderRoute("/recommendations");

  expect(await screen.findByRole("heading", { name: "Recommendations" })).toBeInTheDocument();
  expect(screen.getByText(recommendations[0].type)).toBeInTheDocument();
  expect(screen.getByText(`Suggested action: ${recommendations[0].suggested_action}`)).toBeInTheDocument();
  expect(screen.getByText(`Sidecar note: ${recommendations[0].existing_control_integration_note}`)).toBeInTheDocument();
  expect(screen.getByText("missing_provenance_sensitive_action")).toBeInTheDocument();
});
