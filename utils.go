package main

import (
	"fmt"
	_ "log"
	"net/http"
	"strings"
	"text/template"
)

// per passare funzioni al template creo una variabile di tipo FuncMap che accetta una chiave stringa e una funzione qualsiasi
var fn = template.FuncMap{
	"uppercase": strings.ToUpper,
}

// generate templates based on data and html
func generateHTML(w http.ResponseWriter, data interface{}, files ...string) {
	var a []string
	for _, f := range files {
		a = append(a, fmt.Sprintf("templates/%s.html", f))
	}
	// template.Must si occupa lui di fare l'error checking senza essere ripetitivi e accetta un template come argomento
	// template.PareGlob prende tutti i template dentro una cartella mentre template.ParseFiles uno alla volta dentro slice
	// template.New mi serve per inizializzare il puntatore a template, passargli le funzioni e fargliele trovare inizializzate ai files .gohtml
	templates := template.Must(template.New("").Funcs(fn).ParseFiles(a...))
	templates.ExecuteTemplate(w, "layout", data)
}


func getParam(r *http.Request, s string) string {
	param := strings.TrimPrefix(r.URL.Path, s)
	return param
}