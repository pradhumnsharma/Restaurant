package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MenuItem struct {
	ID          primitive.ObjectID `bson:"_id"`
	FoodName    *string            `json:"foodName"`
	Price       *float64           `json:"price"`
	Description *string            `json:"description"`
}
