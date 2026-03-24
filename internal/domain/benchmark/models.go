package benchmark

type BenchmarkRoundRef struct {
	ID         string `json:"id"`
	Label      string `json:"label"`
	ScenarioID string `json:"scenario_id"`
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
