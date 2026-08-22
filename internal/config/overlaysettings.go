package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func DefaultOverlaySettings() OverlaySettings {
	return OverlaySettings{
		Enabled:        true,
		RefreshSeconds: 1.5,
	}
}

// LoadOverlaySettings reads overlay.yaml, creating it with defaults on first run.
func LoadOverlaySettings() (OverlaySettings, error) {
	path, err := OverlaySettingsPath()
	if err != nil {
		return OverlaySettings{}, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		defaults := DefaultOverlaySettings()
		if werr := SaveOverlaySettings(defaults); werr != nil {
			return defaults, werr
		}
		return defaults, nil
	}
	if err != nil {
		return OverlaySettings{}, err
	}

	var os2 OverlaySettings
	if err := yaml.Unmarshal(data, &os2); err != nil {
		return OverlaySettings{}, err
	}
	if os2.RefreshSeconds <= 0 {
		os2.RefreshSeconds = 1.5
	}
	return os2, nil
}

func SaveOverlaySettings(s OverlaySettings) error {
	path, err := OverlaySettingsPath()
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
