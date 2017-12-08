package event_interpreter

import (
	"fmt"
	"strings"
	"petbot/models"
)

func GetAllPetsAsFormattedString(petDataStore models.PetDataStore) string {
	values, _ := petDataStore.GetAllPets()
	valuesFormatted := fmt.Sprintf("%+v", values)
	valuesFormatted = strings.Trim(valuesFormatted, "[]")
	return valuesFormatted
}