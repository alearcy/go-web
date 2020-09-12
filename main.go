package main

import (
	"database/sql"
	"fmt"
	_ "github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"time"
)

// simulo due DB, uno di sessione e uno degli utenti
var dbUsers = map[string]User{}       // user ID, user
var dbSessions = map[string]session{} // session ID, user ID
var dbSessionsCleaned time.Time

const sessionLength int = 30

var tpl *template.Template
var db *sql.DB

func index(w http.ResponseWriter, r *http.Request) {
	// creo un cookie con un messaggio flash
	flash(w, "Benvenuto")
	// creo una variabile da passare al template
	music := []string{"pop", "rock", "rap", "metal", "classical"}
	generateHTML(w, r, music, "layout", "index", "partial")

}

func dashboard(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	generateHTML(w, r, c.Value, "layout", "dashboard")
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	urs := make([]User, 0)
	for rows.Next() {
		us := User{}
		err := rows.Scan(&us.ID, &us.Name, &us.Surname, &us.Email, &us.Password, &us.Role)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		urs = append(urs, us)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	generateHTML(w, r, urs, "layout", "users")
}

func getUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	param := getParam(r, "/user/")

	row := db.QueryRow("SELECT * FROM users where id = $1", param)

	us := User{}
	err := row.Scan(&us.ID, &us.Name, &us.Surname, &us.Email, &us.Password, &us.Role)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	generateHTML(w, r, us, "layout", "user")
}

// funzione di inizializzazione
func init() {
	var err error
	// inizializzo DB
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
	// carico gli assets statici
	files := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", files))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	// con HandleFunc ottengo un codice più pulito in quanto si occupa lui, una volta passata una funzione con argomenti (res http.ResponseWriter, req *http.Request),
	// di creare il multiplexer
	http.HandleFunc("/", index)
	http.HandleFunc("/dashboard/", protected(dashboard))
	http.HandleFunc("/users/", protected(getUsers))
	http.HandleFunc("/user/", protected(getUser))
	http.HandleFunc("/signup/", signup)
	http.HandleFunc("/logout/", logout)
	log.Println("Listening on :8000...")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
