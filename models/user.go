package model

import (
	"database/sql"
	"fmt"
	"net/http"
	"web/database"
	"web/utils"
)

// User struct
type User struct {
	ID       int
	Name     string
	Surname  string
	Email    string
	Password []byte
	Role     int
}

// GetUsers get all active users from DB
func GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	param := r.URL.Query().Get("id")
	if param != "" {
		row := database.Db.QueryRow("SELECT name, surname, email, role FROM users where id = $1", param)
		us := User{}
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
	} else {

		rows, err := database.Db.Query("SELECT id, name, surname, email, password, role FROM users")
		if err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		urs := make([]User, 0)
		defer rows.Close()
		for rows.Next() {
			us := User{}
			err := rows.Scan(&us.ID, &us.Name, &us.Surname, &us.Email, &us.Password, &us.Role)
			fmt.Println(err)
			if err != nil {
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
			urs = append(urs, us)
		}
		if err = rows.Err(); err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		utils.GenerateHTML(w, r, urs, "layout", "users")
	}

}
