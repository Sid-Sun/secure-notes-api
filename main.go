package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var db = make(map[string]storedData)

func main() {
	fmt.Println("Hello, World!")
	port := "8990"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	fmt.Println("Starting on port:", port)

	// Create new router and define routes
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Use(addContentTypeHeader)
	myRouter.HandleFunc("/", greet).Methods("GET")
	myRouter.HandleFunc("/set", setData).Methods("POST")
	myRouter.HandleFunc("/update/note", updateNote).Methods("PUT")
	myRouter.HandleFunc("/update/pass", updateNotePass).Methods("PATCH")
	myRouter.HandleFunc("/delete", deleteNote).Methods("DELETE")
	myRouter.HandleFunc("/get", getData).Methods("GET")

	http.Handle("/", myRouter)
	err := http.ListenAndServe(":"+port, myRouter)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func greet(w http.ResponseWriter, r *http.Request) {
	output, _ := json.Marshal(greetResponse{
		Message: `Hi there, Time-Appropriate Greetings! Please consult my documentation at https://github.com/sid-sun/passwordless-notes-api to use me, developed by Sidharth Soni (Sid Sun) - sid@sidsun.com. Open Sourced under MIT.`,
	})
	w.WriteHeader(200)
	_, _ = fmt.Fprintf(w, "%+v", string(output))
}

// addContentTypeHeader acts as a middleware and adds the Content-Type header to responses
func addContentTypeHeader(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}