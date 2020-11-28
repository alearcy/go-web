package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"web/auth"
	"web/database"
	"web/models/users"
	"web/pages"
)

func init() {
	database.StartDb()
}

func main() {
	files := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", files))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	r := mux.NewRouter()
	r.HandleFunc("/", pages.Index)
	r.HandleFunc("/dashboard", auth.Protected(pages.Dashboard))
	r.HandleFunc("/users", auth.Protected(users.GetUsers))
	r.HandleFunc("/users/{id:[0-9]+}", auth.Protected(users.GetUser))
	r.HandleFunc("/users/create", auth.Signup)
	r.HandleFunc("/login", auth.Login)
	r.HandleFunc("/logout", auth.Logout)
	log.Println("Listening on localhost:8000...")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal(err)
	}
}
