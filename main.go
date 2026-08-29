package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"peeporun/internal/config"
	"peeporun/internal/ui"
)

// version is set at build time via -ldflags "-X main.version=..."
// (see .goreleaser.yaml). Defaults to "dev" for local `go build`.
var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println("peepoRun " + version)
		return
	}

	kb, err := config.LoadKeybinds()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load keybinds:", err)
		os.Exit(1)
	}

	presets, err := config.LoadPresets()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load presets:", err)
		os.Exit(1)
	}

	save, err := config.LoadSave()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load save data:", err)
		os.Exit(1)
	}

	overlaySettings, err := config.LoadOverlaySettings()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load overlay settings:", err)
		os.Exit(1)
	}

	theme, err := config.LoadTheme()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load theme:", err)
		os.Exit(1)
	}

	app := ui.NewApp(kb, presets, save, overlaySettings, theme)

	p := tea.NewProgram(app, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
