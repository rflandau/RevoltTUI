// The controller is the main, parent tea.Model that drives every other tea.Model and ensures
// control is passed around appropriately.
package controller

import (
	"revolt_tui/broker"
	"revolt_tui/log"
	"revolt_tui/modes"

	tea "github.com/charmbracelet/bubbletea"
)

type controller struct {
	quitting   bool
	mode       modes.Mode
	curAction  modes.Action
	initialCmd tea.Cmd
}

// model needs a logged in Client to proceed
func Initial() controller {
	model := controller{
		mode: modes.ServerSelection,
	}

	// enter the starter (server selection) mode
	log.Writer.Debug("controller entering initial mode", "mode", model.mode)
	model.curAction = modes.Get(model.mode)
	if success, init := model.curAction.Enter(); !success {
		// failure, dying...
		model.quitting = true
		return model
	} else {
		model.initialCmd = init
	}

	return model
}

//#region tea.Model implementation

func (ctl controller) Init() tea.Cmd {
	return ctl.initialCmd
}

func (ctl controller) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Writer.Debug("updating", "msg", msg)

	if ctl.quitting {
		return ctl, nil
	}
	// always handle kill keys, no matter the mode
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		// check for CTRL+C
		if keyMsg.Type == tea.KeyCtrlC {
			// clean up is handled by the program that originally initialized this model
			ctl.quitting = true
			return ctl, tea.Quit
		}
	}

	// capture window size
	if WSMsg, ok := msg.(tea.WindowSizeMsg); ok {
		broker.SetDimensions(WSMsg.Width, WSMsg.Height)
	}

	var cmd tea.Cmd = ctl.curAction.Update(msg)

	// check for a mode change
	if chg, newMode := ctl.curAction.ChangeMode(); chg {
		ctl.mode = newMode
		// fetch the action associated to the new mode
		ctl.curAction = modes.Get(ctl.mode)
		if success, init := ctl.curAction.Enter(); !success {
			// failure, dying...
			ctl.quitting = true
			return ctl, tea.Quit
		} else {
			return ctl, init
		}
	}
	return ctl, cmd
}

func (ctl controller) View() string {
	return ctl.curAction.View()
}

//#endregion
