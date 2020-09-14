package main

import (
	"encoding/base64"
	"fmt"
	_ "log"
	"net/http"
	"strings"
	"text/template"
	"time"
)

// per passare funzioni al template creo una variabile di tipo FuncMap che accetta una chiave stringa e una funzione qualsiasi
var fn = template.FuncMap{
	"uppercase": strings.ToUpper,
}

// generate templates based on data and html
func generateHTML(w http.ResponseWriter, r *http.Request, data interface{}, files ...string) {
	flashMessages, _ := showFlash(w, r)
	dataWithFlashMsgs := map[string]interface{}{
		"FlashMessages": flashMessages,
		"Data":          data,
	}

	var a []string
	for _, f := range files {
		a = append(a, fmt.Sprintf("templates/%s.html", f))
	}
	// template.Must si occupa lui di fare l'error checking senza essere ripetitivi e accetta un template come argomento
	// template.PareGlob prende tutti i template dentro una cartella mentre template.ParseFiles uno alla volta dentro slice
	// template.New mi serve per inizializzare il puntatore a template, passargli le funzioni e fargliele trovare inizializzate ai files .gohtml
	templates := template.Must(template.New("").Funcs(fn).ParseFiles(a...))
	templates.ExecuteTemplate(w, "layout", dataWithFlashMsgs)
}

func getParam(r *http.Request, toStrip string) string {
	param := strings.TrimPrefix(r.URL.Path, toStrip)
	return param
}

func flash(w http.ResponseWriter, s string) {
	msg := []byte(s)
	c := http.Cookie{
		Name:  "flash",
		Value: base64.URLEncoding.EncodeToString(msg),
		Path: "/",
	}
	http.SetCookie(w, &c)
}

func showFlash(w http.ResponseWriter, r *http.Request) (string, error) {
	// TODO: meglio usare questo perché potrei avere più flash messages
	// quindi anche nel template meglio usare un ciclo per mostrare i messaggi
	// for _, cookie := range r.Cookies() {
	// 	fmt.Fprint(w, cookie.Name)
	// }
	c, err := r.Cookie("flash")
	var val []byte
	if err != nil {
		return "", err
	}
	rc := http.Cookie{
		Name:    "flash",
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
		Path: "/",
	}
	http.SetCookie(w, &rc)
	val, _ = base64.URLEncoding.DecodeString(c.Value)
	return string(val), nil
}
