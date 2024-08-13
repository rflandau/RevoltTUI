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
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sentinelb51/revoltgo"
)

type tabConst uint8

const (
	channels tabConst = iota
	chat
)

type Action struct {
	server *revoltgo.Server

	// tab management
	activeTab tabConst
	tabs      []string

	// tab data structs
	channels tabChannels
}

var _ modes.Action = &Action{}

func New() *Action {
	return &Action{}
}

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

	// generate each tab
	a.channels = initTabChannels(a.server) // channels
	// TODO // chat

	return true, textinput.Blink
}

func (a *Action) Update(s *revoltgo.Session, msg tea.Msg) tea.Cmd {
	var (
		cmd    tea.Cmd
		newTab tabConst = a.activeTab
	)
	switch a.activeTab {
	case channels:
		cmd, newTab = a.channels.update(msg)
	default:
		log.Writer.Warn("unknown active tab, restoring to channels", "active tab", a.activeTab)
	}

	if newTab != a.activeTab {
		// change to new tab
		// TODO
	}

	return cmd

}

// Displays the current server, collapsing the channel column automatically if a channel has been
// selected and the terminal is not wide enough.
func (a *Action) View() string {
	var sb strings.Builder
	sb.WriteString(a.drawTabs())
	switch a.activeTab {
	case channels:
		a.channels.view()
	}
	return sb.String()
}

// helper function for View.
// Draws the tabs in their current state.
func (a *Action) drawTabs() string {
	// draw each tab

	// if channels tab does not have a channel selected, do not display the chat tab
	if a.channels.activeChannel == nil {
		//TODO
	}

	// conjoin the drawn tabs
	return ""
}

//#endregion
