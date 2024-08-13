package server

import (
	"revolt_tui/broker"
	"revolt_tui/log"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sentinelb51/revoltgo"
)

/**
 * This file handles the operation and components of the channels tab.
 */

type tabChannels struct {
	list          list.Model
	activeChannel *revoltgo.Channel
	selectionErr  string
}

// Initializes the list of channels
func initTabChannels(server *revoltgo.Server) tabChannels {
	tc := tabChannels{}

	var itms []list.Item = make([]list.Item, len(server.Channels))
	for i, ch := range server.Channels {
		itms[i] = channelItem{name: ch}
	}

	tc.list = list.New(itms, list.NewDefaultDelegate(), broker.Width(), broker.Height())

	return tc
}

const changeChannelErrString string = "an error has occurred changing channel to "

func (tc tabChannels) update(msg tea.Msg) (tea.Cmd, tabConst) {
	// window size updates are handled by the main server Update; only need to check for keymsg
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEnter {
		baseItm := tc.list.SelectedItem()
		itm, ok := baseItm.(channelItem)
		if !ok {
			log.Writer.Warn("Failed to set active channel: base list item failed cast", "base item", baseItm)
			tc.selectionErr = changeChannelErrString + baseItm.FilterValue()
			return nil, channels
		}
		if itm.channel == nil {
			log.Writer.Warn("Failed to set active channel: cast item.channel is nil", "item", itm)
			tc.selectionErr = changeChannelErrString + baseItm.FilterValue()
			return nil, channels
		}
		tc.activeChannel = itm.channel
		// switch to chat channel
		return textinput.Blink, chat
	}

	var cmd tea.Cmd
	tc.list, cmd = tc.list.Update(msg)
	return cmd, channels
}

func (tc tabChannels) view() string {
	var sb strings.Builder
	sb.WriteString(tc.selectionErr + "\n")
	sb.WriteString(tc.list.View())

	return sb.String()
}

// channel representation for the channelList list.Model
type channelItem struct {
	name        string
	description string
	channel     *revoltgo.Channel
}

var _ list.Item = channelItem{} // check interface

func (ci channelItem) Title() string {
	return ci.name
}

func (ci channelItem) Description() string {
	return ci.description
}

func (ci channelItem) FilterValue() string {
	return ci.name
}
