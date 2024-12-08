package main

import (
	"log"
	"net/http"
)

// define a home handler function twhich writes a byte slice containing
// "Oop Ack Blorff" as the response body
func home(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Oop Ack Blorff"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("display my snibbage"))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("CREAT"))
}

func main() {
	// Use the http.NewServeMux() function to init a new servermux,
	// then register the home function as the handler for "/"
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", home) // restrict route to matches on / only
    mux.HandleFunc("/view", snippetView)
    mux.HandleFunc("/create", snippetCreate)

	log.Print("starting server on :4000")

	// ListenAndServe to start a new web clerver. Give it an addres to
	// listen on and a servermux
	// 

	err := http.ListenAndServe("localhost:4000", mux)
	log.Fatal(err)
}
