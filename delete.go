package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func deleteNote(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	deleteNoteInstance := deleteNoteType{}
	json.Unmarshal(reqBody, &deleteNoteInstance)

	// If ID or password are empty, return with empty ID
	// To indicate failure
	if deleteNoteInstance.ID == "" || deleteNoteInstance.Pass == "" {
		output, _ := json.Marshal(deleteNoteResponse{
			ID: "",
		})
		_, _ = fmt.Fprintf(w, "%+v", string(output))
		return
	}

	// Check if there is any data with supplied ID
	if !storedDataEmpty(db[deleteNoteInstance.ID]) {
		_, err := verifyNotePassword(db[deleteNoteInstance.ID], deleteNoteInstance.Pass)
		if err == nil {
			// Save AAD, AAD Hash and Encrypted note in db map
			db[deleteNoteInstance.ID] = storedData{}

			// Verification successful, decrypt data and send response
			output, _ := json.Marshal(deleteNoteResponse{
				ID: deleteNoteInstance.ID,
			})
			_, _ = fmt.Fprintf(w, "%+v", string(output))
			return
		}
	}

	// If ID does not exist in DB / pass is incorrect
	// Return with empty ID to indicate error
	output, _ := json.Marshal(deleteNoteResponse{
		ID: "",
	})
	_, _ = fmt.Fprintf(w, "%+v", string(output))
}
