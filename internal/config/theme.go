package config

import (
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

const DefaultAccentColor = "#FFD400"

var hexColorRe = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

// ValidHexColor reports whether s is a "#RRGGBB" hex color string.
func ValidHexColor(s string) bool {
	return hexColorRe.MatchString(s)
}

func DefaultTheme() ThemeSettings {
	return ThemeSettings{AccentColor: DefaultAccentColor}
}

// LoadTheme reads theme.yaml, creating it with the default yellow on first
// run. If the configured color isn't a valid "#RRGGBB" hex string, it falls
// back to the default rather than failing to start.
func LoadTheme() (ThemeSettings, error) {
	path, err := ThemePath()
	if err != nil {
		return ThemeSettings{}, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		defaults := DefaultTheme()
		if werr := SaveTheme(defaults); werr != nil {
			return defaults, werr
		}
		return defaults, nil
	}
	if err != nil {
		return ThemeSettings{}, err
	}

	var t ThemeSettings
	if err := yaml.Unmarshal(data, &t); err != nil {
		return ThemeSettings{}, err
	}
	if !ValidHexColor(t.AccentColor) {
		t.AccentColor = DefaultAccentColor
	}
	return t, nil
}

func SaveTheme(t ThemeSettings) error {
	path, err := ThemePath()
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(t)
	if err != nil {
		return err
	}
	header := "# The standard theme color is " + DefaultAccentColor + "\n"
	return os.WriteFile(path, append([]byte(header), data...), 0o644)
}
