package utils

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"
)

// Fn is a map with useful template function
var Fn = template.FuncMap{
	"uppercase": strings.ToUpper,
}

// GenerateHTML generate templates based on data and html
func GenerateHTML(w http.ResponseWriter, r *http.Request, data interface{}, files ...string) {
	flashMessages, _ := ShowFlash(w, r)
	dataWithFlashMsgs := map[string]interface{}{
		"FlashMessages": flashMessages,
		"Data":          data,
	}

	var a []string
	for _, f := range files {
		a = append(a, fmt.Sprintf("templates/%s.gohtml", f))
	}
	// template.Must si occupa lui di fare l'error checking senza essere ripetitivi e accetta un template come argomento
	// template.PareGlob prende tutti i template dentro una cartella mentre template.ParseFiles uno alla volta dentro slice
	// template.New mi serve per inizializzare il puntatore a template, passargli le funzioni e fargliele trovare inizializzate ai files .gohtml
	templates := template.Must(template.New("").Funcs(Fn).ParseFiles(a...))
	templates.ExecuteTemplate(w, "layout", dataWithFlashMsgs)
}

// Flash - create flash message passing ResponseWriter and a message
func Flash(w http.ResponseWriter, s string) {
	msg := []byte(s)
	c := http.Cookie{
		Name:  "flash",
		Value: base64.URLEncoding.EncodeToString(msg),
		Path:  "/",
	}
	http.SetCookie(w, &c)
}

// ShowFlash - show a flash message passing the ResponseWriter and the Requiest pointer
func ShowFlash(w http.ResponseWriter, r *http.Request) (string, error) {
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
		Path:    "/",
	}
	http.SetCookie(w, &rc)
	val, _ = base64.URLEncoding.DecodeString(c.Value)
	return string(val), nil
}
