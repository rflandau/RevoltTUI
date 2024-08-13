package credentials

import (
	"fmt"
	"revolt_tui/log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sentinelb51/revoltgo"
)

type selected uint8

const (
	email selected = iota
	pass
)

type Model struct {
	Killed          bool
	emailTI, passTI textinput.Model
	sel             selected
	Session         *revoltgo.Session
	inputErr        string
}

func InitialModel() Model {
	cm := Model{
		emailTI: textinput.New(),
		passTI:  textinput.New(),
	}
	cm.passTI.EchoMode = textinput.EchoPassword
	cm.sel = email
	cm.emailTI.Focus()
	return cm
}

//#region tea.Model implementation

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// check for special keys
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		// clear a preexisting error
		m.inputErr = ""
		switch keyMsg.Type {
		case tea.KeyTab, tea.KeyUp, tea.KeyDown:
			return m, m.switchSelected()
		case tea.KeyCtrlC, tea.KeyEsc:
			m.Killed = true
			return m, tea.Quit
		case tea.KeyEnter:
			// attempt to login
			sess, lr, err := revoltgo.NewWithLogin(revoltgo.LoginData{
				Email:        strings.TrimSpace(m.emailTI.Value()),
				Password:     m.passTI.Value(),
				FriendlyName: "TUIFriendly", // TODO
			})
			if err != nil {
				log.Writer.Error("login failed", "error", err)
				m.inputErr = err.Error()
				return m, textinput.Blink
			}
			log.Writer.Debug("completed login attempt", "loginResponse", lr)
			m.Session = sess
			return m, tea.Quit
		}

	}

	var cmd tea.Cmd
	if m.sel == email {
		m.emailTI, cmd = m.emailTI.Update(msg)
	} else {
		m.passTI, cmd = m.passTI.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf("Email%v\nPassword%v\n%s\n",
		m.emailTI.View(), m.passTI.View(), m.inputErr)
}

//#endregion

func (m *Model) switchSelected() tea.Cmd {
	if m.sel == email {
		m.emailTI.Blur()
		m.sel = pass
		return m.passTI.Focus()
	}
	m.passTI.Blur()
	m.sel = email
	return m.emailTI.Focus()
}
