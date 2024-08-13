package server

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sentinelb51/revoltgo"
)

type ovrvw struct {
	compiledOverview string
}

var _ tab = &ovrvw{}

// user-facing name of the tab
func (*ovrvw) Name() string {
	return "overview"
}

// is this tab currently accessible?
func (*ovrvw) Enabled() bool {
	return true
}

func (o *ovrvw) Enter(s *revoltgo.Server) {
	// s is nil checked prior to call

	centerSty := lipgloss.NewStyle().AlignHorizontal(lipgloss.Center)

	// pre-generate the server overview
	var sb strings.Builder
	sb.WriteString(centerSty.Bold(true).Render(s.Name) + "\n")
	sb.WriteString(centerSty.Italic(true).Render(s.Description) + "\n")
	sb.WriteRune('\n')
	sb.WriteString("Owner:" + s.Owner)
	sb.WriteString(fmt.Sprintf("ID: %v", s.ID))

	o.compiledOverview = sb.String()

	// TODO send out a goroutine to poll for overview updates
}

func (*ovrvw) Update(msg tea.Msg) (tea.Cmd, tabConst) {
	// no work to be done until server changes
	return nil, overview
}
func (o *ovrvw) View() string {
	return o.compiledOverview
}
