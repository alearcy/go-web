package pages

import (
	"log/slog"
	"net/http"
	"web/pkg/utils"
)

// Index is the main web page
func Index(w http.ResponseWriter, r *http.Request) {
	url, err := utils.FormatURL(r)
	if err != nil {
		slog.Error(err.Error())
		utils.GenerateTemplate(w, r, nil, "/404")
		return
	}
	utils.GenerateTemplate(w, r, nil, url)
}

// NotFound is the custom page when an HTML page was not founded
func NotFound(w http.ResponseWriter, r *http.Request) {
	url, _ := utils.FormatURL(r)
	utils.GenerateTemplate(w, r, nil, url)
}

// ServerError is the custom page when server complains
func ServerError(w http.ResponseWriter, r *http.Request) {
	url, err := utils.FormatURL(r)
	if err != nil {
		slog.Error(err.Error())
		http.Redirect(w, r, "/404", http.StatusOK)
		return
	}
	utils.GenerateTemplate(w, r, nil, url)
}

// Dashboard is where user can see his data
func Admin(w http.ResponseWriter, r *http.Request) {
	// TODO: creare una data struct con cookie e T per data
	c, _ := r.Cookie("session")
	url, err := utils.FormatURL(r)
	if err != nil {
		slog.Error(err.Error())
		http.Redirect(w, r, "/404", http.StatusOK)
		return
	}
	utils.GenerateTemplate(w, r, c, url)
}
