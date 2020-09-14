package main

import (
	_ "fmt"
	"golang.org/x/crypto/bcrypt"
	_ "log"
	"net/http"
	_ "time"
)

// User variabili della struct in maiuscolo per renderle esportabili
type User struct {
	ID       int
	Name     string
	Surname  string
	Email    string
	Password []byte
	Role     int
}

func signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		//controllare se l'utente esiste prima di continuare
		rows, err := db.Query("SELECT * FROM users where email = $1", email)
		if err != nil {
			flash(w, "Non è stato possibile interrogare il DB, riprova.")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		if rows.Next() {
			flash(w, "Utente già esistente")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		name := r.FormValue("name")
		surname := r.FormValue("surname")
		password := r.FormValue("password")
		password2 := r.FormValue("password2")
		if password != password2 {
			flash(w, "Le due password non coincidono")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u := User{Name: name, Surname: surname, Email: email, Password: cryptedPassword, Role: 1}
		_, err = db.Exec("insert into users (name, surname, email, password, role) values (?, ?, ?, ?, ?)", &u.Name, &u.Surname, &u.Email, &u.Password, &u.Role)
		if err != nil {
			flash(w, "Non è stato possibile salvare a DB, riprova.")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/users", http.StatusSeeOther)
		return
	}
	if !isLoggedIn(w, r) {
		generateHTML(w, r, nil, "layout", "signup")
		return
	}
	flash(w, "Risulti già loggato!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		remember := r.FormValue("remember")
		//controllare se l'utente esiste prima di continuare
		rows, err := db.Query("SELECT * FROM users where email = $1", email)
		if err != nil {
			flash(w, "Non è stato possibile interrogare il DB, riprova.")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		if rows.Next() {
			//e redirect su /dashboard
			u := User{}
			err := rows.Scan(&u.ID, &u.Name, &u.Surname, &u.Email, &u.Password, &u.Role)
			if err != nil {
				flash(w, "Non è stato possibile interrogare il DB, riprova.")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = bcrypt.CompareHashAndPassword(u.Password, []byte(password))
			if err != nil {
				flash(w, "Email o password non validi")
				http.Error(w, "Email o password non validi", http.StatusForbidden)
				return
			}
			uuid, _ := generateUUID()
			if ok, _ := generateSession(w, u, remember, uuid); ok {
				http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			}
		} else {
			flash(w, "Email o password non validi")
			http.Error(w, "Email o password non validi", http.StatusForbidden)
			return
		}
	}
	if isLoggedIn(w, r) {
		generateHTML(w, r, nil, "layout", "dashboard")
		return
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err := deleteCookie(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func protected(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok := isLoggedIn(w, r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}
		h(w, r)
	}
}
