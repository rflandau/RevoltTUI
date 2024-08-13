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

type chnl struct {
	list          list.Model
	activeChannel *revoltgo.Channel
	selectionErr  string
}

var _ tab = &chnl{}

const changeChannelErrString string = "an error has occurred changing channel to "

func (tc *chnl) Name() string {
	return "channels"
}

func (tc *chnl) Enabled() bool {
	return true
}

func (c *chnl) Enter(s *revoltgo.Server) {
	// s is nil checked prior to call
	var itms []list.Item = make([]list.Item, len(s.Channels))
	for i, ch := range s.Channels {
		itms[i] = channelItem{name: ch}
	}

	c.list = list.New(itms, list.NewDefaultDelegate(), broker.Width(), broker.Height())

}

func (tc *chnl) Update(msg tea.Msg) (tea.Cmd, tabConst) {
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

func (tc *chnl) View() string {
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
