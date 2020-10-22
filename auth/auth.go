package auth

import (
	"database/sql"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
	"web/database"
	"web/forms"
	"web/utils"
)

type User struct {
	ID       int
	Name     string
	Surname  string
	Email    string
	Password []byte
	Role     int
}

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

func GenerateSession(w http.ResponseWriter, user User, remember string) (bool, error) {
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
			utils.GenerateHTML(w, r, sf, "layout", "signup")
			return
		}

		//controllare se l'utente esiste prima di continuare
		row := database.Db.QueryRow("SELECT * FROM users WHERE email = $1", email)
		u := User{}
		err := row.Scan(&u.ID)
		if err != nil {
			if password != password2 {
				utils.Flash(w, "Le due password non coincidono")
				http.Redirect(w, r, "/users/create", http.StatusSeeOther)
				return
			}
			cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			u := User{Name: name, Surname: surname, Email: email, Password: cryptedPassword, Role: 1}
			_, err = database.Db.Exec("insert into users (name, surname, email, password, role) values ($1, $2, $3, $4, $5)", &u.Name, &u.Surname, &u.Email, &u.Password, &u.Role)
			if err != nil {
				utils.Flash(w, "Non è stato possibile salvare a DB, riprova.")
				http.Redirect(w, r, "/users/create", http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		utils.Flash(w, "Utente già esistente")
		http.Redirect(w, r, "/users/create", http.StatusSeeOther)
	}
	if !IsLoggedIn(w, r) {
		utils.GenerateHTML(w, r, nil, "layout", "signup")
		return
	}
	utils.Flash(w, "Risulti già loggato!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

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
			utils.GenerateHTML(w, r, lf, "layout", "login")
			return
		}

		//controllo se l'utente esiste prima di continuare
		row := database.Db.QueryRow("SELECT id, name, surname, email, password, role FROM users WHERE email = $1", email)
		u := User{}
		err := row.Scan(&u.ID, &u.Name, &u.Surname, &u.Email, &u.Password, &u.Role)
		switch {
		case err == sql.ErrNoRows:
			utils.Flash(w, "Nome utente errato o non esistente.")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		case err != nil:
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		err = bcrypt.CompareHashAndPassword(u.Password, []byte(password))
		if err != nil {
			utils.Flash(w, "Password non valida")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if ok, _ := GenerateSession(w, u, remember); ok {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		}
	}
	if r.Method == http.MethodGet {
		if IsLoggedIn(w, r) {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
		utils.GenerateHTML(w, r, nil, "layout", "login")
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if !IsLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err := DeleteCookie(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Protected(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok := IsLoggedIn(w, r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}
		h(w, r)
	}
}
