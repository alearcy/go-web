package utils

import (
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"text/template"
	"time"
)

const (
	layoutsDir   = "templates/layouts"
	templatesDir = "templates"
	extension    = "/*.gohtml"
)

var pages map[string]*template.Template

// GenerateTemplatesFromFiles generate templates based on data and html
func GenerateTemplatesFromFiles(files embed.FS) error {
	if pages == nil {
		pages = make(map[string]*template.Template)
	}
	templates, err := fs.ReadDir(files, templatesDir)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, tmpl := range templates {
		if tmpl.IsDir() {
			continue
		}
		pt, err := template.ParseFS(files, templatesDir+"/"+tmpl.Name(), layoutsDir+extension)
		if err != nil {
			fmt.Println(err)
			return err
		}
		pages[tmpl.Name()] = pt
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

// Execute the template based on the filename and data
func ExecTemplate(w http.ResponseWriter, r *http.Request, data any, fileName string) error {

	// TODO: mappare errore
	flashMessages, _ := ShowFlash(w, r)

	dataWithFlashMsgs := map[string]any{
		"FlashMessages": flashMessages,
		"Data":          data,
	}
	t, ok := pages[fileName]
	if !ok {
		return errors.New("template not found")
	}
	if err := t.Execute(w, dataWithFlashMsgs); err != nil {
		return errors.New("error executing template")
	}
	return nil
}
