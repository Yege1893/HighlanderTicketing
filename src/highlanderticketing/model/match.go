package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Match struct {
	ID                    primitive.ObjectID `bson:"_id, omitempty"`
	ExternalID            int64              `bson:"externalID"`
	InitialTicketAmount   int32              `bson:"initial_ticket_amount"`
	AvailableTicketAmount int32              `bson:"available_ticket_amount"`
	Opponenent            string             `bson:"opponent"`
	LeagueName            string             `bson:"league_name"`
	AwayMatch             bool               `bson:"away_match"`
	Location              string             `bson:"location"`
	Date                  time.Time          `bson:"date"`
	OrderAmount           int32              `bson:"orderamount"`
	Orders                []Order            `bson:"orders"`
}

// muss nach jedem update , erstellen gemacht werden , funktion anpassen
