package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getData(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	getDataInstance := getDataType{}
	json.Unmarshal(reqBody, &getDataInstance)

	// If ID or password are empty, return with 400
	if getDataInstance.ID == "" || getDataInstance.Pass == "" {
		w.WriteHeader(400)
		return
	}

	// Check if there is any data with supplied ID
	if exists, err := dbInstance.Exists(getDataInstance.ID); exists && err == nil {
		// If so, fetch data
		if data, err := dbInstance.Get(getDataInstance.ID); err == nil {
			// Verify pass and store decrypted AAD
			if AAD, err := verifyNotePassword(data, getDataInstance.Pass); err == nil {
				// Verification successful, decrypt data and send response with ID
				output, _ := json.Marshal(getDataResponse{
					ID:   getDataInstance.ID,
					Note: decrypt(data, AAD),
				})
				w.WriteHeader(200)
				_, _ = fmt.Fprintf(w, "%+v", string(output))
				return
			}
		}
	}

	// If ID does not exist in DB / pass is incorrect
	// Return 404
	w.WriteHeader(404)
}
