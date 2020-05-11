package main

type storedData struct {
	AADData []byte
	AADHash [32]byte
	Note    []byte
}

type greetResponse struct {
	Message string
}

type setDataType struct {
	ID   string
	Pass string
	Note string
}

type setDataResponse struct {
	ID string
}

type getDataType struct {
	ID   string
	Pass string
}

type getDataResponse struct {
	ID   string
	Note string
}

type updateNoteType struct {
	ID   string
	Pass string
	Note string
	NewPass string
}

type updateNoteResponse struct {
	ID string
}

type updateNotePassType struct {
	ID      string
	Pass    string
	NewPass string
}

type updateNotePassResponse struct {
	ID string
}

type deleteNoteType struct {
	ID   string
	Pass string
}

type deleteNoteResponse struct {
	ID string
}
