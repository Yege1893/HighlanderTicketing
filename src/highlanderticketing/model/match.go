package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/genproto/googleapis/type/date"
)

type Match struct {
	ID                    primitive.ObjectID `bson:"_id, omitempty"`
	InitialTicketAmount   int32              `bson:"initial_ticket_amount"`
	AvailableTicketAmount int32              `bson:"available_ticket_amount"`
	AwayMatch             bool               `bson:"away_match"`
	Location              string             `bson:"location"`
	Date                  date.Date          `bson:"date"`
	Travel                Travel             `bson:"travel, omitempty"`
	Orders                []Order            `bson:"orders"`
}

/*func calculateAmountDonated(matchID primitive.ObjectID) (int32, error) {
	collection := client.Database("your_db").Collection("orders")
	pipeline := []bson.M{
		{"$match": bson.M{"match_id": matchID}},
		{"$group": bson.M{"_id": nil, "totalAmount": bson.M{"$sum": "$amount"}}},
	}

	var result struct {
		TotalAmount int32 `bson:"totalAmount"`
	}

	err := collection.Aggregate(context.TODO(), pipeline).Decode(&result)
	if err != nil {
		return 0, err
	}

	return result.TotalAmount, nil
}*/ // muss nach jedem update , erstellen gemacht werden , funktion anpassen
