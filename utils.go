package main

import (
	"text/template"
	"net/http"
	"strings"
	"log"
	"fmt"
)

// per passare funzioni al template creo una variabile di tipo FuncMap che accetta una chiave stringa e una funzione qualsiasi
var fn = template.FuncMap{
	"uppercase":  strings.ToUpper,
}

func getParams(r *http.Request) []string {
	keys := strings.Split(r.URL.Path, "/")
	log.Println(r.URL.Path)
        if len(keys) > 0 {
		params := make([]string, 0)
		for _, v := range keys {
			log.Println(v + "\n")
			params = append(params, v)
		}
        	return params
	}
	return nil
}

// generate templates based on data and html
func generateHTML(w http.ResponseWriter, data interface{}, files ...string) {
	var a []string
	for _, f := range files {
		a = append(a, fmt.Sprintf("templates/%s.gohtml", f))
	}
	// template.Must si occupa lui di fare l'error checking senza essere ripetitivi e accetta un template come argomento
	// template.PareGlob prende tutti i template dentro una cartella mentre template.ParseFiles uno alla volta dentro slice
	// template.New mi serve per inizializzare il puntatore a template, passargli le funzioni e fargliele trovare inizializzate ai files .gohtml
	templates := template.Must(template.New("").Funcs(fn).ParseFiles(a...))
	templates.ExecuteTemplate(w, "layout", data)
}