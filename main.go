package main

import (
	"log"
	"net/http"
	"web/auth"
	"web/database"
	model "web/models"
	"web/pages"
)

func init() {
	database.StartDb()
}

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", pages.Index)
	mux.HandleFunc("/admin", auth.Protected(pages.Admin))
	mux.HandleFunc("/users", auth.Protected(model.GetUsers))
	mux.HandleFunc("/signup", auth.Signup)
	mux.HandleFunc("/login", auth.Login)
	mux.HandleFunc("/logout", auth.Logout)
	mux.HandleFunc("/404", pages.NotFound)
	mux.HandleFunc("/500", pages.ServerError)
	log.Println("Listening on localhost:8000...")
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
