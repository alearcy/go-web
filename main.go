package main

import (
	"embed"
	"log"
	"net/http"
	"web/auth"
	"web/database"
	model "web/models"
	"web/pages"
	"web/utils"
)

var (
	//go:embed templates/* templates/layouts/*.gohtml
	files embed.FS
)

func init() {
	database.StartDb()
}

func main() {
	utils.GenerateTemplatesFromFiles(files)
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	// mux.Handle("/favicon.ico", http.NotFoundHandler())
	mux.HandleFunc("/", pages.Index)
	mux.HandleFunc("/dashboard", auth.Protected(pages.Dashboard))
	mux.HandleFunc("/users", auth.Protected(model.GetUsers))
	// mux.HandleFunc("/users/{id:[0-9]+}", auth.Protected(model.GetUser))
	mux.HandleFunc("/signup", auth.Signup)
	mux.HandleFunc("/login", auth.Login)
	mux.HandleFunc("/logout", auth.Logout)
	log.Println("Listening on localhost:8000...")
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal(err)
	}
}

// TODO; 404, 500, 401 pages, html errors header
