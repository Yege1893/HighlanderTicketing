package model

import "google.golang.org/genproto/googleapis/type/date"

type Match struct {
	ID                    uint
	InitialTicketAmount   int32
	AvailableTicketAmount int32
	AwayMatch             bool
	Location              string
	Date                  date.Date
	Travel                Travel
	Orders                []Order
}

// Funktion ins Modell (siehe
//Myaktion), welche den available_ Ticket_Amount berechnet
