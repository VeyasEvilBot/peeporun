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
	return ThemeSettings{AccentColor: DefaultAccentColor, ShowPB: true}
}

// LoadTheme reads theme.yaml, creating it with the defaults on first run.
// If the configured color isn't a valid "#RRGGBB" hex string, it falls
// back to the default rather than failing to start. If show_pb is absent
// from an older theme.yaml (from before this setting existed), it
// defaults to true rather than Go's normal bool zero-value, so upgrading
// doesn't silently hide PB for existing users.
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

	var raw map[string]interface{}
	_ = yaml.Unmarshal(data, &raw) // best-effort, just to check key presence
	_, hasShowPB := raw["show_pb"]

	var t ThemeSettings
	if err := yaml.Unmarshal(data, &t); err != nil {
		return ThemeSettings{}, err
	}
	if !ValidHexColor(t.AccentColor) {
		t.AccentColor = DefaultAccentColor
	}
	if !hasShowPB {
		t.ShowPB = true
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
	header := "" +
		"# The standard theme color is " + DefaultAccentColor + "\n" +
		"# show_pb controls whether Personal Best is shown, in BOTH the TUI\n" +
		"# and the OBS overlay - set to false to hide it in both places.\n"
	return os.WriteFile(path, append([]byte(header), data...), 0o644)
}
