package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
	"web/database"
	"web/forms"
	model "web/models"
	"web/utils"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	"log/slog"
)

type session struct {
	ID        int
	sessionID string
	userID    int
	createdAt time.Time
}

func generateUUID() string {
	sID := uuid.NewV4()
	sIDs := sID.String()
	return sIDs
}

// TODO: aggiungere da file env se è dev o prod e mettere o meno https
// GenerateSession - genmerate a new session passing the ResponseWriter, the user and the remember option
func GenerateSession(w http.ResponseWriter, user model.User, remember string) (bool, error) {
	uuID := generateUUID()
	c := http.Cookie{
		Name:     "session",
		Value:    uuID,
		HttpOnly: true,
		Path:     "/",
		// solo HTTPS
		// Secure: true
	}
	if remember == "remember" {
		// scade dopo un anno, altrimenti a ogni nuova sessione
		c.Expires = time.Now().Add(365 * 24 * time.Hour)
	}
	s := session{sessionID: uuID, userID: user.ID, createdAt: time.Now()}
	_, err := database.Db.Exec("insert into sessions (uuId, user_id, created_at) values ($1, $2, $3)", &s.sessionID, &s.userID, &s.createdAt)
	// TODO: rendere generica la gestione degli errori con relativo flash message e fmt.Errorf
	if err != nil {
		utils.Flash(w, "Non è stato possibile creare la sessione utente, riprova.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false, err
	}
	http.SetCookie(w, &c)
	return true, nil
}

// DeleteCookie - delete the session cookie passing ResponseWrtier and Request pointer
func DeleteCookie(w http.ResponseWriter, r *http.Request) error {
	c, err := r.Cookie("session")
	if err != nil {
		return err
	}
	_, err = database.Db.Exec("DELETE FROM sessions WHERE uuid = $1", c.Value)
	if err != nil {
		fmt.Println(err)
		utils.Flash(w, "Non è stato possibile fare il logout dell'utente.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	c.MaxAge = -1
	c.Value = ""
	http.SetCookie(w, c)
	return nil
}

// IsLoggedIn - check if a session cookie exists and if the user and the session wxists in DB
func IsLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		// non ho trovato il cookie quindi non sono loggato
		return false
	}
	// cerco nella tabella sessions se esiste quella con sessionID del cookie
	rows, err := database.Db.Query("SELECT * FROM sessions where uuid = $1", c.Value)
	if err != nil {
		fmt.Println(err)
		utils.Flash(w, "Non è stato possibile interrogare il DB, riprova.") // TODO: il flash si vede o viene cancellato dall'http error??
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	defer rows.Close()
	// ho trovato una sessione con il sessionID del cookie quindi sono loggato
	return rows.Next()
}

// Signup - create a new user from the form
func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		name := r.FormValue("name")
		surname := r.FormValue("surname")
		password := r.FormValue("password")
		password2 := r.FormValue("password2")

		sf := &forms.SignupForm{
			Email:     email,
			Password:  password,
			Password2: password2,
			Name:      name,
			Surname:   surname,
		}

		if ok := sf.Validate(); !ok {
			url, err := utils.FormatURL(r)
			if err != nil {
				slog.Error(err.Error())
				return
			}
			utils.GenerateTemplate(w, r, nil, url)
			return
		}

		//controllare se l'utente esiste prima di continuare
		row := database.Db.QueryRow("SELECT * FROM users WHERE email = $1", email)
		u := model.User{}
		err := row.Scan(&u.ID)
		if err != nil {
			cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			u := model.User{Name: name, Surname: surname, Email: email, Password: cryptedPassword, Role: 1}
			_, err = database.Db.Exec("insert into users (name, surname, email, password, role) values ($1, $2, $3, $4, $5)", &u.Name, &u.Surname, &u.Email, &u.Password, &u.Role)
			if err != nil {
				fmt.Println(err)
				utils.Flash(w, "Non è stato possibile salvare a DB, riprova.")
				http.Redirect(w, r, "/signup", http.StatusFound)
				return
			}
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		utils.Flash(w, "Utente già esistente")
		http.Redirect(w, r, "/signup", http.StatusFound)
	}
	if !IsLoggedIn(w, r) {
		url, err := utils.FormatURL(r)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		utils.GenerateTemplate(w, r, nil, url)
		return
	}
	utils.Flash(w, "Risulti già loggato!")
	http.Redirect(w, r, "/", http.StatusFound)
}

// Login - login and create session from the user login page
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		remember := r.FormValue("remember")

		lf := &forms.LoginForm{
			Email:    email,
			Password: password,
			Remember: remember,
		}

		if ok := lf.Validate(); !ok {
			url, err := utils.FormatURL(r)
			if err != nil {
				slog.Error(err.Error())
				return
			}
			utils.GenerateTemplate(w, r, nil, url)
			return
		}

		//controllo se l'utente esiste prima di continuare
		row := database.Db.QueryRow("SELECT id, name, surname, email, password, role FROM users WHERE email = $1", email)
		u := model.User{}
		err := row.Scan(&u.ID, &u.Name, &u.Surname, &u.Email, &u.Password, &u.Role)
		switch {
		case err == sql.ErrNoRows:
			utils.Flash(w, "Nome utente errato o non esistente.")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		case err != nil:
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		err = bcrypt.CompareHashAndPassword(u.Password, []byte(password))
		if err != nil {
			utils.Flash(w, "Password non valida")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if ok, _ := GenerateSession(w, u, remember); ok {
			http.Redirect(w, r, "/admin", http.StatusFound)
		}
	}
	if r.Method == http.MethodGet {
		if IsLoggedIn(w, r) {
			http.Redirect(w, r, "/admin", http.StatusFound)
			return
		}
		url, err := utils.FormatURL(r)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		utils.GenerateTemplate(w, r, nil, url)
		return
	}
}

// Logout user and delete session cookie
func Logout(w http.ResponseWriter, r *http.Request) {
	if !IsLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	err := DeleteCookie(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

// Protected - a middleware to protect routes from unauthenticated users
func Protected(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok := IsLoggedIn(w, r)
		if !ok {
			utils.Flash(w, "Devi essere loggato.")
			http.Redirect(w, r, "/login", http.StatusFound)
		}
		next(w, r)
	}
}
