package serverselection

import (
	"revolt_tui/cache"
	"revolt_tui/log"
	"revolt_tui/modes"
	"revolt_tui/terminal"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sentinelb51/revoltgo"
)

type Action struct {
	list        list.Model
	initialized bool
}

func (a *Action) ChangeMode() (bool, modes.Mode) {
	return false, modes.ServerSelection
}

func (a *Action) Enter() (bool, tea.Cmd) {
	var cmd tea.Cmd
	if a.initialized {
		// if we have already been initialized, update current values
		a.list.SetWidth(terminal.Width())
		a.list.SetHeight(terminal.Height())

		// regenerate servers list as []list.Item
		// the server list should never become nil after being initialized, but you never know
		if servers := cache.Servers(); servers != nil {
			cmd = a.list.SetItems(castServersToItems(servers))
		} // else, do nothing

	} else if !a.tryInitialize() {
		log.Writer.Debug("not yet able to initialize")
	}
	return true, cmd
}

func (a *Action) Update(session *revoltgo.Session, msg tea.Msg) tea.Cmd {
	if !a.initialized && !a.tryInitialize() { //retry initialization
		return nil
	}
	var cmd tea.Cmd
	a.list, cmd = a.list.Update(msg)
	return cmd
}

func (a *Action) View() string {
	if !a.initialized {
		return "Initializing..."
	}
	return a.list.View()
}

//#region helper functions

func castServersToItems(servers []*revoltgo.Server) []list.Item {
	var itms []list.Item = make([]list.Item, len(servers))

	for i, server := range servers {
		itms[i] = serverItem{
			title:       server.Name,
			description: server.Description,
			id:          server.ID,
		}
	}

	return itms
}

func (a *Action) tryInitialize() bool {
	w, h := terminal.Width(), terminal.Height()
	if !cache.Ready() || w == 0 || h == 0 {
		return false
	}
	log.Writer.Debug("initializing server selection...")

	// if we have not been initialized, attempt to initialize
	a.list = list.New(castServersToItems(cache.Servers()), list.NewDefaultDelegate(), w, h)
	a.initialized = true

	return true
}

//#endregion

//#region list item definition

type serverItem struct {
	list.Item
	title       string
	description string
	id          string // server item for lookup upon selection
}

func (li *serverItem) Title() string {
	return li.title
}
func (li *serverItem) Description() string {
	return li.description
}

//#endregion
