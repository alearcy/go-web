package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name            string             `json:"name,omitempty" bson:"name,omitempty"`
	Descritpion     string             `json:"description,omitempty" bson:"description,omitempty"`
	Ingredients     string             `json:"ingredients,omitempty" bson:"ingredients,omitempty"`
	FullPrice       string             `json:"fullPrice,omitempty" bson:"fullPrice,omitempty"`
	DiscountedPrice string             `json:"discountedPrice,omitempty" bson:"discountedPrice,omitempty"`
	Tags            []Tag              `json:"tags,omitempty" bson:"tags,omitempty"`
	Allergens       []Allergen         `json:"allergens,omitempty" bson:"allergens,omitempty"`
	CategoryID      primitive.ObjectID `json:"categoryId,omitempty" bson:"categoryId,omitempty"`
	SubCategoryID   primitive.ObjectID `json:"subCategoryId,omitempty" bson:"subCategoryId,omitempty"`
	MenuID          primitive.ObjectID `json:"menuId,omitempty" bson:"menuId,omitempty"`
}
