package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/sentinelb51/revoltgo"
)

type mainModel struct {
	session *revoltgo.Session
	log     *log.Logger
}

// model needs a logged in Client to proceed
func Initial(session *revoltgo.Session, log *log.Logger) mainModel {
	model := mainModel{session: session, log: log}
	// attach ready handler
	model.session.AddHandler(func(session *revoltgo.Session, r *revoltgo.EventReady) {
		model.log.Info("Ready to handle commands from %v user(s) across %d servers from %d channels",
			len(r.Users), len(r.Servers), len(r.Channels))
	})

	// attach message handler
	model.session.AddHandler(func(session *revoltgo.Session, r *revoltgo.Message) {
		// TODO
	})

	return model
}

//#region tea.Model implementation

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m mainModel) View() string {
	return ""
}

//#endregion
