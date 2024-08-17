package stylesheet

import (
	"revolt_tui/stylesheet/colors"

	"github.com/charmbracelet/lipgloss"
)

func TabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var TabBorders = struct {
	Inactive lipgloss.Border
	Active   lipgloss.Border
}{
	Inactive: TabBorderWithBottom("┴", "─", "┴"),
	Active:   TabBorderWithBottom("┘", " ", "└"),
}

var NewMessageComposeArea lipgloss.Style = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colors.TabBorderForeground)
