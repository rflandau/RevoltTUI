package server

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sentinelb51/revoltgo"
)

type overviewTab struct {
	compiledOverview string
}

var _ tab = &overviewTab{}

// user-facing name of the tab
func (*overviewTab) Name() string {
	return "overview"
}

// is this tab currently accessible?
func (*overviewTab) Enabled() bool {
	return true
}

func (o *overviewTab) Init(s *revoltgo.Server, _, _ int) {
	// s is nil checked prior to call

	centerSty := lipgloss.NewStyle().AlignHorizontal(lipgloss.Center)

	// pre-generate the server overview
	var sb strings.Builder
	sb.WriteString(centerSty.Bold(true).Render(s.Name) + "\n")
	sb.WriteString(centerSty.Italic(true).Render(s.Description) + "\n")
	sb.WriteRune('\n')
	sb.WriteString(fmt.Sprintf("Owner: %v\n", s.Owner))
	sb.WriteString(fmt.Sprintf("ID: %v\n", s.ID))

	o.compiledOverview = sb.String()

	// TODO send out a goroutine to poll for overview updates
}

func (*overviewTab) Update(msg tea.Msg) (tea.Cmd, tabConst) {
	// no work to be done until server changes
	return nil, OVERVIEW
}
func (o *overviewTab) View() string {
	return o.compiledOverview
}
