package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cometpuppy/peeporun/internal/config"
)

// updateThemePicker handles key input for the in-TUI theme color picker.
func (a *App) updateThemePicker(msg tea.Msg) (tea.Model, tea.Cmd) {
	kb := a.kb

	// Text input mode (typing a custom hex code)
	if a.themeEditMode {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case config.Matches(msg.String(), kb.Confirm):
				a.input.Blur()
				hex := a.input.Value()
				if applyThemeFromInput(a, hex) {
					a.status = fmt.Sprintf("Theme: %s", hex)
				}
				a.themeEditMode = false
				a.screen = screenPresetSelect
			case config.Matches(msg.String(), kb.Cancel):
				a.themeEditMode = false
				a.input.Blur()
			default:
				// Let the textinput model handle normal typing
				var cmd tea.Cmd
				a.input, cmd = a.input.Update(msg)
				return a, cmd
			}
		}
		return a, nil
	}

	key := ""
	if msg, ok := msg.(tea.KeyMsg); ok {
		key = msg.String()
	}

	lastIdx := len(config.ThemePresets) // index of the "Custom" option

	switch {
	case config.Matches(key, kb.Up):
		if a.themeCursor > 0 {
			a.themeCursor--
		}
	case config.Matches(key, kb.Down):
		if a.themeCursor < lastIdx {
			a.themeCursor++
		}
	case config.Matches(key, kb.Confirm):
		if a.themeCursor < lastIdx {
			// Named color selected
			hex := config.ThemePresets[a.themeCursor].Color
			a.setThemeAccent(hex)
			a.status = fmt.Sprintf("Theme: %s", hex)
			a.screen = screenPresetSelect
		} else {
			// "Custom" option selected — enter text input mode
			a.themeEditMode = true
			a.input.SetValue(a.theme.AccentColor)
			a.input.Focus()
			a.input.CursorEnd()
		}
	case config.Matches(key, kb.Cancel):
		a.screen = screenPresetSelect
	}
	return a, nil
}

// applyThemeFromInput validates and applies a user-typed hex color string.
// Returns true if the color was applied, false if it was invalid.
func applyThemeFromInput(a *App, hex string) bool {
	if !config.ValidHexColor(hex) {
		a.status = "Invalid hex color (use #RRGGBB)"
		return false
	}
	a.setThemeAccent(hex)
	return true
}

// setThemeAccent updates the accent color across the TUI and overlay,
// saves it to theme.toml, and triggers an overlay rewrite. Changes
// take effect immediately — no restart needed.
func (a *App) setThemeAccent(hex string) {
	if !config.ValidHexColor(hex) {
		hex = config.DefaultAccentColor
	}
	a.theme.AccentColor = hex
	SetAccentColor(hex)
	_ = config.SaveTheme(a.theme)
	a.writeOverlay()
}

// viewThemePicker renders the theme color picker screen.
func (a *App) viewThemePicker() string {
	var b strings.Builder

	b.WriteString(StyleAppTitle.Render("Theme Picker"))
	b.WriteString("\n\n")
	b.WriteString(StyleCategory.Render("Pick an accent color for the TUI and overlay:"))
	b.WriteString("\n\n")

	presets := config.ThemePresets
	for i, nc := range presets {
		style := StyleNormalRow
		if i == a.themeCursor {
			style = StyleActiveRow
		}
		name := nc.Name
		if i == a.themeCursor {
			name = "> " + name
		} else {
			name = "  " + name
		}
		line := padRight(name, 13) + nc.Color
		b.WriteString(style.Render(line))
		b.WriteString("\n")
	}

	// Custom (hex) option — last item in the list
	lastIdx := len(presets)
	style := StyleNormalRow
	if a.themeCursor == lastIdx {
		style = StyleActiveRow
	}
	name := "  Custom (#RRGGBB)"
	if a.themeCursor == lastIdx {
		name = "> Custom (#RRGGBB)"
	}
	if a.themeEditMode {
		name += " — type and Enter"
	}
	b.WriteString(style.Render(padRight(name, 22)))
	b.WriteString("\n")

	b.WriteString("\n")
	kb := a.kb
	help := joinHelp(
		pair(kb.Up, kb.Down)+" move",
		label(kb.Confirm, "apply/select"),
		label(kb.Cancel, "back"),
	)
	b.WriteString(StyleHelp.Render(help))
	b.WriteString("\n")

	// Show current color
	b.WriteString(StyleHelp.Render(fmt.Sprintf("Current: %s", a.theme.AccentColor)))
	b.WriteString("\n")

	// Show text input if in custom mode
	if a.themeEditMode {
		b.WriteString(StyleHelp.Render("Enter hex: "))
		b.WriteString(a.input.View())
	}

	return StyleBox.Render(b.String())
}
