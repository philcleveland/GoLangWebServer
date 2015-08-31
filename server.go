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
	character := Character{Static: StaticURL + "img/",
		CharacterName: "LionHeart",
		ClassLevel:    "Fighter 1",
		Background:    "Folk Hero",
		PlayerName:    "Greg",
		Faction:       "Harpers",
		Race:          "Human",
		Alignment:     "Lawful Good",
		XP:            0,
		DCI:           12345,
		Strength:      15,
		Dexterity:     13,
		Constitution:  14,
		Intelligence:  8,
		Wisdom:        12,
		Charisma:      10}
	err := templates.ExecuteTemplate(w, "index.html", character)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
