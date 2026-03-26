package handlers

import (
	"net/http"
	"strings"

	"clawbot-trust-lab/internal/domain/benchmark"
	detectionmodel "clawbot-trust-lab/internal/domain/detection"
	operator "clawbot-trust-lab/internal/services/operator"
)

type OperatorService interface {
	ListRounds() []benchmark.BenchmarkRound
	GetRound(string) (benchmark.BenchmarkRound, error)
	CompareRounds(string, string) (benchmark.RoundComparison, error)
	ListPromotions(string) []operator.PromotionRecord
	GetPromotion(string) (operator.PromotionDetail, error)
	ReviewPromotion(string, operator.ReviewInput) (benchmark.PromotionReview, error)
	GetDetectionResult(string) (detectionmodel.DetectionResult, error)
	ListRecommendations() []benchmark.Recommendation
	GetRecommendation(string) (benchmark.Recommendation, error)
	GetTrendSummary() benchmark.LongRunSummary
	GetReports(string) ([]benchmark.ReportDescriptor, error)
	GetReportArtifact(string, string) (operator.ReportContent, error)
}

type OperatorHandler struct {
	operator OperatorService
}

func NewOperatorHandler(operatorService OperatorService) *OperatorHandler {
	return &OperatorHandler{operator: operatorService}
}

func (h *OperatorHandler) ListRounds(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"data": h.operator.ListRounds()})
}

func (h *OperatorHandler) GetRound(w http.ResponseWriter, r *http.Request) {
	round, err := h.operator.GetRound(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": round})
}

func (h *OperatorHandler) CompareRounds(w http.ResponseWriter, r *http.Request) {
	previous := strings.TrimSpace(r.URL.Query().Get("previous"))
	if previous == "" {
		writeError(w, http.StatusBadRequest, "previous is required")
		return
	}
	comparison, err := h.operator.CompareRounds(r.PathValue("id"), previous)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": comparison})
}

func (h *OperatorHandler) ListPromotions(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"data": h.operator.ListPromotions(strings.TrimSpace(r.URL.Query().Get("status")))})
}

func (h *OperatorHandler) GetPromotion(w http.ResponseWriter, r *http.Request) {
	detail, err := h.operator.GetPromotion(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": detail})
}

func (h *OperatorHandler) ReviewPromotion(w http.ResponseWriter, r *http.Request) {
	var input operator.ReviewInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	review, err := h.operator.ReviewPromotion(r.PathValue("id"), input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"data": review})
}

func (h *OperatorHandler) GetDetectionResult(w http.ResponseWriter, r *http.Request) {
	result, err := h.operator.GetDetectionResult(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": result})
}

func (h *OperatorHandler) ListRecommendations(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"data": h.operator.ListRecommendations()})
}

func (h *OperatorHandler) GetRecommendation(w http.ResponseWriter, r *http.Request) {
	item, err := h.operator.GetRecommendation(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": item})
}

func (h *OperatorHandler) GetTrendSummary(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"data": h.operator.GetTrendSummary()})
}

func (h *OperatorHandler) GetReports(w http.ResponseWriter, r *http.Request) {
	items, err := h.operator.GetReports(r.PathValue("round_id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": items})
}

func (h *OperatorHandler) GetReportArtifact(w http.ResponseWriter, r *http.Request) {
	item, err := h.operator.GetReportArtifact(r.PathValue("round_id"), r.PathValue("artifact_name"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": item})
}
