package db

import (
	"golang.org/x/net/context"

	firebase "firebase.google.com/go"
	database "firebase.google.com/go/db"
)

// DataBase implemets the firebase DB
type DataBase struct {
	client *database.Client
}

// NewDB creates a new firebase database instance
func NewDB() (DataBase, error) {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return DataBase{}, err
	}
	_, err = app.Auth(context.Background())
	if err != nil {
		return DataBase{}, err
	}
	var dbClient DataBase
	dbClient.client, err = app.DatabaseWithURL(context.Background(), "https://notes-api-bb9f0.firebaseio.com/")
	if err != nil {
		return DataBase{}, err
	}
	return dbClient, nil
}

// Set the data at id to the provided data
func (db DataBase) Set(id string, data StoredData) error {
	if err := db.client.NewRef(id).Set(context.Background(), data); err != nil {
		return err
	}
	return nil
}

// Delete the data at id
func (db DataBase) Delete(id string) error {
	if err := db.client.NewRef(id).Delete(context.Background()); err != nil {
		return err
	}
	return nil
}

// Get the data at id
func (db DataBase) Get(id string) (StoredData, error) {
	var data StoredData
	if err := db.client.NewRef(id).Get(context.Background(), &data); err != nil {
		return StoredData{}, err
	}
	return data, nil
}

// Exists checks if data exists at ID in database
func (db DataBase) Exists(id string) (bool, error) {
	var data StoredData
	if err := db.client.NewRef(id).Get(context.Background(), &data); err != nil {
		return false, err
	}
	if storedDataEmpty(data) {
		return false, nil
	}
	return true, nil
}
