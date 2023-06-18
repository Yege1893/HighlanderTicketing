package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Match struct {
	ID                    primitive.ObjectID `bson:"_id"`
	InitialTicketAmount   int32              `bson:"initial_ticket_amount"`
	AvailableTicketAmount int32              `bson:"available_ticket_amount"`
	AwayMatch             bool               `bson:"away_match"`
	Location              string             `bson:"location"`
	//Date                  date.Date
	//Travel                Travel
	//Orders                []Order
}

// Funktion ins Modell (siehe
//Myaktion), welche den available_ Ticket_Amount berechnet
