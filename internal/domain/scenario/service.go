package scenario

import "fmt"

type PackSource interface {
	LoadAll() ([]ScenarioPack, error)
}

type Service struct {
	packs        map[string]ScenarioPack
	ordered      []ScenarioPack
	scenarioByID map[string]Scenario
}

func NewService(source PackSource) (*Service, error) {
	packs, err := source.LoadAll()
	if err != nil {
		return nil, err
	}

	service := &Service{
		packs:        map[string]ScenarioPack{},
		scenarioByID: map[string]Scenario{},
	}
	for _, pack := range packs {
		service.packs[pack.ID] = pack
		service.ordered = append(service.ordered, pack)
		for _, scenario := range pack.Scenarios {
			service.scenarioByID[scenario.ID] = scenario
		}
	}

	return service, nil
}

func (s *Service) ListPacks() []ScenarioPack {
	items := make([]ScenarioPack, len(s.ordered))
	copy(items, s.ordered)
	return items
}

func (s *Service) GetPack(id string) (ScenarioPack, error) {
	pack, ok := s.packs[id]
	if !ok {
		return ScenarioPack{}, fmt.Errorf("scenario pack %s not found", id)
	}
	return pack, nil
}

func (s *Service) GetScenario(id string) (Scenario, error) {
	item, ok := s.scenarioByID[id]
	if !ok {
		return Scenario{}, fmt.Errorf("scenario %s not found", id)
	}
	return item, nil
}
