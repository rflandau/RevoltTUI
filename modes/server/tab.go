package server

import tea "github.com/charmbracelet/bubbletea"

// represents a single tab
type tab interface {
	Name() string  // user-facing name of the tab
	Enabled() bool // is this tab currently accessible?
	Update(msg tea.Msg) (tea.Cmd, tabConst)
	View() string
}
