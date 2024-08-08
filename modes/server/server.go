/*
This package represents the server mode, where a user interacts with a selected server.
While there is not a "primary" user mode, this is likely be an important one.
*/
package server

/*
Server mode is a tabbed implementation of the standard revolt server view.
Tabs:
- Overview
- Channels
- Chat (empty if a channel has not been selected)
- Settings


NOTE: as tab is used as a primary navigation method, users cannot current insert tabs in chat messages.
*/

import (
	"revolt_tui/broker"
	"revolt_tui/log"
	"revolt_tui/modes"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sentinelb51/revoltgo"
)

type Action struct {
	server  *revoltgo.Server
	channel *revoltgo.Channel
}

var _ modes.Action = &Action{}

//#region Action Iface Impl

// Do we want to cede control to another mode?
func (a *Action) ChangeMode() (bool, modes.Mode) {
	return false, modes.Server
}

// Control was just passed to us, initialize as need be.
func (a *Action) Enter() (success bool, init tea.Cmd) {
	a.server = broker.GetServer()
	if a.server == nil {
		log.Writer.Errorf("control passed to server mode, but no server has been declared by Broker")
		return false, nil
	}

	// prepare the list of channels for the channel modal
	for i, ch := range a.server.Channels {

	}

	return true, textinput.Blink
}

func (a *Action) Update(s *revoltgo.Session, msg tea.Msg) tea.Cmd {

}

// Displays the current server, collapsing the channel column automatically if a channel has been
// selected and the terminal is not wide enough.
func (a *Action) View() string {
}

//#endregion
