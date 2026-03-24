package scenario

import "time"

type ScenarioType string

const (
	ScenarioTypeMandateReview ScenarioType = "mandate_review"
	ScenarioTypePolicyDrift   ScenarioType = "policy_drift"
	ScenarioTypeReplayCheck   ScenarioType = "replay_check"
)

type Scenario struct {
	ID               string       `json:"id"`
	Name             string       `json:"name"`
	Type             ScenarioType `json:"type"`
	Description      string       `json:"description"`
	PackID           string       `json:"pack_id"`
	Version          string       `json:"version"`
	Actors           []string     `json:"actors"`
	TrustSignals     []string     `json:"trust_signals"`
	ExpectedOutcomes []string     `json:"expected_outcomes"`
	Tags             []string     `json:"tags"`
	CreatedAt        time.Time    `json:"created_at"`
}

type ScenarioPack struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Version     string         `json:"version"`
	Types       []ScenarioType `json:"types"`
	Scenarios   []Scenario     `json:"scenarios"`
}

func KnownTypes() []ScenarioType {
	return []ScenarioType{
		ScenarioTypeMandateReview,
		ScenarioTypePolicyDrift,
		ScenarioTypeReplayCheck,
	}
}
