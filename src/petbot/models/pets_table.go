package models

import (
	"log"
	"database/sql"
)

//func (dataStore *DataStore) FindOnePetInfo(owner string, petName string) {
//
//}

var tableName = "petinfo"

func (dataStore *DataStore) CreatePetTable() {
	_, err := dataStore.Exec("CREATE TABLE `" + tableName + "` (" +
		"`owner` VARCHAR(64) NOT NULL, " +
		"`pet_name` VARCHAR(64) NOT NULL, " +
		"PRIMARY KEY (`owner`, `pet_name`))")
	if err != nil {
		log.Fatal(err)
	}
}

func (dataStore *DataStore) GetAllPets() ([]*Pet, error) {
	rows, err := dataStore.Query("SELECT * FROM `" + tableName + "`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pets := make([]*Pet, 0)
	for rows.Next() {
		pet := new(Pet)
		err := rows.Scan(&pet.Owner, &pet.PetName)
		if err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}
	if err == rows.Err() && err != nil {
		return nil, err
	}
	return pets, err
}

func (dataStore *DataStore) AddPetInfo(owner string, petName string) (string) {
	if dataStore.HasPetInfo(owner, petName) {
		return "Duplicate"
	} else {
		_, err := dataStore.Exec("INSERT INTO `" + tableName + "` (`owner`,`pet_name`) " +
			"VALUES('" + owner + "','" + petName + "')")
		if err != nil {
			return "Error"
		}
		return "Added"
	}
}

func (dataStore *DataStore) UpdatePetInfo(owner string, petName string, newPetName string) {
	_, err := dataStore.Exec("UPDATE `" + tableName +
		"` SET pet_name='" + newPetName +
		"' WHERE owner='" + owner +
		"' AND pet_name='" + petName + "'")
	if err != nil {
		log.Fatal(err)
	}
}

func (dataStore *DataStore) HasPetInfo(owner string, petName string) (bool) {
	query := dataStore.QueryRow("SELECT * FROM `" + tableName +
		"` WHERE owner='" + owner + "' AND pet_name='" + petName + "'")
	err := query.Scan(&owner, &petName)
	if err != nil && err == sql.ErrNoRows {
		return false
	}
	return true
}

func (dataStore *DataStore) HasPetTable() (bool) {
	//_, err := dataStore.Exec("SELECT * FROM `%s` LIMIT 1", tableName)
	_, err := dataStore.Exec("SELECT * FROM `" + tableName + "` LIMIT 1")
	if err != nil {
		return false
	}
	return true
}