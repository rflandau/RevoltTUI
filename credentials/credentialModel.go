package credentials

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type selected uint8

const (
	email selected = iota
	pass
)

type Model struct {
	Killed          bool
	EmailTI, PassTI textinput.Model
	sel             selected
}

func InitialModel() Model {
	cm := Model{
		EmailTI: textinput.New(),
		PassTI:  textinput.New(),
	}
	cm.PassTI.EchoMode = textinput.EchoPassword
	cm.sel = email
	cm.EmailTI.Focus()
	return cm
}

//#region tea.Model implementation

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// check for special keys
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyTab:
			return m, m.switchSelected()
		case tea.KeyCtrlC, tea.KeyEsc:
			m.Killed = true
			return m, tea.Quit
		}

	}

	// check for enter key to attempt to submit the crednetials and generate a session
	// TODO

	var cmd tea.Cmd
	if m.sel == email {
		m.EmailTI, cmd = m.EmailTI.Update(msg)
	} else {
		m.PassTI, cmd = m.PassTI.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf("Email%v\nPassword%v\n", m.EmailTI.View(), m.PassTI.View())
}

//#endregion

func (m *Model) switchSelected() tea.Cmd {
	if m.sel == email {
		m.EmailTI.Blur()
		m.sel = pass
		return m.PassTI.Focus()
	}
	m.PassTI.Blur()
	m.sel = email
	return m.EmailTI.Focus()
}
