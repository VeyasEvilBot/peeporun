package config

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

func DefaultPresets() []Preset {
	return []Preset{
		{
			ID:       "ds1-any",
			Game:     "Dark Souls",
			Category: "Any%",
			Splits: []string{
				"Asylum",
				"Gargoyles",
				"Quelaag",
				"Iron Golem",
				"Ornstein & Smough",
				"Pinwheel",
				"Sif",
				"Seath",
				"Nito",
				"Bed of Chaos",
				"Four Kings",
				"Gwyn",
			},
		},
		{
			ID:       "ds2-any-shulva",
			Game:     "Dark Souls II",
			Category: "Any% (Shulva)",
			Splits: []string{
				"Dragonrider",
				"Last Giant",
				"Pursuer",
				"Rotten x2",
				"Shulva",
				"Rotten x2",
				"Dragonriders",
				"Mirror Knight",
				"Demon of Song",
				"Velstadt",
				"Guardian Dragon",
				"Giant Lord",
				"Throne Watchers",
				"Nashandra",
			},
		},
		{
			ID:       "ds3-any",
			Game:     "Dark Souls III",
			Category: "Any%",
			Splits: []string{
				"Gundyr",
				"Vordt",
				"Crystal Sage",
				"Abyss Watchers",
				"Wolnir",
				"Dancer",
				"Deacons",
				"Pontiff",
				"Aldrich",
				"Yhorm",
				"Dragonslayer Armor",
				"Twin Princes",
				"Soul of Cinder",
			},
		},
	}
}

// presetFile is the on-disk shape of a single preset file. It deliberately
// has no ID field - a preset's ID is always just its filename (without
// .yaml), so sharing a preset is literally just sharing this one file and
// dropping it into someone else's presets/ folder, with no internal ID to
// worry about colliding or keeping in sync with the filename.
type presetFile struct {
	Game     string   `yaml:"game"`
	Category string   `yaml:"category"`
	Splits   []string `yaml:"splits"`
}

func presetFilePath(dir, id string) string {
	return filepath.Join(dir, sanitizeFilename(id)+".yaml")
}

// sanitizeFilename keeps preset IDs safe to use as filenames - guards
// against path separators or other characters that would escape the
// presets directory.
func sanitizeFilename(id string) string {
	r := strings.NewReplacer("/", "-", "\\", "-", "..", "-")
	return r.Replace(id)
}

// LoadPresets loads every preset from ~/.config/peeporun/presets/*.yaml,
// one preset per file, with the preset's ID always taken from the
// filename. On first run it migrates an old single-file presets.yaml if
// one exists (preserving each preset's original ID, so existing save
// data / PBs still line up correctly), or falls back to the built-in
// DS1/DS2/DS3 defaults if nothing exists yet at all.
func LoadPresets() ([]Preset, error) {
	dir, err := PresetsDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []os.DirEntry
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".yaml") {
			files = append(files, e)
		}
	}

	if len(files) == 0 {
		if migrated, err := migrateLegacyPresets(dir); err != nil {
			return nil, err
		} else if migrated != nil {
			return migrated, nil
		}
		// Truly fresh install - write the defaults as individual files.
		defaults := DefaultPresets()
		if err := SavePresets(defaults); err != nil {
			return defaults, err
		}
		return defaults, nil
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })

	presets := make([]Preset, 0, len(files))
	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, err
		}
		var pf presetFile
		if err := yaml.Unmarshal(data, &pf); err != nil {
			return nil, err
		}
		id := strings.TrimSuffix(f.Name(), ".yaml")
		presets = append(presets, Preset{
			ID:       id,
			Game:     pf.Game,
			Category: pf.Category,
			Splits:   pf.Splits,
		})
	}
	return presets, nil
}

// migrateLegacyPresets checks for an old single-file presets.yaml and, if
// found, splits it into individual files in dir (using each preset's
// existing ID as the filename, so save.yaml's per-preset progress/PB data
// still matches up). The old file is renamed to presets.yaml.bak rather
// than deleted, as a safety net. Returns nil (with no error) if there was
// nothing to migrate.
func migrateLegacyPresets(dir string) ([]Preset, error) {
	legacyPath, err := LegacyPresetsPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(legacyPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var pf PresetsFile
	if err := yaml.Unmarshal(data, &pf); err != nil {
		return nil, err
	}
	if len(pf.Presets) == 0 {
		return nil, nil
	}

	for _, p := range pf.Presets {
		if err := writePresetFile(dir, p); err != nil {
			return nil, err
		}
	}

	// Keep the old file as a backup rather than silently deleting it.
	os.Rename(legacyPath, legacyPath+".bak")

	return pf.Presets, nil
}

func writePresetFile(dir string, p Preset) error {
	data, err := yaml.Marshal(presetFile{
		Game:     p.Game,
		Category: p.Category,
		Splits:   p.Splits,
	})
	if err != nil {
		return err
	}
	return os.WriteFile(presetFilePath(dir, p.ID), data, 0o644)
}

// SavePresets writes every preset to its own file in the presets
// directory, and removes any leftover files for presets that no longer
// exist (e.g. after a deletion in the TUI).
func SavePresets(presets []Preset) error {
	dir, err := PresetsDir()
	if err != nil {
		return err
	}

	keep := make(map[string]bool, len(presets))
	for _, p := range presets {
		if err := writePresetFile(dir, p); err != nil {
			return err
		}
		keep[sanitizeFilename(p.ID)+".yaml"] = true
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		if !keep[e.Name()] {
			os.Remove(filepath.Join(dir, e.Name()))
		}
	}
	return nil
}
