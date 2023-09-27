package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
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

func FormatURL(r *http.Request) (string, error) {
	split := strings.Split(r.URL.Path, string(os.PathSeparator))
	fmt.Println(split)
	if len(split) > 0 {
		if split[0] == "/" {

			return "/index", nil
		}
		return split[0], nil

	}
	return "", errors.New("URL non valido")
}

// GenerateTemplate generate HTML templates based on URL with data
func GenerateTemplate(w http.ResponseWriter, r *http.Request, data any, fileName string) {
	flashMessages, _ := ShowFlash(w, r)
	dataWithFlashMsgs := map[string]any{
		"FlashMessages": flashMessages,
		"Data":          data,
	}
	urlWithoutLastSlash := strings.TrimSuffix(fileName, "/")
	formatteLayoutUrl := "templates/base.gohtml"
	formatteTemplateUrl := fmt.Sprintf("templates%s.gohtml", urlWithoutLastSlash)
	parsedTemplate, err := template.ParseFiles(formatteTemplateUrl, formatteLayoutUrl)
	if err != nil {
		slog.Error(err.Error())
		http.Redirect(w, r, "/404", http.StatusFound)
		return
	}
	templates := template.Must(parsedTemplate, err)
	err = templates.Execute(w, dataWithFlashMsgs)
	if err != nil {
		slog.Error(err.Error())
		http.Redirect(w, r, "/404", http.StatusFound)
		return
	}
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
