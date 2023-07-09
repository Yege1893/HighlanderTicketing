package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Match struct {
	ID                    primitive.ObjectID `bson:"_id, omitempty"`
	InitialTicketAmount   int32              `bson:"initial_ticket_amount"`
	AvailableTicketAmount int32              `bson:"available_ticket_amount"`
	AwayMatch             bool               `bson:"away_match"`
	Location              string             `bson:"location"`
	Date                  string             `bson:"date, omitempty"`
	OrderAmount           int32              `bson:"orderamount, omitempty"`
	Orders                []Order            `bson:"orders"`
}

// muss nach jedem update , erstellen gemacht werden , funktion anpassen
