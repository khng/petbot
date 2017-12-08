package command_interpreter

import (
	"petbot/models"
)

type AddPetInfoCommand struct {}
func (addPetInfoCommand *AddPetInfoCommand) Execute(message *Message, rtm SlackRTM, petDataStore models.PetDataStore) string {
	return AddPetInfo(message, rtm, petDataStore)
}

func AddPetInfo(message *Message, rtm SlackRTM, petDataStore models.PetDataStore) string {
	info := rtm.GetInfo()
	user := info.GetUserByID(message.sender)

	return petDataStore.AddPetInfo(user.RealName, message.body)
}