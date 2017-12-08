package modelsfakes

import (
	"petbot/models"
)

type FakePetTableKey struct {
	OwnerName, PetName string
}

type FakePetTableColumns struct {
	OwnerName, PetName string
}

type FakePetTable struct {
	FakePetTableData    map[FakePetTableKey]FakePetTableColumns
	GetAllPetsCallCount int
	AddPetInfoCallCount int
}

type FakeDataStore struct {
	FakePetTable
}

func Init(_ string, _ string) (*FakeDataStore, error) {
	fakeDataStore := new(FakeDataStore)
	fakeDataStore.FakePetTableData = make(map[FakePetTableKey]FakePetTableColumns)
	fakeDataStore.GetAllPetsCallCount = 0
	return fakeDataStore, nil
}

func (fakeDataStore *FakeDataStore) CreatePetTable() {
}

func (fakeDataStore *FakeDataStore) GetAllPets() ([]*models.Pet, error) {
	fakeDataStore.GetAllPetsCallCount++
	values := fakeDataStore.FakePetTableData
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
		fakeDataStore.FakePetTableData[FakePetTableKey{owner, petName}] = FakePetTableColumns{owner, petName}
		return "Added"
	}
}

func (fakeDataStore *FakeDataStore) UpdatePetInfo(owner string, petName string, newPetName string) {
}

func (fakeDataStore *FakeDataStore) HasPetInfo(owner string, petName string) (bool) {
	value, exists := fakeDataStore.FakePetTableData[FakePetTableKey{owner, petName}]
	if exists && value.PetName == petName {
		return true
	} else {
		fakeDataStore.FakePetTableData[FakePetTableKey{owner,petName}] = FakePetTableColumns{owner,petName}
		return false
	}
}

func (fakeDataStore *FakeDataStore) HasPetTable() (bool) {
	return true
}