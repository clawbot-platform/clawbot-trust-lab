package store

import "clawbot-trust-lab/internal/domain/scenario"

type ScenarioCatalog interface {
	ScenarioTypes() []scenario.ScenarioType
}

type InMemoryScenarioCatalog struct{}

func NewInMemoryScenarioCatalog() *InMemoryScenarioCatalog {
	return &InMemoryScenarioCatalog{}
}

func (c *InMemoryScenarioCatalog) ScenarioTypes() []scenario.ScenarioType {
	return scenario.KnownTypes()
}
