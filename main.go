package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
	// "database/sql"
	// _ "github.com/mattn/go-sqlite3"
)

// variabili della struct in maiuscolo per renderle esportabili
type user struct {
	Email    string
	Password []byte
	Role     string
}

type session struct {
	uID string
	activityTime time.Time
}

// simulo due DB, uno di sessione e uno degli utenti
var dbUsers = map[string]user{}      // user ID, user
var dbSessions = map[string]session{} // session ID, user ID
var dbSessionsCleaned time.Time

// per passare funzioni al template creo una variabile di tipo FuncMap che accetta una chiave stringa e una funzione qualsiasi
var fn = template.FuncMap{
	"uppercase":  strings.ToUpper,
	"firstThree": firstThree,
}

const sessionLength int = 30

func firstThree(s string) string {
	s = strings.TrimSpace(s)
	s = s[:3]
	return s
}

var tpl *template.Template

// funzione di inizializzazione
func init() {
	// template.Must si occupa lui di fare l'error checking senza essere ripetitivi e accetta un template come argomento
	// template.PareGlob prende tutti i template dentro una cartella
	// template.New mi serve per inizializzare il puntatore a template, passargli le funzioni e fargliele trovare inizializzate ai files .gohtml
	tpl = template.Must(template.New("").Funcs(fn).ParseGlob("templates/*.gohtml"))
}

func index(w http.ResponseWriter, r *http.Request) {

	// creo una variabile da passare al template
	music := []string{"pop", "rock", "rap", "metal", "classical"}

	// tpl.ExecuteTemplate esegue il template
	tpl.ExecuteTemplate(w, "index.gohtml", music)
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	log.Println(c)
	tpl.ExecuteTemplate(w, "dashboard.gohtml", c.Value)
}

func checkForm(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("fname")
	f, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	log.Println("il mio nome: ", name, "il mio file: ", f)
	http.Redirect(w, r, "/read-cookie", http.StatusSeeOther)
}

func main() {
	// con HandleFunc ottengo un codice più pulito in quanto si occupa lui, una volta passata una funzione con argomenti (res http.ResponseWriter, req *http.Request),
	// di creare il multiplexer
	http.HandleFunc("/", index)
	http.HandleFunc("/check-form", checkForm)
	http.HandleFunc("/dashboard", dashboard)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/logout", logout)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8000", nil)
}
