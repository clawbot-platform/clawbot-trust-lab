package scenario

import "time"

type ScenarioType string
type ScenarioSetRole string

const (
	ScenarioTypeMandateReview        ScenarioType = "mandate_review"
	ScenarioTypePolicyDrift          ScenarioType = "policy_drift"
	ScenarioTypeReplayCheck          ScenarioType = "replay_check"
	ScenarioTypeCommercePurchase     ScenarioType = "commerce_purchase"
	ScenarioTypeCommerceRefundReview ScenarioType = "commerce_refund_review"
)

const (
	ScenarioSetRoleStable ScenarioSetRole = "stable"
	ScenarioSetRoleLiving ScenarioSetRole = "living"
)

type FeatureTierModel struct {
	TierA []string `json:"tier_a"`
	TierB []string `json:"tier_b"`
	TierC []string `json:"tier_c,omitempty"`
}

type Scenario struct {
	ID               string           `json:"id"`
	Code             string           `json:"code,omitempty"`
	Name             string           `json:"name"`
	Type             ScenarioType     `json:"type"`
	Family           string           `json:"family,omitempty"`
	SetRole          ScenarioSetRole  `json:"set_role,omitempty"`
	VariantID        string           `json:"variant_id,omitempty"`
	Description      string           `json:"description"`
	PackID           string           `json:"pack_id"`
	Version          string           `json:"version"`
	Actors           []string         `json:"actors"`
	TrustSignals     []string         `json:"trust_signals"`
	ExpectedOutcomes []string         `json:"expected_outcomes"`
	Tags             []string         `json:"tags"`
	FeatureModel     FeatureTierModel `json:"feature_model"`
	CreatedAt        time.Time        `json:"created_at"`
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
		ScenarioTypeCommercePurchase,
		ScenarioTypeCommerceRefundReview,
	}
}
