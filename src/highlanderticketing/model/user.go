package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	GoogleID   string             `json:"id,omitempty" bson:"google_id"`
	Email      string             `json:"email" bson:"email"`
	Name       string             `json:"name" bson:"name"`
	FamilyName string             `json:"family_name" bson:"family_name"`
	IsAdmin    bool               `json:"is_admin" bson:"is_admin"`
}
