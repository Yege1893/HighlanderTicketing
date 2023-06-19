package service

import (
	"context"
	"fmt"

	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/db"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// noch testen
func AddMatchOrder(matchID primitive.ObjectID, order *model.Order) error {
	filter := bson.D{primitive.E{Key: "_id", Value: matchID}}
	order.ID = primitive.NewObjectID()

	updater := bson.D{primitive.E{Key: "$push", Value: bson.D{
		primitive.E{Key: "orders", Value: order},
	}}}

	client, err := db.GetMongoClient()
	if err != nil {
		return err
	}
	collection := client.Database(db.DB).Collection(db.MATCHES)

	updateResult, err := collection.UpdateOne(context.TODO(), filter, updater)
	if err != nil {
		return err
	}

	if updateResult.ModifiedCount == 0 {
		return fmt.Errorf("no document was updated")
	}

	return nil
}

func AddTravelOrder(matchID primitive.ObjectID, order *model.Order) error {
	filter := bson.D{primitive.E{Key: "_id", Value: matchID}}

	updater := bson.M{"$push": bson.M{"travel.orders": order}}

	client, err := db.GetMongoClient()
	if err != nil {
		return err
	}
	collection := client.Database(db.DB).Collection(db.MATCHES)

	updateResult, err := collection.UpdateOne(context.TODO(), filter, updater)
	if err != nil {
		return err
	}

	if updateResult.ModifiedCount == 0 {
		return fmt.Errorf("no document was updated")
	}

	return nil
}

func UpdateOrder() {

}
func DeleteOrder() {

}
