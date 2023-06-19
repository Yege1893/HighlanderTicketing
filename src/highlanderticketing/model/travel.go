package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/genproto/googleapis/type/date"
)

type Travel struct {
	ID                  primitive.ObjectID `bson:"_id, omitempty"`
	TravelType          TravelType
	InitialSeatAmount   int32
	AvailableSeatAmount int32
	StartLocation       string
	EndLocation         string
	StartDate           date.Date
	Orders              []Order
}

type TravelType string

const (
	CAR   TravelType = "CAR"
	BUS   TravelType = "BUS"
	PLANE TravelType = "PLANE"
)

// Funktion ins Modell (siehe
//Myaktion), welche den available_ Ticket_Amount berechnet
