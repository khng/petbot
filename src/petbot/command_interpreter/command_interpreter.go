package command_interpreter

import (
	"github.com/nlopes/slack"
	"fmt"
	"petbot/models"
)

const MissingCommandMessage = "Please input a command..."
const InvalidCommandMessage = "This is not a valid command..."

type SlackRTM interface {
	GetInfo() *slack.Info
	SendMessage(message *slack.OutgoingMessage)
	NewOutgoingMessage(string, string) *slack.OutgoingMessage
}

type Command interface {
	Execute(message *Message, rtm SlackRTM, petDataStore models.PetDataStore) string
}

func InterpretCommand(messageEvent *slack.MessageEvent, rtm SlackRTM, petDataStore models.PetDataStore) {
	message := Message{}
	message.ParseMessageEvent(messageEvent)
	info := rtm.GetInfo()
	botName := info.User.ID

	if MessageAddressedToPetBot(botName, message.addressee, message.sender) {
		commands := map[string]Command{
			"/all": &GetAllCommand{},
			"/add": &AddPetInfoCommand{},
		}

		if _, exists := commands[message.command]; exists {
			value := commands[message.command].Execute(&message, rtm, petDataStore)
			rtm.SendMessage(rtm.NewOutgoingMessage(value, messageEvent.Channel))
		} else if message.command == ""{
			rtm.SendMessage(rtm.NewOutgoingMessage(MissingCommandMessage, messageEvent.Channel))
		} else {
			rtm.SendMessage(rtm.NewOutgoingMessage(InvalidCommandMessage, messageEvent.Channel))
		}
	}
}

func MessageAddressedToPetBot(botName string, messageReceiver string, messageSender string) bool {
	botNameReferenceFormat := fmt.Sprintf("<@%s>", botName)
	return messageReceiver == botNameReferenceFormat && messageSender != botName
}