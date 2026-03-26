import type {
  BenchmarkRound,
  BenchmarkRecommendation,
  DetectionResult,
  LongRunSummary,
  PromotionDetail,
  PromotionRecord,
  PromotionReview,
  ReportContent,
  ReportDescriptor,
  RoundComparison
} from "../types/api";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, {
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers ?? {})
    },
    ...init
  });

  const payload = (await response.json()) as { data?: T; error?: { message?: string } };
  if (!response.ok) {
    throw new Error(payload.error?.message ?? `request failed for ${path}`);
  }
  if (payload.data === undefined) {
    throw new Error(`missing data for ${path}`);
  }
  return payload.data;
}

export const api = {
  listRounds: () => request<BenchmarkRound[]>("/api/v1/operator/rounds"),
  getRound: (id: string) => request<BenchmarkRound>(`/api/v1/operator/rounds/${id}`),
  compareRounds: (current: string, previous: string) =>
    request<RoundComparison>(`/api/v1/operator/rounds/${current}/compare?previous=${encodeURIComponent(previous)}`),
  listPromotions: (status?: string) =>
    request<PromotionRecord[]>(status ? `/api/v1/operator/promotions?status=${encodeURIComponent(status)}` : "/api/v1/operator/promotions"),
  getPromotion: (id: string) => request<PromotionDetail>(`/api/v1/operator/promotions/${id}`),
  reviewPromotion: (id: string, status: string, note: string) =>
    request<PromotionReview>(`/api/v1/operator/promotions/${id}/review`, {
      method: "POST",
      body: JSON.stringify({ status, note })
    }),
  getDetectionResult: (id: string) => request<DetectionResult>(`/api/v1/operator/detection/results/${id}`),
  listRecommendations: () => request<BenchmarkRecommendation[]>("/api/v1/operator/recommendations"),
  getTrendSummary: () => request<LongRunSummary>("/api/v1/operator/trends/summary"),
  getReports: (roundId: string) => request<ReportDescriptor[]>(`/api/v1/operator/reports/${roundId}`),
  getReportArtifact: (roundId: string, artifactName: string) =>
    request<ReportContent>(`/api/v1/operator/reports/${roundId}/${encodeURIComponent(artifactName)}`)
};
