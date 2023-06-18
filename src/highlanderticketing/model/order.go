package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID        primitive.ObjectID `bson:"_id, omitempty"`
	OrderType OrderType          `bson:"order_type"`
	Amount    int32              `bson:"amount"`
	User      User               `bson:"user"`
}
type OrderType string

const (
	MATCHTICKET OrderType = "MATCHTICKET"
	BUSTICKET   OrderType = "TRAVELTICKET"
)
