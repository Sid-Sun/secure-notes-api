package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"secure-notes-api/db"
)

func updateNote(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	updateNoteInstance := updateNoteType{}
	json.Unmarshal(reqBody, &updateNoteInstance)

	// If note, password or id are empty, return with empty ID To suggest failure
	if updateNoteInstance.Note == "" || updateNoteInstance.Pass == "" || updateNoteInstance.ID == "" {
		emptyUpdateResponse(w)
		return
	}

	if exists, err := dbInstance.Exists(updateNoteInstance.ID); exists && err == nil {
		// If note with ID exists in DB, verify its password
		// And throw the AAD - we don't need it to replace note
		data, err := dbInstance.Get(updateNoteInstance.ID)
		if err != nil {
			emptyUpdateResponse(w)
			return
		}
		_, err = verifyNotePassword(data, updateNoteInstance.Pass)
		if err == nil {
			// If NewPass is supplied, set pass as newpass before encrypting
			if updateNoteInstance.NewPass != "" {
				updateNoteInstance.Pass = updateNoteInstance.NewPass
			}
			AAD, hash, encryptedNote := encrypt(updateNoteInstance.Note, updateNoteInstance.Pass)

			// Save new AAD, AAD Hash and Encrypted note in DB map

			if err = dbInstance.Set(updateNoteInstance.ID, db.StoredData{
				AADData: AAD,
				AADHash: hash,
				Note:    encryptedNote,
			}); err != nil {
				emptyUpdateResponse(w)
				return
			}

			// On success, respond with original ID
			output, _ := json.Marshal(updateNoteResponse{
				ID: updateNoteInstance.ID,
			})

			_, _ = fmt.Fprintf(w, "%+v", string(output))
			return
		}
	}

	// If ID does not exist in DB / pass is incorrect
	// Return with empty ID to indicate error
	emptyUpdateResponse(w)
}

func updateNotePass(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	updateNoteInstance := updateNotePassType{}
	json.Unmarshal(reqBody, &updateNoteInstance)

	// If ID, Pass or newPass are empty or newpass is the same as old pass
	// Return with empty ID To indicate failure
	if updateNoteInstance.ID == "" || updateNoteInstance.Pass == "" || updateNoteInstance.NewPass == "" || updateNoteInstance.NewPass == updateNoteInstance.Pass {
		emptyUpdateResponse(w)
		return
	}

	// Check if there is any data with supplied ID
	if exists, err := dbInstance.Exists(updateNoteInstance.ID); exists && err == nil {
		// If note with ID exists in DB, verify its password
		// And take decrypted AAD
		data, err := dbInstance.Get(updateNoteInstance.ID)
		if err != nil {
			emptyUpdateResponse(w)
			return
		}
		AAD, err := verifyNotePassword(data, updateNoteInstance.Pass)
		if err == nil {
			// If verfication was successful, decrypt the note and ecrypt it with new pass
			AAD, hash, encryptedNote := encrypt(decrypt(data, AAD), updateNoteInstance.NewPass)

			// Save new AAD, AAD Hash and Encrypted note in DB map
			dbInstance.Set(updateNoteInstance.ID, db.StoredData{
				AADData: AAD,
				AADHash: hash,
				Note:    encryptedNote,
			})

			// On success, respond with original ID
			output, _ := json.Marshal(updateNotePassResponse{
				ID: updateNoteInstance.ID,
			})

			_, _ = fmt.Fprintf(w, "%+v", string(output))
			return
		}
	}

	// If ID does not exist in DB / pass is incorrect
	// Return with empty ID to indicate error
	emptyUpdateResponse(w)
}

func emptyUpdateResponse(w http.ResponseWriter) {
	output, _ := json.Marshal(updateNotePassResponse{
		ID: "",
	})
	_, _ = fmt.Fprintf(w, "%+v", string(output))
}
