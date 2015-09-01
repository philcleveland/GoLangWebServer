package main

import (
	"html/template"
	"net/http"
	"os"
	"log"
	"path"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
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
	db gorm.DB
)

func main() {
	var err error
	db, err = gorm.Open("sqlite3", "./characters.db")

	if err != nil {
		log.Fatal(err)
	}

	// defer db.Close()
	db.LogMode(true)

	db.DropTable("character")
	db.CreateTable(&Character{})

	character := Character{
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

	if db.NewRecord(&character){
		db.Create(character)
	}

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
	var characterData Character
	db.First(&characterData)
	err := templates.ExecuteTemplate(w, "index.html", characterData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
