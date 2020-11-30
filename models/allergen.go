package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Allergen struct {
	ID   primitive.ObjectID
	Name string
}
