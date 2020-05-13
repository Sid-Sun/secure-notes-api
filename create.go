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
		emptySetDataResponse(w)
		return
	}

	// If ID is empty, create an 8 character random one
	if exists, err := existsInDB(setDataInstance.ID); exists || setDataInstance.ID == "" {
		if err != nil {
			emptySetDataResponse(w)
			return
		}
		setDataInstance.ID = randString(8)
	}

	// Encrypt note with pass and fetch AAD, hash and ecryptedNote
	AAD, hash, encryptedNote := encrypt(setDataInstance.Note, setDataInstance.Pass)

	err := insertIntoDB(setDataInstance.ID, storedData{
		AADData: AAD,
		AADHash: hash[:],
		Note:    encryptedNote,
	})

	if err != nil {
		emptySetDataResponse(w)
		return
	}

	// On success, respond with proper ID
	output, _ := json.Marshal(setDataResponse{
		ID: setDataInstance.ID,
	})

	_, _ = fmt.Fprintf(w, "%+v", string(output))
}

func emptySetDataResponse(w http.ResponseWriter) {
	output, _ := json.Marshal(setDataResponse{
		ID: "",
	})
	_, _ = fmt.Fprintf(w, "%+v", string(output))
}
