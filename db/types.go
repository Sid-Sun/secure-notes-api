package db

// StoredData species the data format stored in Db
type StoredData struct {
	AADData []byte
	AADHash [32]byte
	Note    []byte
}
