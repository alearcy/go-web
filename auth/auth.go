package auth

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
	"web/database"
	"web/forms"
	"web/models"
	"web/utils"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type session struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SessionID string             `json:"sessionID,omitempty" bson:"sessionID,omitempty"`
	UserID    primitive.ObjectID `json:"userID,omitempty" bson:"userID,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}

func generateUUID() (string, error) {
	sID := uuid.NewV4()
	sIDs := sID.String()
	return sIDs, nil
}

// GenerateSession - generate a new session passing the ResponseWriter, the user and the remember option
func GenerateSession(w http.ResponseWriter, user models.User, remember string) (bool, error) {
	uuID, _ := generateUUID()
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
	s := session{SessionID: uuID, UserID: user.ID, CreatedAt: time.Now()}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := database.CallDb().Collection("sessions")
	_, err := collection.InsertOne(ctx, s)

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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := database.CallDb().Collection("sessions")
	_, err = collection.DeleteOne(ctx, bson.M{"sessionID": c.Value})
	if err != nil {
		utils.Flash(w, "Non è stato possibile eseguire il logout dell'utente.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	c.MaxAge = -1
	c.Value = ""
	http.SetCookie(w, c)
	return nil
}

// IsLoggedIn - check if a session cookie exists and if the user and the session wxists in CallDb()
func IsLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		// non ho trovato il cookie quindi non sono loggato
		return false
	}
	// cerco nella tabella sessions se esiste quella con sessionID del cookie
	var s session
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := database.CallDb().Collection("sessions")
	err = collection.FindOne(ctx, bson.M{"sessionID": c.Value}).Decode(&s)
	if err != nil {
		utils.Flash(w, "Non è stato possibile interrogare il CallDb(), riprova.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	return true
}

// Signup - create a new user from the form
func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		name := r.FormValue("name")
		surname := r.FormValue("surname")
		password := r.FormValue("password")
		password2 := r.FormValue("password2")

		if password != password2 {
			utils.Flash(w, "Le due password non coincidono")
			http.Redirect(w, r, "/users/create", http.StatusSeeOther)
			return
		}

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
		var u models.User

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		collection := database.CallDb().Collection("users")
		err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&u)

		if err != nil {
			cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			u := models.User{Name: name, Surname: surname, Email: email, Password: cryptedPassword, Role: 1}
			_, err = collection.InsertOne(ctx, u)
			if err != nil {
				utils.Flash(w, "Non è stato possibile salvare a database, riprova.")
				http.Redirect(w, r, "/users/create", http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} else {
			utils.Flash(w, "Utente già esistente.")
			http.Redirect(w, r, "/users/create", http.StatusSeeOther)
			return
		}

	}
	if !IsLoggedIn(w, r) {
		utils.GenerateHTML(w, r, nil, "layout", "signup")
		return
	}
	utils.Flash(w, "Risulti già loggato!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
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
			utils.GenerateHTML(w, r, lf, "layout", "login")
			return
		}

		//controllo se l'utente esiste prima di continuare
		var u models.User
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		collection := database.CallDb().Collection("users")
		err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&u)
		if err != nil {
			utils.Flash(w, "Nome utente errato o non esistente.")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
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

// Logout user and delete session cookie
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

// Protected - a middleware to protect routes
func Protected(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok := IsLoggedIn(w, r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
