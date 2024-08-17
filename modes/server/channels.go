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

type channelTab struct {
	list          list.Model
	activeChannel *revoltgo.Channel
	selectionErr  string
}

var _ tab = &channelTab{}

const changeChannelErrString string = "an error has occurred changing channel to "

func (tc *channelTab) Name() string {
	return "channels"
}

func (tc *channelTab) Enabled() bool {
	return true
}

func (c *channelTab) Init(s *revoltgo.Server, width, height int) {
	// s is nil checked prior to call
	var itms []list.Item // = make([]list.Item, len(s.Channels))
	for _, chID := range s.Channels {
		ci := channelItem{channelID: chID}
		if channel, err := broker.Session.Channel(chID); err != nil {
			log.Writer.Warn("failed to fetch channel", "id", chID, "error", err)
			ci.name = "[unknown]"
			ci.description = "failed to retrieve channel information"
		} else if channel.ChannelType == revoltgo.ChannelTypeText {
			ci.name = channel.Name
			ci.description = channel.Description
			ci.channel = channel

			itms = append(itms, ci)
		}

		//itms[i] = ci
	}

	c.list = list.New(itms, list.NewDefaultDelegate(), width, 30)

}

func (tc *channelTab) Update(msg tea.Msg) (tea.Cmd, tabConst) {
	// window size updates are handled by the main server Update; only need to check for keymsg
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEnter {
		baseItm := tc.list.SelectedItem()
		itm, ok := baseItm.(channelItem)
		if !ok {
			log.Writer.Warn("Failed to set active channel: base list item failed cast", "base item", baseItm)
			tc.selectionErr = changeChannelErrString + baseItm.FilterValue()
			return nil, CHANNELS
		}
		if itm.channel == nil {
			log.Writer.Warn("Failed to set active channel: cast item.channel is nil", "item", itm)
			tc.selectionErr = changeChannelErrString + baseItm.FilterValue()
			return nil, CHANNELS
		}
		tc.activeChannel = itm.channel
		// switch to chat channel
		return textinput.Blink, CHAT
	}

	var cmd tea.Cmd
	tc.list, cmd = tc.list.Update(msg)
	return cmd, CHANNELS
}

func (tc *channelTab) View() string {
	var sb strings.Builder
	sb.WriteString(tc.selectionErr + "\n")
	sb.WriteString(tc.list.View())

	return sb.String()
}

// channel representation for the channelList list.Model
type channelItem struct {
	name        string
	description string
	channelID   string
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
