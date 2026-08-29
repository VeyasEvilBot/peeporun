package config

import (
	"os"
	"path/filepath"
)

// Dir returns ~/.config/peeporun (Linux/macOS) or %AppData%/peeporun (Windows),
// creating it if it doesn't exist yet.
func Dir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		home, herr := os.UserHomeDir()
		if herr != nil {
			return "", err
		}
		base = filepath.Join(home, ".config")
	}
	dir := filepath.Join(base, "peeporun")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

func KeysPath() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "keys.yaml"), nil
}

// PresetsDir returns ~/.config/peeporun/presets/ - one .yaml file per
// preset, so a single preset can be shared/imported just by handing
// someone that one file.
func PresetsDir() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	presetsDir := filepath.Join(dir, "presets")
	if err := os.MkdirAll(presetsDir, 0o755); err != nil {
		return "", err
	}
	return presetsDir, nil
}

// LegacyPresetsPath is the old single-file presets.yaml location, kept
// around only so LoadPresets can detect and migrate it on first run
// after upgrading to the folder-per-preset layout.
func LegacyPresetsPath() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "presets.yaml"), nil
}

func SavePath() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "save.yaml"), nil
}

func OverlaySettingsPath() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "overlay.yaml"), nil
}

func OverlayHTMLPath() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "overlay.html"), nil
}

func ThemePath() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "theme.yaml"), nil
}
