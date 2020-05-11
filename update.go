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

	// If note or password are empty, return with empty ID
	// To suggest failure
	if updateNoteInstance.Note == "" || updateNoteInstance.Pass == "" || updateNoteInstance.ID == "" {
		output, _ := json.Marshal(updateNoteResponse{
			ID: "",
		})
		_, _ = fmt.Fprintf(w, "%+v", string(output))
		return
	}

	if !storedDataEmpty(db[updateNoteInstance.ID]) {
		_, err := verifyNotePassword(db[updateNoteInstance.ID], updateNoteInstance.Pass)
		if err == nil {
			// If NewPass is supplied, set pass as newpass before encrypting
			if updateNoteInstance.NewPass != "" {
				updateNoteInstance.Pass = updateNoteInstance.NewPass
			}
			AAD, hash, encryptedNote := encrypt(updateNoteInstance.Note, updateNoteInstance.Pass)

			// Save AAD, AAD Hash and Encrypted note in db map
			db[updateNoteInstance.ID] = storedData{
				AADData: AAD,
				AADHash: hash,
				Note:    encryptedNote,
			}

			// On success, respond with proper ID
			output, _ := json.Marshal(updateNoteResponse{
				ID: updateNoteInstance.ID,
			})

			_, _ = fmt.Fprintf(w, "%+v", string(output))
			return
		}
	}

	// If ID does not exist in DB / pass is incorrect
	// Return with empty ID to indicate error
	output, _ := json.Marshal(updateNoteResponse{
		ID: "",
	})

	_, _ = fmt.Fprintf(w, "%+v", string(output))
}

func updateNotePass(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	updateNoteInstance := updateNotePassType{}
	json.Unmarshal(reqBody, &updateNoteInstance)

	// If ID or password are empty, return with empty ID
	// To indicate failure
	if updateNoteInstance.ID == "" || updateNoteInstance.Pass == "" || updateNoteInstance.NewPass == "" || updateNoteInstance.NewPass == updateNoteInstance.Pass {
		output, _ := json.Marshal(updateNotePassResponse{
			ID: "",
		})
		_, _ = fmt.Fprintf(w, "%+v", string(output))
		return
	}

	// Check if there is any data with supplied ID
	if !storedDataEmpty(db[updateNoteInstance.ID]) {
		AAD, err := verifyNotePassword(db[updateNoteInstance.ID], updateNoteInstance.Pass)
		if err == nil {
			AAD, hash, encryptedNote := encrypt(decrypt(db[updateNoteInstance.ID], AAD), updateNoteInstance.NewPass)

			// Save AAD, AAD Hash and Encrypted note in db map
			db[updateNoteInstance.ID] = storedData{
				AADData: AAD,
				AADHash: hash,
				Note:    encryptedNote,
			}

			// Verification successful, decrypt data and send response
			output, _ := json.Marshal(updateNotePassResponse{
				ID: updateNoteInstance.ID,
			})
			_, _ = fmt.Fprintf(w, "%+v", string(output))
			return
		}
	}

	// If ID does not exist in DB / pass is incorrect
	// Return with empty ID to indicate error
	output, _ := json.Marshal(updateNotePassResponse{
		ID: "",
	})
	_, _ = fmt.Fprintf(w, "%+v", string(output))
}
