package config

// Preset defines a game + category and its ordered list of splits (bosses).
type Preset struct {
	ID       string   `toml:"id"`
	Game     string   `toml:"game"`
	Category string   `toml:"category"`
	Splits   []string `toml:"splits"`
}

// PresetsFile is the on-disk shape of presets.toml.
type PresetsFile struct {
	Presets []Preset `toml:"presets"`
}

// Keybinds maps actions to one or more key strings (bubbletea key.String() values).
type Keybinds struct {
	Up          []string `toml:"up"`
	Down        []string `toml:"down"`
	Hit         []string `toml:"hit"`
	Undo        []string `toml:"undo"`
	Split       []string `toml:"split"`
	Unsplit     []string `toml:"unsplit"`
	SavePB      []string `toml:"save_pb"`
	DeletePB    []string `toml:"delete_pb"`
	Reset       []string `toml:"reset"`
	Presets     []string `toml:"presets"`
	Quit        []string `toml:"quit"`
	Confirm     []string `toml:"confirm"`
	Cancel      []string `toml:"cancel"`
	New         []string `toml:"new"`
	Edit        []string `toml:"edit"`
	Delete      []string `toml:"delete"`
	Add         []string `toml:"add"`
	Rename      []string `toml:"rename"`
	MoveUp      []string `toml:"move_up"`
	MoveDn      []string `toml:"move_down"`
	ThemePicker []string `toml:"theme_picker"`
}

// Matches reports whether the given key string matches any bind in the set.
func Matches(key string, binds []string) bool {
	for _, b := range binds {
		if b == key {
			return true
		}
	}
	return false
}

// SplitState is the live/current run state for a single split.
type SplitState struct {
	Hits   int  `toml:"hits"`
	Beaten bool `toml:"beaten"`
}

// PBEntry is the personal-best hit count for a single split.
type PBEntry struct {
	Hits int `toml:"hits"`
}

// PresetSave is the persisted state (current run + PB) for one preset.
type PresetSave struct {
	Current []SplitState `toml:"current"`
	PB      []PBEntry    `toml:"pb"`
	HasPB   bool         `toml:"has_pb"`
	PBTotal int          `toml:"pb_total"`
}

// OverlaySettings controls the OBS-friendly HTML overlay output.
type OverlaySettings struct {
	Enabled        bool    `toml:"enabled"`
	RefreshSeconds float64 `toml:"refresh_seconds"`
}

// ThemeSettings controls the accent color shared by the TUI and the overlay,
// plus other display preferences shared between the two.
type ThemeSettings struct {
	AccentColor string `toml:"accent_color"`
	ShowPB      bool   `toml:"show_pb"`
}

// SaveFile is the on-disk shape of save.toml.
type SaveFile struct {
	LastPreset string                `toml:"last_preset"`
	Presets    map[string]PresetSave `toml:"presets"`
}
