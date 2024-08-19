package server

import (
	"fmt"
	"revolt_tui/broker"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
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

	o.compiledOverview = generateOverview(s)

	// TODO send out a goroutine to poll for overview updates
}

// generates the overview display for the given server
func generateOverview(s *revoltgo.Server) string {
	// pre-generate the server overview
	var sb strings.Builder
	sb.WriteString(titleSty.Render(s.Name) + "\n")
	sb.WriteString(subtitleSty.Render(s.Description) + "\n")
	sb.WriteRune('\n')
	var notDiscover string
	if s.Discoverable != nil && !(*s.Discoverable) {
		notDiscover = "not "
	}

	// fetch data for the overview

	owner, err := broker.Session.User(s.Owner)
	if err != nil {
		log.Warn("failed to fetch server owner", "owner ID", s.Owner, "server ID", s.ID)
	}

	// generate and pair up fields+values
	left := lipgloss.JoinVertical(lipgloss.Right, leftAlignerSty.Render("Owner:"), leftAlignerSty.Render("ID:"))
	right := owner.Username + "\n" + s.ID

	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Center, left, right))

	sb.WriteString(fmt.Sprintf("\n\nThis server is %sdiscoverable.", notDiscover))

	return sb.String()
}

func (*overviewTab) Update(msg tea.Msg) (tea.Cmd, tabConst) {
	// no work to be done until server changes
	return nil, OVERVIEW
}
func (o *overviewTab) View() string {
	return o.compiledOverview
}

// #region styles

var (
	titleSty       lipgloss.Style = lipgloss.NewStyle().Bold(true)
	subtitleSty    lipgloss.Style = lipgloss.NewStyle().Italic(true)
	leftAlignerSty lipgloss.Style = lipgloss.NewStyle().AlignHorizontal(lipgloss.Right).Width(7).PaddingRight(1)
)

//#endregion styles
