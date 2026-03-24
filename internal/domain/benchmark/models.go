package benchmark

import "time"

type BenchmarkRoundRef struct {
	ID         string    `json:"id"`
	Label      string    `json:"label"`
	ScenarioID string    `json:"scenario_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type StableSuiteRef struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type LivingSuiteRef struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	MutationPolicy string `json:"mutation_policy"`
}

type RegistrationRequest struct {
	StableSuite         StableSuiteRef `json:"stable_suite"`
	LivingSuite         LivingSuiteRef `json:"living_suite"`
	ScenarioPackID      string         `json:"scenario_pack_id"`
	ScenarioPackVersion string         `json:"scenario_pack_version"`
	ReplayCaseRefs      []string       `json:"replay_case_refs"`
	Notes               string         `json:"notes"`
}

type RegistrationResult struct {
	RegistrationID string    `json:"registration_id"`
	Status         string    `json:"status"`
	RegisteredAt   time.Time `json:"registered_at"`
}
