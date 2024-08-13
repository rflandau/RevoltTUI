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
	"github.com/charmbracelet/lipgloss"
	"github.com/sentinelb51/revoltgo"
)

type Action struct {
	server *revoltgo.Server

	// tab management
	activeTab tabConst
	tabs      []tab
	tabCount  uint8 // set on startup, as tab count (enabled & disabled) should not change
}

var _ modes.Action = &Action{}

func New() *Action {
	a := &Action{}
	a.tabs = []tab{
		&ovrvw{},
		&chnl{},
	}
	a.tabCount = uint8(len(a.tabs))
	// check that we have an enumeration for each tab; this must be updated whenever a new tab enumeration is appended
	if a.tabCount != lastTabConst {
		log.Writer.Fatal("tab array count does not match enumeration count", "tab count", a.tabCount, "tabs", a.tabs, "last enumeration", lastTabConst)
	}

	return a
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

	// initialize each tab
	for _, tb := range a.tabs {
		tb.Enter(a.server)
	}

	// ensure we start on the always-enabled overview tab
	a.activeTab = overview

	return true, textinput.Blink
}

func (a *Action) Update(s *revoltgo.Session, msg tea.Msg) tea.Cmd {
	// consume tab cycle keys
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyTab:
			a.nextTab()
			return textinput.Blink
		case tea.KeyShiftTab:
			a.previousTab()
			return textinput.Blink
		} // all other inputs are unhandled
	}

	var cmd tea.Cmd
	cmd, a.activeTab = a.tabs[a.activeTab].Update(msg)

	return cmd
}

// Displays the current server, collapsing the channel column automatically if a channel has been
// selected and the terminal is not wide enough.
func (a *Action) View() string {
	var sb strings.Builder
	sb.WriteString(a.drawTabs() + "\n")
	sb.WriteString(a.tabs[a.activeTab].View())
	return sb.String()
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
)

// helper function for View.
// Draws the tabs in their current state.
func (a *Action) drawTabs() string {
	var (
		renderedTabs []string
		margin       int = 2
		tabWidth     int = (broker.Width() - (margin*int(a.tabCount) - 1)) / int(a.tabCount)
	)

	// draw each tab
	for i, t := range a.tabs {
		var style lipgloss.Style = inactiveTabStyle
		isFirst, isLast, isActive := i == 0, i == int(a.tabCount)-1, i == int(a.activeTab)
		if isActive {
			style = activeTabStyle
		}
		style.Width(tabWidth)
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t.Name()))

		if !t.Enabled() { // do not display disabled tabs
			continue
		}
	}

	// conjoin the drawn tabs
	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

//#endregion
