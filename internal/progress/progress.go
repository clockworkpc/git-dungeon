package progress

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Progress struct {
	CompletedLevels map[string]bool `json:"completed_levels"`
}

func dataPath() (string, error) {
	base, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, ".local", "share", "gitdungeon", "progress.json"), nil
}

func Load() (Progress, error) {
	p := Progress{CompletedLevels: make(map[string]bool)}
	path, err := dataPath()
	if err != nil {
		return p, nil
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return p, nil
	}
	if err != nil {
		return p, err
	}
	if err := json.Unmarshal(data, &p); err != nil {
		return p, err
	}
	if p.CompletedLevels == nil {
		p.CompletedLevels = make(map[string]bool)
	}
	return p, nil
}

func (p Progress) Save() error {
	path, err := dataPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (p Progress) IsCompleted(levelID string) bool {
	return p.CompletedLevels[levelID]
}

func (p *Progress) MarkCompleted(levelID string) {
	p.CompletedLevels[levelID] = true
}
