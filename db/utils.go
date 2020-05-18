package db

import "bytes"

func storedDataEmpty(a StoredData) bool {
	// If any of the fields are empty
	// The data is considered as empty
	if bytes.Equal(a.AADData, []byte{}) {
		return true
	}
	if bytes.Equal(a.AADHash[:], []byte{}) {
		return true
	}
	if bytes.Equal(a.Note, []byte{}) {
		return true
	}
	return false
}
