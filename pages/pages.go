package pages

import (
	"fmt"
	"net/http"
	"web/utils"
)

// Index is the main web page
func Index(w http.ResponseWriter, r *http.Request) {
	if err := utils.ExecTemplate(w, r, nil, "index.gohtml"); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// Dashboard is where user can see his data
func Dashboard(w http.ResponseWriter, r *http.Request) {
	// TODO: creare una data struct con cookie e T per data
	c, _ := r.Cookie("session")
	utils.ExecTemplate(w, r, c, "dashboard.gohtml")
}
