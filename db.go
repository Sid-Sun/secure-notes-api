package main

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var (
	host     = "postgres"
	port     = 5432
	user     = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = os.Getenv("POSTGRES_DB")
)

func initDB() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	sqlStatement := `
CREATE TABLE IF NOT EXISTS notes (
  id TEXT UNIQUE NOT NULL,
  aad TEXT NOT NULL,
  aad_hash TEXT NOT NULL,
  note TEXT NOT NULL,
  PRIMARY KEY (id)
);`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		return err
	}
	return nil
}

func insertIntoDB(ID string, data storedData) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	sqlStatement := `INSERT INTO notes (id, aad, aad_hash, note)
VALUES ($1, $2, $3, $4)
RETURNING id;`
	var id string
	err = db.QueryRow(sqlStatement, ID, base64.StdEncoding.EncodeToString(data.AADData), base64.StdEncoding.EncodeToString(data.AADHash[:]), base64.StdEncoding.EncodeToString(data.Note)).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

func fetchFromDB(ID string) (storedData, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return storedData{}, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return storedData{}, err
	}

	sqlStatement := `SELECT id, aad, aad_hash, note FROM notes WHERE id=$1;`
	var id string
	var rawAAD, rawHash, rawNote string
	row := db.QueryRow(sqlStatement, ID)
	err = row.Scan(&id, &rawAAD, &rawHash, &rawNote)
	if err != nil {
		return storedData{}, err
	}

	AAD, err := base64.StdEncoding.DecodeString(rawAAD)
	if err != nil {
		return storedData{}, errors.New("Could not parse stored data properly")
	}
	Hash, err := base64.StdEncoding.DecodeString(rawHash)
	if err != nil {
		return storedData{}, errors.New("Could not parse stored data properly")
	}
	Note, err := base64.StdEncoding.DecodeString(rawNote)
	if err != nil {
		return storedData{}, errors.New("Could not parse stored data properly")
	}

	return storedData{
		AADData: AAD,
		AADHash: Hash,
		Note:    Note,
	}, nil

}

func updateInDB(ID string, data storedData) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}
	sqlStatement := `
	UPDATE notes
	SET aad = $2, aad_hash = $3, note = $4
	WHERE id = $1
	RETURNING id;`
	var id string
	err = db.QueryRow(sqlStatement, ID, base64.StdEncoding.EncodeToString(data.AADData), base64.StdEncoding.EncodeToString(data.AADHash[:]), base64.StdEncoding.EncodeToString(data.Note)).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

func deleteFromDB(ID string) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	sqlStatement := `DELETE FROM notes WHERE id=$1 RETURNING id;`
	var id string
	err = db.QueryRow(sqlStatement, ID).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

func existsInDB(ID string) (bool, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return true, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return true, err
	}

	sqlStatement := `SELECT id FROM notes WHERE id=$1;`
	var id string
	err = db.QueryRow(sqlStatement, ID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return true, err
	}
	return true, nil
}
