package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"peeporun/internal/config"
)

func (a *App) editingPreset() *config.Preset {
	if a.editPresetIdx < 0 || a.editPresetIdx >= len(a.presets) {
		return nil
	}
	return &a.presets[a.editPresetIdx]
}

func (a *App) updatePresetEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if a.editMode {
		return a.updatePresetEditInput(msg)
	}

	key := msg.String()
	kb := a.kb
	p := a.editingPreset()
	if p == nil {
		a.screen = screenPresetSelect
		return a, nil
	}
	maxCursor := len(p.Splits) // "add split" virtual row

	switch {
	case config.Matches(key, kb.Up):
		a.editCursor--
		if a.editCursor < -2 {
			a.editCursor = maxCursor
		}
	case config.Matches(key, kb.Down):
		a.editCursor++
		if a.editCursor > maxCursor {
			a.editCursor = -2
		}
	case config.Matches(key, kb.Confirm), config.Matches(key, kb.Rename):
		a.beginEditField()
	case config.Matches(key, kb.Add):
		a.beginAddSplit()
	case config.Matches(key, kb.Delete):
		if a.editCursor >= 0 && a.editCursor < len(p.Splits) {
			a.confirm = confirmDeleteSplit
		}
	case config.Matches(key, kb.MoveUp):
		if a.editCursor > 0 && a.editCursor < len(p.Splits) {
			p.Splits[a.editCursor-1], p.Splits[a.editCursor] = p.Splits[a.editCursor], p.Splits[a.editCursor-1]
			a.editCursor--
			a.persist()
		}
	case config.Matches(key, kb.MoveDn):
		if a.editCursor >= 0 && a.editCursor < len(p.Splits)-1 {
			p.Splits[a.editCursor+1], p.Splits[a.editCursor] = p.Splits[a.editCursor], p.Splits[a.editCursor+1]
			a.editCursor++
			a.persist()
		}
	case config.Matches(key, kb.Cancel):
		a.screen = screenPresetSelect
		a.syncSaveShape()
	}
	return a, nil
}

func (a *App) beginEditField() {
	p := a.editingPreset()
	if p == nil {
		return
	}
	switch {
	case a.editCursor == -2:
		a.editAction = "game"
		a.input.SetValue(p.Game)
	case a.editCursor == -1:
		a.editAction = "category"
		a.input.SetValue(p.Category)
	case a.editCursor >= 0 && a.editCursor < len(p.Splits):
		a.editAction = "rename-split"
		a.input.SetValue(p.Splits[a.editCursor])
	default:
		a.beginAddSplit()
		return
	}
	a.editMode = true
	a.input.Focus()
	a.input.CursorEnd()
}

func (a *App) beginAddSplit() {
	a.editAction = "add-split"
	a.input.SetValue("")
	a.editMode = true
	a.input.Focus()
}

func (a *App) updatePresetEditInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		a.commitEditField()
		a.editMode = false
		a.input.Blur()
		return a, nil
	case "esc":
		a.editMode = false
		a.input.Blur()
		return a, nil
	}
	var cmd tea.Cmd
	a.input, cmd = a.input.Update(msg)
	return a, cmd
}

func (a *App) commitEditField() {
	p := a.editingPreset()
	if p == nil {
		return
	}
	val := strings.TrimSpace(a.input.Value())
	if val == "" {
		return
	}
	switch a.editAction {
	case "game":
		p.Game = val
	case "category":
		p.Category = val
	case "rename-split":
		if a.editCursor >= 0 && a.editCursor < len(p.Splits) {
			p.Splits[a.editCursor] = val
		}
	case "add-split":
		p.Splits = append(p.Splits, val)
		a.editCursor = len(p.Splits) - 1
	}
	a.syncSaveShapeFor(*p)
	a.persist()
}

func (a *App) doDeleteSplit() {
	p := a.editingPreset()
	if p == nil || a.editCursor < 0 || a.editCursor >= len(p.Splits) {
		return
	}
	idx := a.editCursor
	p.Splits = append(p.Splits[:idx], p.Splits[idx+1:]...)
	if a.editCursor >= len(p.Splits) {
		a.editCursor = len(p.Splits) - 1
	}
	if a.editCursor < -2 {
		a.editCursor = -2
	}
	a.syncSaveShapeFor(*p)
	a.persist()
}

func (a *App) viewPresetEdit() string {
	p := a.editingPreset()
	if p == nil {
		return StyleHelp.Render("No preset selected.")
	}

	var b strings.Builder
	b.WriteString(StyleGameTitle.Render("Editing Preset"))
	b.WriteString("\n\n")

	b.WriteString(fieldRow("Game:", p.Game, a.editCursor == -2 && !a.editMode))
	b.WriteString("\n")
	b.WriteString(fieldRow("Category:", p.Category, a.editCursor == -1 && !a.editMode))
	b.WriteString("\n\n")

	if a.editMode && (a.editAction == "game" || a.editAction == "category") {
		b.WriteString(StyleActiveRow.Render(" > " + a.input.View()))
		b.WriteString("\n\n")
	}

	b.WriteString(StyleHeader.Render("Splits:"))
	b.WriteString("\n")

	for i, name := range p.Splits {
		active := a.editCursor == i && !a.editMode
		row := "  " + name
		if active {
			row = "> " + name
		}
		if a.editMode && a.editAction == "rename-split" && a.editCursor == i {
			b.WriteString(StyleActiveRow.Render("> " + a.input.View()))
		} else if active {
			b.WriteString(StyleActiveRow.Render(row))
		} else {
			b.WriteString(StyleNormalRow.Render(row))
		}
		b.WriteString("\n")
	}

	addRow := "  + Add split"
	if a.editCursor == len(p.Splits) && !a.editMode {
		addRow = "> + Add split"
	}
	if a.editMode && a.editAction == "add-split" {
		b.WriteString(StyleActiveRow.Render("> " + a.input.View()))
	} else if a.editCursor == len(p.Splits) {
		b.WriteString(StyleActiveRow.Render(addRow))
	} else {
		b.WriteString(StyleNormalRow.Render(addRow))
	}
	b.WriteString("\n\n")

	if a.editMode {
		b.WriteString(StyleHelp.Render(joinHelp("enter confirm", "esc cancel")))
	} else {
		kb := a.kb
		b.WriteString(StyleHelp.Render(joinHelp(
			pair(kb.Up, kb.Down)+" move",
			altPair(kb.Confirm, kb.Rename)+" rename",
			label(kb.Add, "add"),
			label(kb.Delete, "delete"),
			pair(kb.MoveUp, kb.MoveDn)+" reorder",
			label(kb.Cancel, "back"),
		)))
	}

	return StyleBox.Render(b.String())
}

func fieldRow(label, value string, active bool) string {
	row := "  " + label + " " + value
	if active {
		row = "> " + label + " " + value
		return StyleActiveRow.Render(row)
	}
	return StyleNormalRow.Render(row)
}
