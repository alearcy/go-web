package models

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
	"web/database"
	"web/utils"
)

// User struct
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Surname  string             `json:"surname,omitempty" bson:"surname,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password []byte             `json:"password,omitempty" bson:"password,omitempty"`
	Role     int                `json:"role,omitempty" bson:"role,omitempty"`
}

// GetUsers get all active users from DB
func GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	var users []User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := database.CallDb().Collection("users")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	err = cursor.All(ctx, &users)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	utils.GenerateHTML(w, r, users, "layout", "users")
}

// GetUser get a single user from thd DB passing the known ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	params := mux.Vars(r)
	param := params["id"]

	var u User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := database.CallDb().Collection("users")
	err := collection.FindOne(ctx, bson.M{"_id": param}).Decode(&u)

	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	utils.GenerateHTML(w, r, u, "layout", "user")
}
