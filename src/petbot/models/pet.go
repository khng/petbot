package models

type Pet struct {
	Owner string
	PetName string
	//isHere bool
}

func (pet *Pet) String() string {
	return "owner: '" + pet.Owner + "', pet name: '" + pet.PetName + "'\n"
}