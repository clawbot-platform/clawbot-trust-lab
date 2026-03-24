package scenario

import "time"

type ScenarioType string

const (
	ScenarioTypeMandateReview ScenarioType = "mandate_review"
	ScenarioTypePolicyDrift   ScenarioType = "policy_drift"
	ScenarioTypeReplayCheck   ScenarioType = "replay_check"
)

type Scenario struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Type        ScenarioType `json:"type"`
	Description string       `json:"description"`
	PackID      string       `json:"pack_id"`
	CreatedAt   time.Time    `json:"created_at"`
}

type ScenarioPack struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Types       []ScenarioType `json:"types"`
	Version     string         `json:"version"`
}

func KnownTypes() []ScenarioType {
	return []ScenarioType{
		ScenarioTypeMandateReview,
		ScenarioTypePolicyDrift,
		ScenarioTypeReplayCheck,
	}
}
