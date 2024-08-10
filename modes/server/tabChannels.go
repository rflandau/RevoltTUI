package server

import (
	"revolt_tui/broker"

	"github.com/charmbracelet/bubbles/list"
	"github.com/sentinelb51/revoltgo"
)

/**
 * This file handles the operation and components of the channels tab.
 */

type tabChannels struct {
	list          list.Model
	activeChannel *revoltgo.Channel
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

// channel representation for the channelList list.Model
type channelItem struct {
	name        string
	description string
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
