package session

import (
	"github.com/satori/go.uuid"
	"net/http"
	"time"
	"web/auth"
	"web/database"
	"web/utils"
)

type session struct {
	ID        int
	sessionID string
	userID    int
	createdAt time.Time
}

func generateUUID() (string, error) {
	sID := uuid.NewV4()
	sIDs := sID.String()
	return sIDs, nil
}

func GenerateSession(w http.ResponseWriter, user auth.User, remember string) (bool, error) {
	uuId, _ := generateUUID()
	c := http.Cookie{
		Name:     "session",
		Value:    uuId,
		HttpOnly: true,
		Path:     "/",
		// solo HTTPS
		// Secure: true
	}
	if remember == "remember" {
		// scade dopo un anno, altrimenti a ogni nuova sessione
		c.Expires = time.Now().Add(365 * 24 * time.Hour)
	}
	s := session{sessionID: uuId, userID: user.ID, createdAt: time.Now()}
	_, err := database.Db.Exec("insert into sessions (uuId, user_id, created_at) values ($1, $2, $3)", &s.sessionID, &s.userID, &s.createdAt)
	if err != nil {
		utils.Flash(w, "Non è stato possibile creare la sessione utente, riprova.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false, err
	}
	http.SetCookie(w, &c)
	return true, nil
}

func DeleteCookie(w http.ResponseWriter, r *http.Request) error {
	c, err := r.Cookie("session")
	if err != nil {
		return err
	}
	_, err = database.Db.Exec("DELETE FROM sessions WHERE uuid = $1", c.Value)
	if err != nil {
		utils.Flash(w, "Non è stato possibile sloggare l'utente.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	c.MaxAge = -1
	c.Value = ""
	http.SetCookie(w, c)
	return nil
}

func IsLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		// non ho trovato il cookie quindi non sono loggato
		return false
	}
	// cerco nella tabella sessions se esiste quella con sessionID del cookie
	rows, err := database.Db.Query("SELECT * FROM sessions where uuid = $1", c.Value)
	if err != nil {
		utils.Flash(w, "Non è stato possibile interrogare il DB, riprova.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	if rows.Next() {
		// ho trovato una sessione con il sessionID del cookie quindi sono loggato
		return true
	}
	return false
}
