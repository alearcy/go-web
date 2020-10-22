package pages

import (
	"net/http"
	"web/utils"
)

func Index(w http.ResponseWriter, r *http.Request) {
	music := []string{"pop", "rock", "rap", "metal", "classical"}
	utils.GenerateHTML(w, r, music, "layout", "index", "partial")
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("session")
	utils.GenerateHTML(w, r, c.Value, "layout", "dashboard")
}
