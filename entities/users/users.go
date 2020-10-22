package users

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"web/auth"
	"web/database"
	"web/utils"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	rows, err := database.Db.Query("SELECT * FROM users")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	urs := make([]auth.User, 0)
	for rows.Next() {
		us := auth.User{}
		err := rows.Scan(&us.ID, &us.Name, &us.Surname, &us.Email, &us.Password, &us.Role)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		urs = append(urs, us)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	utils.GenerateHTML(w, r, urs, "layout", "users")
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	params := mux.Vars(r)
	param := params["id"]

	row := database.Db.QueryRow("SELECT name, surname, email, role FROM users where id = $1", param)

	us := auth.User{}
	err := row.Scan(&us.Name, &us.Surname, &us.Email, &us.Role)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		fmt.Println(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	utils.GenerateHTML(w, r, us, "layout", "user")
}
