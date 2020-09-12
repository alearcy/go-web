package main

import (
	"errors"
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

type session struct {
	ID        int
	UUID      string
	email     string
	userID    int
	createdAt time.Time
}

func generateUUID() (string, error) {
	sID, err := uuid.NewV4()
	if err != nil {
		return "", errors.New("Error generating UUID")
	}
	sIDs := sID.String()
	return sIDs, nil
}

func generateSession(w http.ResponseWriter, user User, remember bool, uuid string) error {
	c := http.Cookie{
		Name:     "session",
		Value:    uuid,
		HttpOnly: true,
		// Secure: true solo HTTPS
	}
	if remember {
		// scade dopo un anno, altrimenti a ogni nuova sessione
		c.Expires = time.Now().Add(365 * 24 * time.Hour)
	}
	s := session{UUID: uuid, email: user.Email, userID: user.ID, createdAt: time.Now()}
	_, err := db.Exec("insert into sessions (uuid, email, user_id, created_at) values (?, ?, ?, ?)", &s.UUID, &s.email, &s.userID, &s.createdAt)
	if err != nil {
		flash(w, "Non è stato possibile creare la sessione utente, riprova.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	http.SetCookie(w, &c)
	return nil
}

func deleteCookie(w http.ResponseWriter, r *http.Request) (bool, error) {
	c, err := r.Cookie("session")
	if err != nil {
		return false, err
	}
	//TODO: cancellarla da DB
	c.MaxAge = -1
	c.Value = ""
	http.SetCookie(w, c)
	return true, nil
}

func isLoggedIn(r *http.Request) bool {
	_, err := r.Cookie("session")
	if err != nil {
		return false
	}
	//TODO: check also in DB session table
	return true
}
