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

	// If ID or password are empty, return with empty ID to indicate failure
	if deleteNoteInstance.ID == "" || deleteNoteInstance.Pass == "" {
		output, _ := json.Marshal(deleteNoteResponse{
			ID: "",
		})

		_, _ = fmt.Fprintf(w, "%+v", string(output))
		return
	}

	// Check if there is any data with supplied ID
	if !storedDataEmpty(db[deleteNoteInstance.ID]) {
		// If note with ID exists in DB, verify its password 
		// And throw the AAD - we don't need it to delete note
		_, err := verifyNotePassword(db[deleteNoteInstance.ID], deleteNoteInstance.Pass)
		if err == nil {
			// If verification was successful, replace note with empty zer-valued storedData
			db[deleteNoteInstance.ID] = storedData{}

			// Respond with original ID to indicate success
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
