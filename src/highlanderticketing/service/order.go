package service

import (
	"context"
	"fmt"

	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/db"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddMatchOrder(matchID primitive.ObjectID, order *model.Order) error {
	filter := bson.D{primitive.E{Key: "_id", Value: matchID}}
	order.ID = primitive.NewObjectID()
	matchToFind := model.Match{}

	updater := bson.D{primitive.E{Key: "$push", Value: bson.D{
		primitive.E{Key: "orders", Value: order},
	}}}

	client, err := db.GetMongoClient()
	if err != nil {
		return err
	}

	collection := client.Database(db.DB).Collection(db.MATCHES)

	err = collection.FindOne(context.TODO(), filter).Decode(&matchToFind)
	if err != nil {
		return err
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, updater)
	if err != nil {
		return err
	}

	if updateResult.ModifiedCount == 0 {
		return fmt.Errorf("no document was updated, please send order again")
	}

	natsServer, err := ConnectToNats()
	defer natsServer.Nc.Close()

	emailContent := model.EmialContent{Name: order.User.Name, AwayMatch: matchToFind.AwayMatch, Location: matchToFind.Location, Date: matchToFind.Date, Emailadress: order.User.Email, OrderID: matchToFind.ID.String()}
	if err := natsServer.ConfirmOrder(&emailContent); err != nil {
		// hier warten und nochmal versuchen zu senden
		order.Ordernotified = false
		return fmt.Errorf("error sending confirm email: %v", err)
	} else {
		order.Ordernotified = true
	}

	updaterNotification := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "orders.$[element]", Value: order},
	}}}

	options := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{Key: "element._id", Value: order.ID}},
		},
	})

	updateNotification, err := collection.UpdateOne(context.TODO(), filter, updaterNotification, options)
	if err != nil {
		return fmt.Errorf("no document was updated, please send order again")
	}

	if updateNotification.ModifiedCount == 0 {
		return fmt.Errorf("no document was updated, please send order again")
	}

	return nil
}

func AddTravelOrder(matchID primitive.ObjectID, order *model.Order) error {
	filter := bson.D{primitive.E{Key: "_id", Value: matchID}}
	order.ID = primitive.NewObjectID()
	matchToFind := model.Match{}

	updater := bson.M{"$push": bson.M{"travel.orders": order}}

	client, err := db.GetMongoClient()
	if err != nil {
		return err
	}

	collection := client.Database(db.DB).Collection(db.MATCHES)

	err = collection.FindOne(context.TODO(), filter).Decode(&matchToFind)
	if err != nil {
		return err
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, updater)
	if err != nil {
		return err
	}

	if updateResult.ModifiedCount == 0 {
		return fmt.Errorf("no document was updated, please send order again")
	}

	natsServer, err := ConnectToNats()
	defer natsServer.Nc.Close()

	emailContent := model.EmialContent{Name: order.User.Name, AwayMatch: matchToFind.AwayMatch, Location: matchToFind.Location, Date: matchToFind.Date, Emailadress: order.User.Email, OrderID: matchToFind.ID.String()}
	if err := natsServer.ConfirmOrder(&emailContent); err != nil {
		// hier warten und nochmal versuchen zu senden
		order.Ordernotified = false
		return fmt.Errorf("error sending confirm email %v", err)
	} else {
		order.Ordernotified = true
	}

	updaterNotification := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "orders.$[element]", Value: order},
	}}}

	options := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{Key: "element._id", Value: order.ID}},
		},
	})

	updateNotification, err := collection.UpdateOne(context.TODO(), filter, updaterNotification, options)
	if err != nil {
		return fmt.Errorf("no document was updated, please send order again")
	}

	if updateNotification.ModifiedCount == 0 {
		return fmt.Errorf("no document was updated, please send order again")
	}

	return nil
}

var isMatchOrder bool = true

func CancelOrder(matchID primitive.ObjectID, order *model.Order) error {
	filter := bson.D{primitive.E{Key: "_id", Value: matchID}}
	matchToFind, err := GetMatchByID(matchID)
	if err != nil {
		return err
	}

	client, err := db.GetMongoClient()
	if err != nil {
		return err
	}

	collection := client.Database(db.DB).Collection(db.MATCHES)

	order.Canceled = true

	updaterMatchCancel := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "orders.$[element]", Value: order},
	}}}

	updaterTravelCancel := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "travel.orders.$[element]", Value: order},
	}}}

	options := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{Key: "element._id", Value: order.ID}},
		},
	})

	updateMatchCancel, err := collection.UpdateOne(context.TODO(), filter, updaterMatchCancel, options)
	if err != nil {
		return err
	}
	if updateMatchCancel.ModifiedCount == 0 {
		isMatchOrder = false
		updateTravelCancel, err := collection.UpdateOne(context.TODO(), filter, updaterTravelCancel, options)
		if err != nil {
			return err
		}
		if updateTravelCancel.ModifiedCount == 0 {
			return fmt.Errorf("orderid not in system")
		}

	}

	natsServer, err := ConnectToNats()
	defer natsServer.Nc.Close()

	emailContent := model.EmialContent{Name: order.User.Name, AwayMatch: matchToFind.AwayMatch, Location: matchToFind.Location, Date: matchToFind.Date, Emailadress: order.User.Email, OrderID: order.ID.String()}
	if err := natsServer.ConfirmCancel(&emailContent); err != nil {
		return fmt.Errorf("error sending confirm email %v", err)
	}

	order.Cancelnotified = true

	if isMatchOrder == false {
		updateTravelCancelNotifi, err := collection.UpdateOne(context.TODO(), filter, updaterTravelCancel, options)
		if err != nil {
			return err
		}
		if updateTravelCancelNotifi.ModifiedCount == 0 {
			return fmt.Errorf("orderid not in system")
		}
		isMatchOrder = true
	} else {
		updateMatchCancelNotifi, err := collection.UpdateOne(context.TODO(), filter, updaterMatchCancel, options)
		if err != nil {
			return err
		}
		if updateMatchCancelNotifi.ModifiedCount == 0 {
			return fmt.Errorf("orderid not in system")
		}
	}
	return nil

}
func deleteOrder(orderID primitive.ObjectID) error {
	filter := bson.D{primitive.E{Key: "_id", Value: orderID}}
	client, err := db.GetMongoClient()
	if err != nil {
		return err
	}
	collection := client.Database(db.DB).Collection(db.MATCHES)
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}
