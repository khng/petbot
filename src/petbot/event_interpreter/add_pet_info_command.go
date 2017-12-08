package event_interpreter

import (
	"petbot/models"
)

func AddPetInfo(message *Message, rtm SlackRTM, petDataStore models.PetDataStore) string {
	info := rtm.GetInfo()
	user := info.GetUserByID(message.sender)

	return petDataStore.AddPetInfo(user.RealName, message.body)
}