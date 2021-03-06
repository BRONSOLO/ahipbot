package ahipbot

import (
	"fmt"

	"github.com/tkawachi/hipchat"
)
import "strings"

type BotReply struct {
	To      string
	Message string
}

type BotMessage struct {
	*hipchat.Message
	BotMentioned bool
	FromUser     *User
	FromRoom     *Room
}

func (msg *BotMessage) IsPrivate() bool {
	return msg.FromRoom == nil
}

func (msg *BotMessage) ContainsAnyCased(strs []string) bool {
	for _, s := range strs {
		if strings.Contains(msg.Body, s) {
			return true
		}
	}
	return false
}

func (msg *BotMessage) ContainsAny(strs []string) bool {
	lowerStr := strings.ToLower(msg.Body)

	for _, s := range strs {
		lowerInput := strings.ToLower(s)

		if strings.Contains(lowerStr, lowerInput) {
			return true
		}
	}
	return false
}

func (msg *BotMessage) Contains(s string) bool {
	lowerStr := strings.ToLower(msg.Body)
	lowerInput := strings.ToLower(s)

	if strings.Contains(lowerStr, lowerInput) {
		return true
	}
	return false
}

func (msg *BotMessage) Reply(s string) *BotReply {
	return &BotReply{
		To:      msg.From,
		Message: s,
	}
}

func (msg *BotMessage) ReplyPrivate(s string) *BotReply {
	return &BotReply{
		To:      msg.FromUser.JID,
		Message: s,
	}
}

func (msg *BotMessage) String() string {
	fromUser := "<unknown>"
	if msg.FromUser != nil {
		fromUser = msg.FromUser.Name
	}
	fromRoom := "<none>"
	if msg.FromRoom != nil {
		fromRoom = msg.FromRoom.Name
	}
	return fmt.Sprintf(`BotMessage{"%s", from_user=%s, from_room=%s, mentioned=%v, private=%v}`, msg.Body, fromUser, fromRoom, msg.BotMentioned, msg.IsPrivate())
}
