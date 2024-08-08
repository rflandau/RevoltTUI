package serverselection

import (
	"revolt_tui/broker"
	"revolt_tui/cache"
	"revolt_tui/log"
	"revolt_tui/modes"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sentinelb51/revoltgo"
)

type Action struct {
	list         list.Model
	initialized  bool
	selectionErr bool // an error occurred on last input
}

// Is this mode ready to change? If so, to what mode?
func (a *Action) ChangeMode() (bool, modes.Mode) {
	return false, modes.ServerSelection
}

// On user first entering this mode.
func (a *Action) Enter() (bool, tea.Cmd) {
	var cmd tea.Cmd
	if a.initialized {
		// if we have already been initialized, update current values
		a.list.SetWidth(broker.Width())
		a.list.SetHeight(broker.Height())

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
	if !a.initialized {
		if !a.tryInitialize() { //retry initialization
			return nil
		}
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		a.selectionErr = false
		if keyMsg.Type == tea.KeyEnter { // check for enter key
			// fetch the chosen server
			if serverItm, ok := a.list.SelectedItem().(serverItem); ok {
				server, err := session.Server(serverItm.id)
				if err != nil {
					log.Writer.Error("failed to fetch server", "error", err, "id", serverItm.id)
					a.selectionErr = true
					return nil
				}
				// pass the server to the app data broker
			} else {
				log.Writer.Warn("failed to cast item to server item", "item", a.list.SelectedItem())
				a.selectionErr = true
				return nil
			}
		}
	}

	var cmd tea.Cmd
	a.list, cmd = a.list.Update(msg)
	return cmd
}

func (a *Action) View() string {
	if !a.initialized {
		return "Initializing..."
	}
	l := a.list.View()

	// append error text, if an error had occurred
	if a.selectionErr {
		l += "\n An error has occurred, please try a different server."
	}

	return l
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
	w, h := broker.Width(), broker.Height()
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
	title       string
	description string
	id          string // server item for lookup upon selection
}

func (li serverItem) Title() string {
	return li.title
}
func (li serverItem) Description() string {
	return li.description
}
func (li serverItem) FilterValue() string {
	return li.title + li.description
}

//#endregion
