package event_interpreter

import (
	"github.com/nlopes/slack"
	"fmt"
)

var dictionary = make(map[string]int)

type SlackRTM interface {
	GetInfo() *slack.Info
	SendMessage(message *slack.OutgoingMessage)
	NewOutgoingMessage(string, string) *slack.OutgoingMessage
}

func ParseTypingEvent(ev *slack.UserTypingEvent, rtm SlackRTM) {
	fmt.Printf("Typing Event Received from %s\n", ev.User)
	dictionary[ev.User]++

	if dictionary[ev.User] == 3 {
		dictionary[ev.User] = 0
		info := rtm.GetInfo()
		user := info.GetUserByID(ev.User)
		msg := fmt.Sprintf("looks like you have something to say, @%s. What's up?", user.Name)
		rtm.SendMessage(rtm.NewOutgoingMessage(msg, ev.Channel))
	}
}
