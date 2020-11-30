package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	SubCategories []SubCategory      `json:"subCategories,omitempty" bson:"subCategories,omitempty"`
}
