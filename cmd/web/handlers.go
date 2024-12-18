package main

import (
    "fmt"
    "html/template"
    "net/http"
    "strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Server", "FORTRAN")

    files := []string {
        "./ui/html/base.tmpl",
            "./ui/html/partials/nav.tmpl",
            "./ui/html/pages/home.tmpl",
        }

    // read source file into a template set
    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.logger.Error(err.Error(), "method", r.Method, "url", r.URL.RequestURI())
        http.Error(w, "infernal server error", http.StatusInternalServerError)
        return
    }

    // then use execute on the template set to write to as
    // the respons body.  Last parameter is any dymanic data 
    // gets passed in
    err = ts.ExecuteTemplate(w, "base", nil)
    if err != nil {
        app.logger.Error(err.Error(), "method", r.Method, "url", r.URL.RequestURI())
        http.Error(w, "Infernal server error", http.StatusInternalServerError)
    }
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }
    fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Display a form for creating a new snippet..."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Save a new snippet..."))
}
