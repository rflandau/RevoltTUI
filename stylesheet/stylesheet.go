package stylesheet

import "github.com/charmbracelet/lipgloss"

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
