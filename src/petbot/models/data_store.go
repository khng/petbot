package models

import (
	"log"
	"database/sql"
)

const TableName = "petinfo"
const OwnerColumn = "owner"
const PetNameColumn = "pet_name"

type PetDataStore interface {
	CreatePetTable()

	GetAllPets() ([]*Pet, error)
	AddPetInfo(owner string, petName string) (string)
	UpdatePetInfo(owner string, petName string, newPetName string)
	HasPetInfo(owner string, petName string) (bool)
	HasPetTable() (bool)
}

type DataStore struct {
	*sql.DB
}

func Init(driverName string, dataSourceName string) (*DataStore, error) {
	database, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = database.Ping(); err != nil {
		return nil, err
	}

	return &DataStore{database}, nil
}

func (dataStore *DataStore) CreatePetTable() {
	_, err := dataStore.Exec("CREATE TABLE `" + TableName + "` (" +
		"`" + OwnerColumn + "` VARCHAR(64) NOT NULL, " +
		"`" + PetNameColumn + "` VARCHAR(64) NOT NULL, " +
		"PRIMARY KEY (`" + OwnerColumn + "`, `" + PetNameColumn + "`))")

	if err != nil {
		log.Fatal(err)
	}
}

func (dataStore *DataStore) GetAllPets() ([]*Pet, error) {
	rows, err := dataStore.Query("SELECT * FROM `" + TableName + "`")

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
		_, err := dataStore.Exec("INSERT INTO `" + TableName +
			"` (`" + OwnerColumn + "`,`" + PetNameColumn + "`) " +
			"VALUES('" + owner + "','" + petName + "')")
		if err != nil {
			return "Error"
		}
		return "Added"
	}
}

func (dataStore *DataStore) UpdatePetInfo(owner string, petName string, newPetName string) {
	_, err := dataStore.Exec("UPDATE `" + TableName +
		"` SET " + PetNameColumn + "='" + newPetName +
		"' WHERE " + OwnerColumn + "='" + owner +
		"' AND " + PetNameColumn + "='" + petName + "'")
	if err != nil {
		log.Fatal(err)
	}
}

func (dataStore *DataStore) HasPetInfo(owner string, petName string) (bool) {
	query := dataStore.QueryRow("SELECT * FROM `" + TableName +
		"` WHERE " + OwnerColumn + "='" + owner +
		"' AND " + PetNameColumn + "='" + petName + "'")
	err := query.Scan(&owner, &petName)

	return err == nil || err != sql.ErrNoRows
}

func (dataStore *DataStore) HasPetTable() (bool) {
	_, err := dataStore.Exec("SELECT * FROM `" + TableName + "` LIMIT 1")

	return err == nil
}