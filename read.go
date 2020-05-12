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

	// If ID or password are empty, return with empty ID to indicate failure
	if getDataInstance.ID == "" || getDataInstance.Pass == "" {
		output, _ := json.Marshal(getDataResponse{
			ID:   "",
			Note: "",
		})
		_, _ = fmt.Fprintf(w, "%+v", string(output))
		return
	}

	// Check if there is any data with supplied ID
	if !storedDataEmpty(db[getDataInstance.ID]) {
		// If so, verify pass and store decrypted AAD
		AAD, err := verifyNotePassword(db[getDataInstance.ID], getDataInstance.Pass)
		if err == nil {
			// Verification successful, decrypt data and send response with ID
			output, _ := json.Marshal(getDataResponse{
				ID:   getDataInstance.ID,
				Note: decrypt(db[getDataInstance.ID], AAD),
			})

			_, _ = fmt.Fprintf(w, "%+v", string(output))
			return
		}
	}

	// If ID does not exist in DB / pass is incorrect
	// Return with supplied ID and empty note
	output, _ := json.Marshal(getDataResponse{
		ID:   getDataInstance.ID,
		Note: "",
	})
	_, _ = fmt.Fprintf(w, "%+v", string(output))
}
