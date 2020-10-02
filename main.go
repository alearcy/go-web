package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	_ "time"
)

var db *sql.DB

func index(w http.ResponseWriter, r *http.Request) {
	music := []string{"pop", "rock", "rap", "metal", "classical"}
	generateHTML(w, r, music, "layout", "index", "partial")
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("session")
	generateHTML(w, r, c.Value, "layout", "dashboard")
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	urs := make([]User, 0)
	for rows.Next() {
		us := User{}
		err := rows.Scan(&us.ID, &us.Name, &us.Surname, &us.Email, &us.Password, &us.Role)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		urs = append(urs, us)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	generateHTML(w, r, urs, "layout", "users")
}

func getUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	params := mux.Vars(r)
	param := params["id"]

	row := db.QueryRow("SELECT name, surname, email, role FROM users where id = $1", param)

	us := User{}
	err := row.Scan(&us.Name, &us.Surname, &us.Email, &us.Role)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		fmt.Println(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	generateHTML(w, r, us, "layout", "user")
}

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://arcy:Aleedenny10@localhost/go?sslmode=disable")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to the DB")
}

func main() {
	files := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", files))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/dashboard", protected(dashboard))
	r.HandleFunc("/users", protected(getUsers))
	r.HandleFunc("/users/{id:[0-9]+}", protected(getUser))
	r.HandleFunc("/users/create", signup)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	log.Println("Listening on :8000...")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal(err)
	}
}
