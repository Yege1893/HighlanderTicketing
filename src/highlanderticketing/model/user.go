package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID         primitive.ObjectID `bson:"_id, omitempty"`
	Email      string             `json:"email" bson:"email"`
	GoogleID   string             `json:"id" bson:"google_id"`
	Name       string             `json:"name" bson:"name"`
	FamilyName string             `json:"family_name" bson:"family_name"`
	IsAdmin    bool               `json:"is_admin" bson:"is_admin"`
}
