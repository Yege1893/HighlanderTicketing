package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID             primitive.ObjectID `bson:"_id, omitempty"`
	Amount         int32              `bson:"amount"`
	User           User               `bson:"user, omitempty"`
	Ordernotified  bool               `bson:"ordernotified, omitempty"`
	Canceled       bool               `bson:"canceled, omitempty"`
	Cancelnotified bool               `bson:"cancelnotified, omitempty"`
}
