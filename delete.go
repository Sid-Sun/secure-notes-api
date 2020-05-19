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
	if exists, err := dbInstance.Exists(deleteNoteInstance.ID); exists && err == nil {
		// If note with ID exists in DB, fetch it
		if data, err := dbInstance.Get(deleteNoteInstance.ID); err == nil {
			// Verify pass supplied is note's pass
			// And throw the AAD - we don't need it to delete note
			if _, err = verifyNotePassword(data, deleteNoteInstance.Pass); err == nil {
				// If verification was successful, delete note
				if err := dbInstance.Delete(deleteNoteInstance.ID); err == nil {
					output, _ := json.Marshal(deleteNoteResponse{
						ID: deleteNoteInstance.ID,
					})
					w.WriteHeader(200)
					_, _ = fmt.Fprintf(w, "%+v", string(output))
				}
			}
		}

	}

	// If ID does not exist in DB / pass is incorrect
	// Return 404
	w.WriteHeader(404)
}
