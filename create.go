package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func setData(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	setDataInstance := setDataType{}
	json.Unmarshal(reqBody, &setDataInstance)

	// If note or password are empty, return with empty ID to suggest failure
	if setDataInstance.Note == "" || setDataInstance.Pass == "" {
		w.WriteHeader(400)
		return
	}

	// If ID is empty, create an 8 character random one
	if setDataInstance.ID == "" || !storedDataEmpty(db[setDataInstance.ID]) {
		setDataInstance.ID = randString(8)
	}

	// Encrypt note with pass and fetch AAD, hash and ecryptedNote
	AAD, hash, encryptedNote := encrypt(setDataInstance.Note, setDataInstance.Pass)

	// Save AAD, AAD Hash and Encrypted note in db map
	db[setDataInstance.ID] = storedData{
		AADData: AAD,
		AADHash: hash,
		Note:    encryptedNote,
	}

	// On success, respond with proper ID
	output, _ := json.Marshal(setDataResponse{
		ID: setDataInstance.ID,
	})
	w.WriteHeader(200)
	_, _ = fmt.Fprintf(w, "%+v", string(output))
}
