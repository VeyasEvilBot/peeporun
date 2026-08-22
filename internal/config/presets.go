package config

import (
	"os"

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

// LoadPresets reads presets.yaml, creating it with the DS1/DS2/DS3 defaults on first run.
func LoadPresets() ([]Preset, error) {
	path, err := PresetsPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		defaults := DefaultPresets()
		if werr := SavePresets(defaults); werr != nil {
			return defaults, werr
		}
		return defaults, nil
	}
	if err != nil {
		return nil, err
	}

	var pf PresetsFile
	if err := yaml.Unmarshal(data, &pf); err != nil {
		return nil, err
	}
	if len(pf.Presets) == 0 {
		return DefaultPresets(), nil
	}
	return pf.Presets, nil
}

func SavePresets(presets []Preset) error {
	path, err := PresetsPath()
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(PresetsFile{Presets: presets})
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
