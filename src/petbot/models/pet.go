package models

type Pet struct {
	Owner string
	PetName string
}

func (pet *Pet) String() string {
	return "owner: '" + pet.Owner + "', pet name: '" + pet.PetName + "'\n"
}