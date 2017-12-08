package modelsfakes

import (
	"petbot/models"
)

type Key struct {
	OwnerName, PetName string
}

type Columns struct {
	OwnerName, PetName string
}

type FakeDataStore struct {
	Data                map[Key] Columns
	GetAllPetsCallCount int
	AddPetInfoCallCount int
}

func Init(_ string, _ string) (*FakeDataStore, error) {
	fakeDataStore := new(FakeDataStore)
	fakeDataStore.Data = make(map[Key] Columns)
	fakeDataStore.GetAllPetsCallCount = 0
	return fakeDataStore, nil
}

func (fakeDataStore *FakeDataStore) CreatePetTable() {
}

func (fakeDataStore *FakeDataStore) GetAllPets() ([]*models.Pet, error) {
	fakeDataStore.GetAllPetsCallCount++
	values := fakeDataStore.Data
	petArray := make([]*models.Pet, 0)
	for _, columns := range values {
		owner := columns.OwnerName
		petName := columns.PetName
		pet := models.Pet{owner, petName}
		petArray = append(petArray, &pet)
	}
	return petArray, nil
}

func (fakeDataStore *FakeDataStore) AddPetInfo(owner string, petName string) (string) {
	fakeDataStore.AddPetInfoCallCount++
	if fakeDataStore.HasPetInfo(owner, petName) {
		return "Duplicate"
	} else {
		fakeDataStore.Data[Key{owner, petName}] = Columns{owner, petName}
		return "Added"
	}
}

func (fakeDataStore *FakeDataStore) UpdatePetInfo(owner string, petName string, newPetName string) {
}

func (fakeDataStore *FakeDataStore) HasPetInfo(owner string, petName string) (bool) {
	value, exists := fakeDataStore.Data[Key{owner, petName}]
	if exists && value.PetName == petName {
		return true
	} else {
		fakeDataStore.Data[Key{owner,petName}] = Columns{owner,petName}
		return false
	}
}

func (fakeDataStore *FakeDataStore) HasPetTable() (bool) {
	return true
}