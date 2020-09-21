package main

import (
	"github.com/satori/go.uuid"
	"net/http"
	"time"
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

func generateSession(w http.ResponseWriter, user User, remember string) (bool, error) {
	uuid, _ := generateUUID()
	c := http.Cookie{
		Name:     "session",
		Value:    uuid,
		HttpOnly: true,
		Path: "/",
		// Secure: true solo HTTPS
	}
	if remember == "remember" {
		// scade dopo un anno, altrimenti a ogni nuova sessione
		c.Expires = time.Now().Add(365 * 24 * time.Hour)
	}
	s := session{sessionID: uuid, userID: user.ID, createdAt: time.Now()}
	_, err := db.Exec("insert into sessions (uuid, user_id, created_at) values ($1, $2, $3)", &s.sessionID, &s.userID, &s.createdAt)
	if err != nil {
		flash(w, "Non è stato possibile creare la sessione utente, riprova.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false, err
	}
	http.SetCookie(w, &c)
	return true, nil
}

func deleteCookie(w http.ResponseWriter, r *http.Request) error {
	c, err := r.Cookie("session")
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM sessions WHERE uuid = $1", c.Value)
	if err != nil {
		flash(w, "Non è stato possibile sloggare l'utente.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	c.MaxAge = -1
	c.Value = ""
	http.SetCookie(w, c)
	return nil
}

func isLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		// non ho trovato il cookie quindi non sono loggato
		return false
	}
	// cerco nella tabella sessions se esiste quella con sessionID del cookie
	rows, err := db.Query("SELECT * FROM sessions where uuid = $1", c.Value)
	if err != nil {
		flash(w, "Non è stato possibile interrogare il DB, riprova.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	defer rows.Close()
	if rows.Next() {
		// ho trovato una sessione con il sessionID del cookie quindi sono loggato
		return true
	}
	return false
}
