package utils

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	LayoutTemplatesFolder = "templates/layouts"
	AdminTemplatesFolder  = "templates/admin"
	TemplatesFolder       = "templates"
)

// GenerateTemplate generate templates based on data and html
func GenerateTemplate(w http.ResponseWriter, r *http.Request, data any, fileName string) error {
	var layoutType string
	flashMessages, _ := ShowFlash(w, r)
	dataWithFlashMsgs := map[string]interface{}{
		"FlashMessages": flashMessages,
		"Data":          data,
	}

	split := strings.Split(r.URL.Path, string(os.PathSeparator))
	level := split[1]
	fmt.Println(level)
	if level != "" {
		layoutType = fmt.Sprintf("templates/%s/", level)
	} else {
		layoutType = "templates"
	}
	formatteLayoutdUrl := fmt.Sprintf("%s/base.gohtml", layoutType)
	formatteTemmplateUrl := fmt.Sprintf("%s/%s.gohtml", layoutType, fileName)
	templates := template.Must(template.ParseFiles(formatteTemmplateUrl, formatteLayoutdUrl))
	fmt.Println(templates.Name())
	err := templates.Execute(w, dataWithFlashMsgs)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
	return nil
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

// ShowFlash - show a flash message passing the ResponseWriter and the Request pointer
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
