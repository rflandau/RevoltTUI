package server

import (
	"fmt"
	"revolt_tui/broker"
	"revolt_tui/log"
	"revolt_tui/stylesheet"
	"revolt_tui/stylesheet/colors"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sentinelb51/revoltgo"
)

const (
	initialMessageFetchLimit int = 30
	messageRefreshLimit      int = 75
	viewportMessageLimit     int = 15
)

type chatTab struct {
	channelTab    *channelTab // this must be set on creation
	msgView       viewport.Model
	newMessageBox textarea.Model
	err           error
	msgs          messageStore
}

var _ tab = &chatTab{}

// Create the empty chat chat skeleton, with reference to its "parent" channel tab.
// Most initialization should go in Init, but this is required because it needs reference to the channel tab.
func LinkedChatTab(chTab *channelTab) *chatTab {
	return &chatTab{channelTab: chTab}
}

func (*chatTab) Name() string {
	return "chat"
}

// The chat tab is only available if a channel has been selected on the channel tab
func (cht *chatTab) Enabled() bool {
	return cht.channelTab.activeChannel != nil
}

func (cht *chatTab) Init(s *revoltgo.Server, width, height int) {
	// spawn a thread to watch for updates for when a channel is selected
	go func() {
		for {
			time.Sleep(time.Second * 15)
			if cht.channelTab.activeChannel == nil {
				continue
			}
			// TODO check for channel change; clear (or archive to SQLite) cache on channel change

			// if we have yet to fetch any messages, fetch the most recent set
			if cht.msgs.newestMessageID == "" {
				msgs, err := broker.Session.ChannelMessages(cht.channelTab.activeChannel.ID,
					revoltgo.ChannelMessagesParams{
						Limit: initialMessageFetchLimit,
						Sort:  revoltgo.ChannelMessagesParamsSortTypeLatest,
					})
				if err != nil {
					log.Writer.Warn("failed to fetch base channel message set",
						"channelID", cht.channelTab.activeChannel.ID,
						"error", err)
					continue
				}
				msgCount := len(msgs)
				cht.msgs.newestMessageID = msgs[msgCount-1].ID
				log.Writer.Debug("fetched base latest channel message set",
					"requested", initialMessageFetchLimit,
					"received", msgCount,
					"newest ID", cht.msgs.newestMessageID,
				)
				// set these messages as our current set
				cht.msgs.messages = msgs

				cht.populateViewport()
			} else { // fetch any new messages
				msgs, err := broker.Session.ChannelMessages(cht.channelTab.activeChannel.ID,
					revoltgo.ChannelMessagesParams{
						Limit: messageRefreshLimit,
						Sort:  revoltgo.ChannelMessagesParamsSortTypeLatest,
						After: cht.msgs.newestMessageID,
					})
				if err != nil {
					log.Writer.Warn("failed to refresh latest channel messages",
						"channelID", cht.channelTab.activeChannel.ID,
						"error", err)
					continue
				}
				msgCount := len(msgs)
				if msgCount == 0 { // no new messages
					continue
				}
				afterMessageID := cht.msgs.newestMessageID
				cht.msgs.newestMessageID = msgs[msgCount-1].ID
				log.Writer.Debug("refreshed latest channel messages",
					"requested", messageRefreshLimit,
					"received", msgCount,
					"previous newest ID", afterMessageID,
					"newest ID", cht.msgs.newestMessageID,
				)
				// if we reached up to our limit, a shit load of messages arrived while sleeping
				// we may need to query for older messages and insert that set into the middle of our array
				// TODO

				// prepend this new message set
				cht.msgs.messages = append(msgs, cht.msgs.messages...)

				cht.populateViewport()
			}
		}
	}()

	cht.newMessageBox = textarea.New()
	cht.newMessageBox.MaxHeight = 4
	cht.newMessageBox.Focus()

	// include height margins in the viewport
	cht.msgView = viewport.New(width, height-cht.newMessageBox.MaxHeight-1)

}

func (cht *chatTab) Update(msg tea.Msg) (tea.Cmd, tabConst) {

	// check for an enter key to submit the current state of the message compose area
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEnter {
		msgText := cht.newMessageBox.Value()
		if strings.TrimSpace(msgText) == "" {
			log.Writer.Debug("refusing to send empty message")
			return textarea.Blink, CHAT
		}
		// attempt to submit the message
		msg := revoltgo.MessageSend{Content: msgText}

		newMsg, err := broker.Session.ChannelMessageSend(cht.channelTab.activeChannel.ID, msg)
		if err != nil {
			log.Writer.Warn("failed to send message", "error", err)
			return textarea.Blink, CHAT
		}

		cht.msgs.messages = append(cht.msgs.messages, newMsg) // attach the message to our list of displayed messages
		cht.msgs.newestMessageID = newMsg.ID
		cht.populateViewport()

		cht.newMessageBox.SetValue("") // clear out the existing message
	}

	cmds := make([]tea.Cmd, 2)
	cht.msgView, cmds[0] = cht.msgView.Update(msg)
	cht.newMessageBox, cmds[1] = cht.newMessageBox.Update(msg)
	return tea.Batch(cmds...), CHAT
}

func (cht *chatTab) View() string {
	// draw a border around the message box to represent that it is highlighted

	existingMsgs := cht.msgView.View()
	compose := stylesheet.NewMessageComposeArea.Render(cht.newMessageBox.View())
	var errStr string
	if cht.err != nil {
		errStr = cht.err.Error()
	}

	return existingMsgs + "\n" + compose + "\n" + errStr
}

// The message store represents the current, local cache of messages able to be displayed
type messageStore struct {
	newestMessageID string              // newest, local message to fetch after
	messages        []*revoltgo.Message // local message cache, sorted newest [0] -> oldest [len]
}

var (
	timestampStyle lipgloss.Style = lipgloss.NewStyle().Foreground(colors.MessageTimestamp).Italic(true)
	authorStyle    lipgloss.Style = lipgloss.NewStyle().Foreground(colors.MessageAuthor)
)

// sets the content in chat's viewport, automatically jumping to the newest message (end of the VP) whenever called
func (cht *chatTab) populateViewport() {
	var sb strings.Builder
	totalMsgCount := len(cht.msgs.messages)
	for i := 0; i < viewportMessageLimit && i < totalMsgCount; i++ {
		if cht.msgs.messages[i] == nil {
			continue
		}
		sb.WriteString(displayMessage(cht.msgs.messages[i]) + "\n")
	}

	cht.msgView.GotoBottom()
}

// helper function for populateViewport(). Given a singular message, it returns a formatted string corresponding to its type.
// Note the lack of suffixed newlines.
func displayMessage(msg *revoltgo.Message) string {
	if msg == nil || msg.System == nil {
		return "undefined message"
	}

	switch msg.System.Type {
	case revoltgo.MessageSystemTypeText:
		return fmt.Sprintf("%s%s: %s", timestampStyle.Render(msg.Edited.Format(time.Stamp)), authorStyle.Render(msg.Author), msg.Content)
	case revoltgo.MessageSystemTypeChannelIconChanged:
		return fmt.Sprintf("%s changed their icon. Content: %s", msg.Author, msg.Content)
	default:
		log.Writer.Warn("unknown message type",
			"type", msg.System.Type, "mID", msg.ID)
		return "message of unknown type. Content: " + msg.Content
	}

}
