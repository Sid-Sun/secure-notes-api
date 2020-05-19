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

	// If note, password or id are empty, return 400
	if updateNoteInstance.Note == "" || updateNoteInstance.Pass == "" || updateNoteInstance.ID == "" {
		w.WriteHeader(400)
		return
	}

	// Check if data with ID exist in DB
	if exists, err := dbInstance.Exists(updateNoteInstance.ID); exists && err == nil {
		// If data exists, fetch it
		if data, err := dbInstance.Get(updateNoteInstance.ID); err == nil {
			// Verify provides password against data
			// And throw the AAD - we don't need it to replace note
			if _, err = verifyNotePassword(data, updateNoteInstance.Pass); err == nil {
				// Verifies successfully, we can replace note!
				// If NewPass is supplied, set pass as newpass before encrypting
				if updateNoteInstance.NewPass != "" {
					updateNoteInstance.Pass = updateNoteInstance.NewPass
				}
				// Encrypt new note with pass
				AAD, hash, encryptedNote := encrypt(updateNoteInstance.Note, updateNoteInstance.Pass)
				// Save new AAD, AAD Hash and Encrypted note in DB
				if err = dbInstance.Set(updateNoteInstance.ID, db.StoredData{
					AADData: AAD,
					AADHash: hash,
					Note:    encryptedNote,
				}); err == nil {
					// If successful, return with original IS
					output, _ := json.Marshal(updateNotePassResponse{
						ID: updateNoteInstance.ID,
					})
					w.WriteHeader(200)
					_, _ = fmt.Fprintf(w, "%+v", string(output))
					return
				}
			}
		}
	}

	// If ID does not exist in DB / pass is incorrect
	// Return 404
	w.WriteHeader(404)
}

func updateNotePass(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	updateNoteInstance := updateNotePassType{}
	json.Unmarshal(reqBody, &updateNoteInstance)

	// If ID, Pass or newPass are empty or newpass is the same as old pass
	// Return 400
	if updateNoteInstance.ID == "" || updateNoteInstance.Pass == "" || updateNoteInstance.NewPass == "" || updateNoteInstance.NewPass == updateNoteInstance.Pass {
		w.WriteHeader(400)
		return
	}

	// Check if there is any data with supplied ID
	if exists, err := dbInstance.Exists(updateNoteInstance.ID); exists && err == nil {
		// If note with ID exists in DB, fetch it
		if data, err := dbInstance.Get(updateNoteInstance.ID); err == nil {
			// If fetches without errors, verify pass provides against data
			if AAD, err := verifyNotePassword(data, updateNoteInstance.Pass); err == nil {
				// If verfication was successful, decrypt the note and ecrypt it with new pass
				AAD, hash, encryptedNote := encrypt(decrypt(data, AAD), updateNoteInstance.NewPass)
				// Save new AAD, AAD Hash and Encrypted note in DB map
				if err = dbInstance.Set(updateNoteInstance.ID, db.StoredData{
					AADData: AAD,
					AADHash: hash,
					Note:    encryptedNote,
				}); err == nil {
					// On success, respond with original ID
					output, _ := json.Marshal(updateNotePassResponse{
						ID: updateNoteInstance.ID,
					})
					w.WriteHeader(200)
					_, _ = fmt.Fprintf(w, "%+v", string(output))
					return
				}
			}
		}
	}

	// If ID does not exist in DB / pass is incorrect
	// Return 404
	w.WriteHeader(404)
}
