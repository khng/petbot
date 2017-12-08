package event_interpreter

import (
	"github.com/nlopes/slack"
	"fmt"
	"strings"
	"petbot/models"
)

var MissingCommandMessage = "Please input a command..."
var InvalidCommandMessage = "This is not a valid command..."

type Command interface {
	Execute(message *Message, rtm SlackRTM, petDataStore models.PetDataStore) string
}

type GetAllCommand struct {}
func (getAllCommand *GetAllCommand) Execute(message *Message, rtm SlackRTM, petDataStore models.PetDataStore) string {
	return GetAllPetsAsFormattedString(petDataStore)
}

type AddPetInfoCommand struct {}
func (addPetInfoCommand *AddPetInfoCommand) Execute(message *Message, rtm SlackRTM, petDataStore models.PetDataStore) string {
	return AddPetInfo(message, rtm, petDataStore)
}

func ExecuteCommand(messageEvent *slack.MessageEvent, rtm SlackRTM, petDataStore models.PetDataStore) {
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
