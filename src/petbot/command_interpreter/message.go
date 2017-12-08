package command_interpreter

import (
	"github.com/nlopes/slack"
	"strings"
)

type Message struct {
	addressee, command, body, sender string
}

func (message *Message) ParseMessageEvent(messageEvent *slack.MessageEvent) {
	messageEventText := messageEvent.Text
	messageWords := strings.SplitN(messageEventText, " ", 3)
	messageWords = messageWords[:3]

	message.addressee = messageWords[0]
	message.command = messageWords[1]
	message.body = messageWords[2]
	message.sender = messageEvent.User
}
