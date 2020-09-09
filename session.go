package main

import (
	"errors"
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

func generateUUID() (string, error) {
	sID, err := uuid.NewV4()
	if err != nil {
		return "", errors.New("Error generating UUID")
	}
	sIDs := sID.String()
	return sIDs, nil
}

func generateCookie(w http.ResponseWriter, uuid string) *http.Cookie {
	c := &http.Cookie{
		Name:     "session",
		Value:    uuid,
		HttpOnly: true,
		// Secure: true solo HTTPS
	}
	c.MaxAge = sessionLength
	http.SetCookie(w, c)
	return c
}

func deleteCookie(w http.ResponseWriter, r *http.Request) (bool, error) {
	c, err := r.Cookie("session")
	if err != nil {
		return false, err
	}
	c.MaxAge = -1
	c.Value = ""
	http.SetCookie(w, c)
	return true, nil
}

func cleanSessions() {
	for k, v := range dbSessions {
		if time.Now().Sub(v.activityTime) > (time.Second * 30) {
			delete(dbSessions, k)
		}
	}
	dbSessionsCleaned = time.Now()
}

func isLoggedIn(r *http.Request) bool {
	_, err := r.Cookie("session")
	//TODO: check also in DB session table
	if err != nil {
		return false
	}
	return true
}
