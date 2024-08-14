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
	"revolt_tui/stylesheet"
	"revolt_tui/stylesheet/colors"
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

	// determine the margins we need to reserve
	var w, h int = broker.Width(), broker.Height() - (lipgloss.Height(a.drawTabs()) + 2)

	// initialize each tab
	for _, tb := range a.tabs {
		tb.Init(a.server, w, h)
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

	// window size messages must be passed to every tab, lest they lost if a tab is unfocused
	if WSMsg, ok := msg.(tea.WindowSizeMsg); ok {
		// modify the height and width to fit within our content window beneath the tabs
		WSMsg.Width -= 2                                   // TODO calculate required margin
		WSMsg.Height = (lipgloss.Height(a.drawTabs()) + 2) // TODO extract to save cycles

		var (
			cmd    tea.Cmd
			newTab tabConst
		)
		for i, tb := range a.tabs {
			// NOTE: results are thrown away for all but the active tab
			c, t := tb.Update(WSMsg)
			if i == int(a.activeTab) {
				cmd = c
				newTab = t
			}
		}
		a.activeTab = newTab
		return cmd
	}

	var cmd tea.Cmd
	cmd, a.activeTab = a.tabs[a.activeTab].Update(msg)

	return cmd
}

var windowStyle = lipgloss.NewStyle().BorderForeground(colors.TabBorderForeground).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()

// Displays the current server, collapsing the channel column automatically if a channel has been
// selected and the terminal is not wide enough.
func (a *Action) View() string {
	var sb strings.Builder
	tabs := a.drawTabs()
	sb.WriteString(tabs + "\n")
	// box the entire display
	content := a.tabs[a.activeTab].View()
	sb.WriteString(windowStyle.Width(broker.Width()).Height(broker.Height() - (lipgloss.Height(tabs) + 1)).Render(content))
	return sb.String()
}

var (
	inactiveTabStyle = lipgloss.NewStyle().Border(stylesheet.TabBorders.Inactive, true).BorderForeground(colors.TabBorderForeground).Padding(0, 1)
	activeTabStyle   = inactiveTabStyle.Border(stylesheet.TabBorders.Active, true)
)

// helper function for View.
// Draws the tabs in their current state.
func (a *Action) drawTabs() string {
	var (
		renderedTabs []string
		margin       int = 2
		tabWidth     int = (broker.Width() - (margin * int(a.tabCount))) / int(a.tabCount)
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
