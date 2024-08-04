package credentials

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	EmailTI, PassTI textinput.Model
}

func InitialModel() Model {
	cm := Model{
		EmailTI: textinput.New(),
		PassTI:  textinput.New(),
	}
	cm.PassTI.EchoMode = textinput.EchoPassword
	return cm
}

//#region tea.Model implementation

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return fmt.Sprintf("Email%v\nPassword%v\n", m.EmailTI.View(), m.PassTI.View())
}

//#endregion
