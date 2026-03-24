package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"clawbot-trust-lab/internal/domain/replay"
)

type ReplayStore interface {
	Create(replay.ReplayCase) (replay.ReplayCase, error)
	List() []replay.ReplayCase
}

type FileReplayStore struct {
	mu    sync.RWMutex
	dir   string
	items map[string]replay.ReplayCase
}

func NewFileReplayStore(dir string) (*FileReplayStore, error) {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, fmt.Errorf("create replay archive dir: %w", err)
	}

	store := &FileReplayStore{
		dir:   dir,
		items: map[string]replay.ReplayCase{},
	}
	if err := store.loadExisting(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *FileReplayStore) Create(item replay.ReplayCase) (replay.ReplayCase, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fileName, err := replayFileName(item.ID)
	if err != nil {
		return replay.ReplayCase{}, err
	}

	filePath := filepath.Join(s.dir, fileName)
	raw, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return replay.ReplayCase{}, fmt.Errorf("marshal replay case: %w", err)
	}
	if err := os.WriteFile(filePath, raw, 0o600); err != nil {
		return replay.ReplayCase{}, fmt.Errorf("write replay archive: %w", err)
	}

	s.items[item.ID] = item
	return item, nil
}

func (s *FileReplayStore) List() []replay.ReplayCase {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]replay.ReplayCase, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].RecordedAt.Before(items[j].RecordedAt)
	})
	return items
}

func (s *FileReplayStore) loadExisting() error {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return fmt.Errorf("read replay archive dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		fileName, err := replayFileName(strings.TrimSuffix(entry.Name(), ".json"))
		if err != nil {
			return err
		}

		filePath := filepath.Join(s.dir, fileName)

		// #nosec G304 -- path is constrained to a validated filename within the configured replay archive directory
		raw, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("read replay archive %s: %w", fileName, err)
		}

		var item replay.ReplayCase
		if err := json.Unmarshal(raw, &item); err != nil {
			return fmt.Errorf("unmarshal replay archive %s: %w", fileName, err)
		}
		s.items[item.ID] = item
	}

	return nil
}

func replayFileName(id string) (string, error) {
	clean := filepath.Clean(strings.TrimSpace(id))
	if clean == "." || clean == "" {
		return "", fmt.Errorf("replay case id is required")
	}
	if filepath.IsAbs(clean) {
		return "", fmt.Errorf("absolute replay case ids are not allowed: %q", id)
	}
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("replay case id escapes root: %q", id)
	}
	if clean != filepath.Base(clean) {
		return "", fmt.Errorf("nested replay case ids are not allowed: %q", id)
	}

	return clean + ".json", nil
}
