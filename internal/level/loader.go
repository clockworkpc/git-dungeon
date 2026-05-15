package level

import (
	"fmt"
	"io/fs"
	"sort"

	"gopkg.in/yaml.v3"
)

func LoadAll(fsys fs.FS) ([]LevelDef, error) {
	entries, err := fs.ReadDir(fsys, "levels")
	if err != nil {
		return nil, fmt.Errorf("reading levels dir: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	var levels []LevelDef
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := fs.ReadFile(fsys, "levels/"+entry.Name())
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", entry.Name(), err)
		}
		var def LevelDef
		if err := yaml.Unmarshal(data, &def); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", entry.Name(), err)
		}
		levels = append(levels, def)
	}
	return levels, nil
}
