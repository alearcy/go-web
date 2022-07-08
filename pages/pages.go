package pages

import (
	"net/http"
	"web/utils"
)

// Index is the main web page
func Index(w http.ResponseWriter, r *http.Request) {
	utils.GenerateHTML(w, r, nil, "layout", "index")
}

// Dashboard is where user can see his data
func Dashboard(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("session")
	utils.GenerateHTML(w, r, c.Value, "layout", "dashboard")
}
