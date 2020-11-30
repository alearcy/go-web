package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
	"web/database"
	"web/forms"
	"web/utils"
)

type Menu struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	LogoUrl     string             `json:"logoUrl,omitempty" bson:"logoUrl,omitempty"`
	CompanyName string             `json:"companyName,omitempty" bson:"companyName,omitempty"`
	Address     string             `json:"address,omitempty" bson:"address,omitempty"`
	PIva        string             `json:"piva,omitempty" bson:"piva,omitempty"`
	UserId      primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
}

// CreateMenu - create a new user from the form
func CreateMenu(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("Name")
		logoUrl := r.FormValue("LogoUrl")
		companyName := r.FormValue("CompanyName")
		address := r.FormValue("Address")
		pIva := r.FormValue("PIva")

		mf := &forms.MenuForm{
			Name:        name,
			CompanyName: companyName,
			Address:     address,
			PIva:        pIva,
		}

		if ok := mf.Validate(); !ok {
			utils.GenerateHTML(w, r, mf, "layout", "createMenu")
			return
		}

		m := Menu{Name: name, LogoUrl: logoUrl, CompanyName: companyName, Address: address, PIva: pIva}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		collection := database.CallDb().Collection("menus")
		_, err := collection.InsertOne(ctx, m)
		if err != nil {
			utils.Flash(w, "Non è stato possibile salvare a database, riprova.")
			http.Redirect(w, r, "/menus/create", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/menus", http.StatusSeeOther)
		return

	}
	if r.Method == http.MethodGet {
		utils.GenerateHTML(w, r, nil, "layout", "createMenu")
		return
	}
}

//ListMenu: get all menus
func ListMenu(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var menus []Menu
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		collection := database.CallDb().Collection("menus")
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		err = cursor.All(ctx, &menus)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		utils.GenerateHTML(w, r, menus, "layout", "menus")
		return
	}
}
