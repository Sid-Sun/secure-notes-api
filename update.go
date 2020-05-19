package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func updateNote(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	updateNoteInstance := updateNoteType{}
	json.Unmarshal(reqBody, &updateNoteInstance)

	// If note, password or id are empty, return with empty ID To suggest failure
	if updateNoteInstance.Note == "" || updateNoteInstance.Pass == "" || updateNoteInstance.ID == "" {
		w.WriteHeader(400)
		return
	}

	if !storedDataEmpty(db[updateNoteInstance.ID]) {
		// If note with ID exists in DB, verify its password 
		// And throw the AAD - we don't need it to replace note
		_, err := verifyNotePassword(db[updateNoteInstance.ID], updateNoteInstance.Pass)
		if err == nil {
			// If NewPass is supplied, set pass as newpass before encrypting
			if updateNoteInstance.NewPass != "" {
				updateNoteInstance.Pass = updateNoteInstance.NewPass
			}
			AAD, hash, encryptedNote := encrypt(updateNoteInstance.Note, updateNoteInstance.Pass)

			// Save new AAD, AAD Hash and Encrypted note in DB map
			db[updateNoteInstance.ID] = storedData{
				AADData: AAD,
				AADHash: hash,
				Note:    encryptedNote,
			}

			// On success, respond with original ID
			output, _ := json.Marshal(updateNoteResponse{
				ID: updateNoteInstance.ID,
			})
			w.WriteHeader(200)
			_, _ = fmt.Fprintf(w, "%+v", string(output))
			return
		}
	}

	// If ID does not exist in DB / pass is incorrect
	// Return with empty ID to indicate error
	w.WriteHeader(404)
}

func updateNotePass(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	updateNoteInstance := updateNotePassType{}
	json.Unmarshal(reqBody, &updateNoteInstance)

	// If ID, Pass or newPass are empty or newpass is the same as old pass
	// Return with empty ID To indicate failure
	if updateNoteInstance.ID == "" || updateNoteInstance.Pass == "" || updateNoteInstance.NewPass == "" || updateNoteInstance.NewPass == updateNoteInstance.Pass {
		w.WriteHeader(400)
		return
	}

	// Check if there is any data with supplied ID
	if !storedDataEmpty(db[updateNoteInstance.ID]) {
		// If note with ID exists in DB, verify its password 
		// And take decrypted AAD
		AAD, err := verifyNotePassword(db[updateNoteInstance.ID], updateNoteInstance.Pass)
		if err == nil {
			// If verfication was successful, decrypt the note and ecrypt it with new pass
			AAD, hash, encryptedNote := encrypt(decrypt(db[updateNoteInstance.ID], AAD), updateNoteInstance.NewPass)

			// Save new AAD, AAD Hash and Encrypted note in DB map
			db[updateNoteInstance.ID] = storedData{
				AADData: AAD,
				AADHash: hash,
				Note:    encryptedNote,
			}

			// On success, respond with original ID
			output, _ := json.Marshal(updateNotePassResponse{
				ID: updateNoteInstance.ID,
			})
			w.WriteHeader(200)
			_, _ = fmt.Fprintf(w, "%+v", string(output))
			return
		}
	}

	// If ID does not exist in DB / pass is incorrect
	// Return with empty ID to indicate error
	w.WriteHeader(404)
}
