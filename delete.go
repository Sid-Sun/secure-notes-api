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
		emptyDeleteResponse(w)
		return
	}

	// Check if there is any data with supplied ID
	if exists, err := existsInDB(deleteNoteInstance.ID); exists {
		if err != nil {
			emptyDeleteResponse(w)
			return
		}
		data, err := fetchFromDB(deleteNoteInstance.ID)
		if err != nil {
			emptyDeleteResponse(w)
			return
		}
		// If note with ID exists in DB, verify its password
		// And throw the AAD - we don't need it to delete note
		_, err = verifyNotePassword(data, deleteNoteInstance.Pass)
		if err == nil {
			// If verification was successful, replace note with empty zer-valued storedData
			err = deleteFromDB(deleteNoteInstance.ID)
			if err != nil {
				emptyDeleteResponse(w)
				return
			}

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
	emptyDeleteResponse(w)
}

func emptyDeleteResponse(w http.ResponseWriter) {
	output, _ := json.Marshal(deleteNoteResponse{
		ID: "",
	})

	_, _ = fmt.Fprintf(w, "%+v", string(output))
}
