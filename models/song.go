package model

import (
	"net/http"
	"web/database"
	"web/utils"
)

// Song struct
type Song struct {
	ID       int
	Title    string
	UserId   string
	Image    string
	Filename string
}

// GetUsers get all active users from DB
func GetSongs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	rows, err := database.Db.Query("SELECT is, title, userId, image, filename FROM songs")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	songs := make([]Song, 0)
	defer rows.Close()
	for rows.Next() {
		song := Song{}
		err := rows.Scan(&song.ID, &song.Title, &song.UserId, &song.Image, &song.Filename)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		songs = append(songs, song)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	utils.ExecTemplate(w, r, songs, "songs.gohtml")
}

// GetUser get a single user from thd DB passing the known ID
// func GetUser(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
// 		return
// 	}
// 	params := mux.Vars(r)
// 	param := params["id"]

// 	row := database.Db.QueryRow("SELECT name, surname, email, role FROM users where id = $1", param)

// 	us := User{}
// 	err := row.Scan(&us.Name, &us.Surname, &us.Email, &us.Role)
// 	switch {
// 	case err == sql.ErrNoRows:
// 		http.NotFound(w, r)
// 		return
// 	case err != nil:
// 		fmt.Println(err)
// 		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
// 		return
// 	}
// 	utils.GenerateHTML(w, r, us, "layout", "user")
// }
