package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cometpuppy/peeporun/internal/config"
)

func (a *App) confirmMessage() string {
	switch a.confirm {
	case confirmReset:
		return "Reset this run to 0 hits?\n\n(PB will NOT be affected)"
	case confirmSavePB:
		return "Save this run as new Personal Best?"
	case confirmDeletePreset:
		return "Delete this preset? This cannot be undone."
	case confirmDeleteSplit:
		return "Delete this split from the preset?"
	case confirmDeletePB:
		return "Delete this preset's Personal Best?\n\n(This cannot be undone)"
	}
	return ""
}

func (a *App) updateConfirm(key string) (tea.Model, tea.Cmd) {
	kind := a.confirm

	switch {
	case config.Matches(key, a.kb.Confirm):
		a.confirm = confirmNone
		switch kind {
		case confirmReset:
			a.resetRun()
		case confirmSavePB:
			a.savePB()
		case confirmDeletePreset:
			a.doDeletePreset()
		case confirmDeleteSplit:
			a.doDeleteSplit()
		case confirmDeletePB:
			a.deletePB()
		}
	case config.Matches(key, a.kb.Cancel):
		a.confirm = confirmNone
	}
	return a, nil
}

func (a *App) viewConfirmOverlay(_ string) string {
	msg := a.confirmMessage()
	footer := "\n\n" + label(a.kb.Confirm, "confirm") + "   " + label(a.kb.Cancel, "cancel")
	return StyleModal.Render(msg + footer)
}
