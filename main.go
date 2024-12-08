package main

import (
    "fmt"
	"log"
	"net/http"
    "strconv"
)

// define a home handler function twhich writes a byte slice containing
// "Oop Ack Blorff" as the response body
func home(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Oop Ack Blorff"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }

    msg := fmt.Sprintf("Oop greeble bork %d", id)
    w.Write([]byte(msg))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("CREAT"))
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("save muh new snibbage"))
}

func main() {
	// Use the http.NewServeMux() function to init a new servermux,
	// then register the home function as the handler for "/"
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home) // restrict route to matches on / only
    mux.HandleFunc("GET /snippet/view/{id}", snippetView)
    mux.HandleFunc("GET /snippet/create", snippetCreate)
    mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Print("starting server on :4000")

	// ListenAndServe to start a new web clerver. Give it an addres to
	// listen on and a servermux
	// 

	err := http.ListenAndServe("localhost:4000", mux)
	log.Fatal(err)
}
