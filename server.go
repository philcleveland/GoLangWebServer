package main

import (
	"html/template"
	"net/http"
	"os"
	"path"
)

const (
	//StaticURL is the root of all static content
	StaticURL string = "/static/"
)

//Context is data descibing a page title and location of static content
type Context struct {
	Title  string
	Static string
}

var (
	templates *template.Template
)

func main() {
	rootdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	templates = template.Must(template.ParseFiles("templates/index.html"))

	// This is the only way I have found to be able to serve images requested in the templates
	http.Handle("/static/img/", http.StripPrefix("/static/img/",
		http.FileServer(http.Dir(path.Join(rootdir, "/static/img/")))))

	http.Handle("/static/css/", http.StripPrefix("/static/css/",
		http.FileServer(http.Dir(path.Join(rootdir, "/static/css/")))))

	http.HandleFunc("/", index)

	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	context := Context{Title: "Char Sheet", Static: StaticURL + "img/"}
	err := templates.ExecuteTemplate(w, "index.html", context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
