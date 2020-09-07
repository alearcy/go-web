package main

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

// User variabili della struct in maiuscolo per renderle esportabili
type User struct {
	ID       string
	Name     string
	Surname  string
	Email    string
	Password []byte
	Role     int
}

func signup(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(r) {
		tpl.ExecuteTemplate(w, "signup.gohtml", nil)
		return
	}
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		// TODO: controllare che le due passwords coincidino
		// TODO: controllare se l'utente esiste prima di continuare
		password := r.FormValue("password")
		cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u := User{"", "", "", email, cryptedPassword, 1}
		uuid, _ := generateUUID()
		generateCookie(w, uuid)

		// momentanemanete sostituisce il DB
		dbSessions[uuid] = session{email, time.Now()}
		dbUsers[email] = u

		log.Println(u, uuid)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
	log.Println("Sei già loggato")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func login(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(r) {
		tpl.ExecuteTemplate(w, "dashboard.gohtml", nil)
		return
	}
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		u, ok := dbUsers[email]
		if !ok {
			http.Error(w, "Email o password non validi", http.StatusForbidden)
			return
		}
		err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
		if err != nil {
			http.Error(w, "Email o password non validi", http.StatusForbidden)
			return
		}
		uuid, _ := generateUUID()
		c := generateCookie(w, uuid)

		// momentanemanete sostituisce il DB
		dbSessions[c.Value] = session{email, time.Now()}

		log.Println(u, uuid)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	_, err := deleteCookie(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// clean up dbSessions
	if time.Now().Sub(dbSessionsCleaned) > (time.Second * 30) {
		go cleanSessions()
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
