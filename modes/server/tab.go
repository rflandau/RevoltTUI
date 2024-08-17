package server

import (
	"revolt_tui/broker"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sentinelb51/revoltgo"
)

// enumeration of the tabs in the tab array, for clearer code
// NOTE: the order the structs are defined in the New() (on action creation) must match the order of enumeration
type tabConst = uint8

const (
	OVERVIEW tabConst = iota
	CHANNELS
	CHAT
)

const lastTabConst = CHAT // used by new to validate tab struct count

// represents a single tab
type tab interface {
	Name() string  // user-facing name of the tab
	Enabled() bool // is this tab currently accessible?
	// called on every tab when *server* is first entered, NOT when the tab is swapped to
	// provides the server so each enter does not have to nil check broker
	Init(server *revoltgo.Server, width, height int)
	Update(msg tea.Msg) (tea.Cmd, tabConst)
	View() string
}

// activate the next, enabled tab in index order
func (a *Action) nextTab() {
	// cycle through tabs until we find an enabled one
	for i := 0; i < int(a.tabCount); i++ { // catch infinite loop
		a.activeTab += 1
		if a.activeTab == a.tabCount {
			a.activeTab = 0
		}
		if a.tabs[a.activeTab].Enabled() {
			break
		}
	}
}

// activate the next, enabled tab in reverse index order
func (a *Action) previousTab() {
	// cycle through tabs until we find an enabled one
	for i := 0; i < int(a.tabCount); i++ { // catch infinite loop
		if a.activeTab == 0 { // if we are at the beginning, jump to the end
			a.activeTab = a.tabCount - 1
		} else {
			a.activeTab -= 1
		}

		if a.tabs[a.activeTab].Enabled() {
			break
		}
	}
}

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
		if !t.Enabled() { // do not display disabled tabs
			continue
		}
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

		var nameTxt string
		if isActive {
			nameTxt = activeTabTextStyle.Render(t.Name())
		} else {
			nameTxt = t.Name()
		}

		renderedTabs = append(renderedTabs, style.Render(nameTxt))
	}

	// conjoin the drawn tabs
	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}
