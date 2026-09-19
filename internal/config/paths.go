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
	return filepath.Join(dir, "keys.toml"), nil
}

// PresetsDir returns ~/.config/peeporun/presets/ - one .toml file per
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

// LegacyPresetsPath is the old single-file presets.toml location, kept
// around only so LoadPresets can detect and migrate it on first run
// after upgrading to the folder-per-preset layout.
func LegacyPresetsPath() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "presets.toml"), nil
}

func SavePath() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "save.toml"), nil
}

func OverlaySettingsPath() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "overlay.toml"), nil
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
	return filepath.Join(dir, "theme.toml"), nil
}

// SocketPath returns the local Unix socket path used for CLI/hotkey
// control (`peeporun hit`, `peeporun split`, etc. talking to an already
// running instance). Prefers $XDG_RUNTIME_DIR (ephemeral, per-login-session,
// cleaned up automatically by the OS) when set, falling back to the config
// dir otherwise. Windows has no real equivalent to Unix sockets in the same
// way, so this feature is effectively Linux/macOS-only for now.
func SocketPath() string {
	if rt := os.Getenv("XDG_RUNTIME_DIR"); rt != "" {
		return filepath.Join(rt, "peeporun.sock")
	}
	if dir, err := Dir(); err == nil {
		return filepath.Join(dir, "peeporun.sock")
	}
	return filepath.Join(os.TempDir(), "peeporun.sock")
}
