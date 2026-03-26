package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"clawbot-trust-lab/internal/domain/scenario"
)

type Loader struct {
	dir string
}

func New(dir string) *Loader {
	return &Loader{dir: dir}
}

func (l *Loader) LoadAll() ([]scenario.ScenarioPack, error) {
	entries, err := os.ReadDir(l.dir)
	if err != nil {
		return nil, fmt.Errorf("read scenario pack dir %s: %w", l.dir, err)
	}

	packs := make([]scenario.ScenarioPack, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		pack, err := l.LoadFile(entry.Name())
		if err != nil {
			return nil, err
		}
		packs = append(packs, pack)
	}

	sort.Slice(packs, func(i, j int) bool {
		return packs[i].ID < packs[j].ID
	})

	return packs, nil
}

func (l *Loader) LoadFile(name string) (scenario.ScenarioPack, error) {
	clean, err := safeScenarioPackName(name)
	if err != nil {
		return scenario.ScenarioPack{}, err
	}

	fullPath := filepath.Join(l.dir, clean)

	// #nosec G304 -- path is constrained to a validated filename within the configured scenario-pack directory
	raw, err := os.ReadFile(fullPath)
	if err != nil {
		return scenario.ScenarioPack{}, fmt.Errorf("read scenario pack %s: %w", clean, err)
	}

	var file scenarioPackFile
	if err := json.Unmarshal(raw, &file); err != nil {
		return scenario.ScenarioPack{}, fmt.Errorf("unmarshal scenario pack %s: %w", clean, err)
	}

	if strings.TrimSpace(file.ID) == "" {
		return scenario.ScenarioPack{}, fmt.Errorf("scenario pack %s is missing id", clean)
	}
	if strings.TrimSpace(file.Name) == "" {
		return scenario.ScenarioPack{}, fmt.Errorf("scenario pack %s is missing name", clean)
	}
	if strings.TrimSpace(file.Version) == "" {
		return scenario.ScenarioPack{}, fmt.Errorf("scenario pack %s is missing version", clean)
	}
	if len(file.Scenarios) == 0 {
		return scenario.ScenarioPack{}, fmt.Errorf("scenario pack %s must contain at least one scenario", clean)
	}

	pack := scenario.ScenarioPack{
		ID:          file.ID,
		Name:        file.Name,
		Description: file.Description,
		Version:     file.Version,
	}

	typeSet := map[scenario.ScenarioType]struct{}{}
	now := time.Now().UTC()

	for _, item := range file.Scenarios {
		if strings.TrimSpace(item.ID) == "" || strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.Version) == "" {
			return scenario.ScenarioPack{}, fmt.Errorf("scenario pack %s contains an invalid scenario entry", clean)
		}
		if strings.TrimSpace(item.ScenarioType) == "" {
			return scenario.ScenarioPack{}, fmt.Errorf("scenario %s is missing scenario_type", item.ID)
		}

		entry := scenario.Scenario{
			ID:               item.ID,
			Code:             strings.TrimSpace(item.Code),
			Name:             item.Name,
			Type:             scenario.ScenarioType(item.ScenarioType),
			Family:           strings.TrimSpace(item.Family),
			SetRole:          scenario.ScenarioSetRole(strings.TrimSpace(item.SetRole)),
			VariantID:        strings.TrimSpace(item.VariantID),
			Description:      item.Description,
			PackID:           file.ID,
			Version:          item.Version,
			Actors:           append([]string(nil), item.Actors...),
			TrustSignals:     append([]string(nil), item.TrustSignals...),
			ExpectedOutcomes: append([]string(nil), item.ExpectedOutcomes...),
			Tags:             append([]string(nil), item.Tags...),
			FeatureModel: scenario.FeatureTierModel{
				TierA: append([]string(nil), item.FeatureModel.TierA...),
				TierB: append([]string(nil), item.FeatureModel.TierB...),
				TierC: append([]string(nil), item.FeatureModel.TierC...),
			},
			CreatedAt: now,
		}

		pack.Scenarios = append(pack.Scenarios, entry)
		typeSet[entry.Type] = struct{}{}
	}

	for kind := range typeSet {
		pack.Types = append(pack.Types, kind)
	}
	sort.Slice(pack.Types, func(i, j int) bool {
		return pack.Types[i] < pack.Types[j]
	})

	return pack, nil
}

func safeScenarioPackName(name string) (string, error) {
	clean := filepath.Clean(strings.TrimSpace(name))
	if clean == "." || clean == "" {
		return "", fmt.Errorf("scenario pack name is required")
	}
	if filepath.IsAbs(clean) {
		return "", fmt.Errorf("absolute scenario pack paths are not allowed: %q", name)
	}
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("scenario pack path escapes root: %q", name)
	}
	if clean != filepath.Base(clean) {
		return "", fmt.Errorf("nested scenario pack paths are not allowed: %q", name)
	}
	if filepath.Ext(clean) != ".json" {
		return "", fmt.Errorf("scenario pack must be a .json file: %q", name)
	}
	return clean, nil
}

type scenarioPackFile struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Version     string             `json:"version"`
	Description string             `json:"description"`
	Scenarios   []scenarioFileItem `json:"scenarios"`
}

type scenarioFileItem struct {
	ID               string   `json:"id"`
	Code             string   `json:"code"`
	Name             string   `json:"name"`
	Version          string   `json:"version"`
	ScenarioType     string   `json:"scenario_type"`
	Family           string   `json:"family"`
	SetRole          string   `json:"set_role"`
	VariantID        string   `json:"variant_id"`
	Description      string   `json:"description"`
	Actors           []string `json:"actors"`
	TrustSignals     []string `json:"trust_signals"`
	ExpectedOutcomes []string `json:"expected_outcomes"`
	Tags             []string `json:"tags"`
	FeatureModel     struct {
		TierA []string `json:"tier_a"`
		TierB []string `json:"tier_b"`
		TierC []string `json:"tier_c"`
	} `json:"feature_model"`
}
