package config

// Preset defines a game + category and its ordered list of splits (bosses).
type Preset struct {
	ID       string   `yaml:"id"`
	Game     string   `yaml:"game"`
	Category string   `yaml:"category"`
	Splits   []string `yaml:"splits"`
}

// PresetsFile is the on-disk shape of presets.yaml.
type PresetsFile struct {
	Presets []Preset `yaml:"presets"`
}

// Keybinds maps actions to one or more key strings (bubbletea key.String() values).
type Keybinds struct {
	Up       []string `yaml:"up"`
	Down     []string `yaml:"down"`
	Hit      []string `yaml:"hit"`
	Undo     []string `yaml:"undo"`
	Split    []string `yaml:"split"`
	Unsplit  []string `yaml:"unsplit"`
	SavePB   []string `yaml:"save_pb"`
	DeletePB []string `yaml:"delete_pb"`
	Reset    []string `yaml:"reset"`
	Presets  []string `yaml:"presets"`
	Quit     []string `yaml:"quit"`
	Confirm  []string `yaml:"confirm"`
	Cancel   []string `yaml:"cancel"`
	New      []string `yaml:"new"`
	Edit     []string `yaml:"edit"`
	Delete   []string `yaml:"delete"`
	Add      []string `yaml:"add"`
	Rename   []string `yaml:"rename"`
	MoveUp   []string `yaml:"move_up"`
	MoveDn   []string `yaml:"move_down"`
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
	Hits   int  `yaml:"hits"`
	Beaten bool `yaml:"beaten"`
}

// PBEntry is the personal-best hit count for a single split.
type PBEntry struct {
	Hits int `yaml:"hits"`
}

// PresetSave is the persisted state (current run + PB) for one preset.
type PresetSave struct {
	Current []SplitState `yaml:"current"`
	PB      []PBEntry    `yaml:"pb"`
	HasPB   bool         `yaml:"has_pb"`
	PBTotal int          `yaml:"pb_total"`
}

// OverlaySettings controls the OBS-friendly HTML overlay output.
type OverlaySettings struct {
	Enabled        bool    `yaml:"enabled"`
	RefreshSeconds float64 `yaml:"refresh_seconds"`
}

// ThemeSettings controls the accent color shared by the TUI and the overlay.
type ThemeSettings struct {
	AccentColor string `yaml:"accent_color"`
}

// SaveFile is the on-disk shape of save.yaml.
type SaveFile struct {
	LastPreset string                `yaml:"last_preset"`
	Presets    map[string]PresetSave `yaml:"presets"`
}
