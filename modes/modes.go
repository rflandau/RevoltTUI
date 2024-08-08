/*
The modes package organizes the various modes available to the user and their subroutines.
Each mode must register the functions required for it to work:
 1. ChangeMode. ChangeMode is called by the controller at the end of each update cycle, to see if the current mode is ready to be swapped out.
    If it returns true, the controller will swap to the new mode and call its Initialize.
 2. Enter. Enter is called whenever this mode is switched *into* to set up any required data within the mode itself.
    Returns whether or not the setup was successful (the action should do its own logging, hence a simple bool).
 3. Update. Update is the function that the controller calls in place of its own Update() function.
 4. View. Like Update, but for View functions.
*/
package modes

import (
	"revolt_tui/log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sentinelb51/revoltgo"
)

type Mode uint8

const (
	// Selecting a server to interact with from all available servers
	ServerSelection Mode = iota
	// Interacting with a selected server
	Server
)

type Action interface {
	// Called at the end of controller's update to determine if the current mode wants to yield control and to whom.
	ChangeMode() (bool, Mode)
	// Called after ChangeMode, to allow the new mode to ready itself
	Enter() (success bool, init tea.Cmd)
	Update(s *revoltgo.Session, msg tea.Msg) tea.Cmd
	View() string
}

var modes map[Mode]Action = make(map[Mode]Action)

func Add(mode Mode, action Action) {
	modes[mode] = action
}

func Get(mode Mode) Action {
	log.Writer.Debug("fetching action from mode", "map", modes, "mode", mode)
	action := modes[mode]
	if action == nil {
		log.Writer.Fatal("no actions associated", "map", modes, "mode", mode)
	}
	return modes[mode]
}
