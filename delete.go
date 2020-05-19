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

	// If ID or password are empty, return 400
	if deleteNoteInstance.ID == "" || deleteNoteInstance.Pass == "" {
		w.WriteHeader(400)
		return
	}

	// Check if there is any data with supplied ID
	if exists, err := existsInDB(deleteNoteInstance.ID); exists && err == nil {
		// If so, fetch data
		var data storedData
		if data, err = fetchFromDB(deleteNoteInstance.ID); err == nil {
			// Verify its password and throw the AAD - we don't need it to delete note
			if _, err = verifyNotePassword(data, deleteNoteInstance.Pass); err == nil {
				// If verification was successful, delete
				if err = deleteFromDB(deleteNoteInstance.ID); err == nil {
					// Respond with original ID to indicate success
					output, _ := json.Marshal(deleteNoteResponse{
						ID: deleteNoteInstance.ID,
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
