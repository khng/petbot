package command_interpreter

import (
	"fmt"
	"strings"
	"petbot/models"
)

type GetAllCommand struct {}
func (getAllCommand *GetAllCommand) Execute(message *Message, rtm SlackRTM, petDataStore models.PetDataStore) string {
	return GetAllPetsAsFormattedString(petDataStore)
}

func GetAllPetsAsFormattedString(petDataStore models.PetDataStore) string {
	values, _ := petDataStore.GetAllPets()
	valuesFormatted := fmt.Sprintf("%+v", values)
	valuesFormatted = strings.Trim(valuesFormatted, "[]")
	return valuesFormatted
}