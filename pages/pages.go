package pages

import (
	"net/http"
	"web/utils"
)

// Index is the main web page
func Index(w http.ResponseWriter, r *http.Request) {
	utils.GenerateTemplate(w, r, nil, "index")
}

// NotFound is the custom page when an HTML page was not founded
func NotFound(w http.ResponseWriter, r *http.Request) {
	utils.GenerateTemplate(w, r, nil, "404")
}

// ServerError is the custom page when server complains
func ServerError(w http.ResponseWriter, r *http.Request) {
	utils.GenerateTemplate(w, r, nil, "500")
}

// Dashboard is where user can see his data
func Admin(w http.ResponseWriter, r *http.Request) {
	// TODO: creare una data struct con cookie e T per data
	c, _ := r.Cookie("session")
	utils.GenerateTemplate(w, r, c, "admin")
}
