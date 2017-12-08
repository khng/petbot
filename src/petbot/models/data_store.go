package models

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

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