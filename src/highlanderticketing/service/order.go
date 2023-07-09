package service

import (
	"context"
	"errors"
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
	matchToFind := &model.Match{}
	emailContent := model.EmialContent{Name: order.User.Name, AwayMatch: matchToFind.AwayMatch, Location: matchToFind.Location, Date: matchToFind.Date, Emailadress: order.User.Email, OrderID: matchToFind.ID.String()}

	updater := bson.D{primitive.E{Key: "$push", Value: bson.D{
		primitive.E{Key: "orders", Value: order},
	}}}

	updaterNotification := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "orders.$[element]", Value: order},
	}}}

	options := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{Key: "element._id", Value: order.ID}},
		},
	})

	client, err := db.GetMongoClient()
	if err != nil {
		return err
	}

	collection := client.Database(db.DB).Collection(db.MATCHES)

	err = collection.FindOne(context.TODO(), filter).Decode(&matchToFind)
	if err != nil {
		return err
	}
	if matchToFind.AvailableTicketAmount < order.Amount {
		return fmt.Errorf("ticket amount not available")
	} else {
		matchToFind.AvailableTicketAmount = matchToFind.AvailableTicketAmount - order.Amount
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, updater)
	if err != nil {
		return err
	}
	if updateResult.ModifiedCount == 0 {
		return fmt.Errorf("no document was updated, please send order again")
	}

	_, errUpdate := UpdateMatch(matchToFind.ID, matchToFind)
	if errUpdate != nil {
		errUpdate = fmt.Errorf("can not update match amount, please send order again")
		err := deleteOrder(order.ID, matchToFind.ID)
		natsServer, err := ConnectToNats()
		if err != nil {
			return err
		}
		defer natsServer.Nc.Close()
		if err := natsServer.ConfirmCancel(&emailContent); err != nil {
			// hier warten und nochmal versuchen zu senden
			order.Cancelnotified = false
			err := fmt.Errorf("error sending confirm email: %v", err)
			return err
		}
	}

	natsServer, err := ConnectToNats()
	if err != nil {
		return err
	}

	defer natsServer.Nc.Close()

	if err := natsServer.ConfirmOrder(&emailContent); err != nil {
		// hier warten und nochmal versuchen zu senden
		order.Ordernotified = false
		err = fmt.Errorf("error sending confirm email: %v", err)
		return err
	} else {
		order.Ordernotified = true
	}

	updateNotification, err := collection.UpdateOne(context.TODO(), filter, updaterNotification, options)
	if err != nil {
		err = fmt.Errorf("no document was updated, please send order again")
		return err
	}

	if updateNotification.ModifiedCount == 0 {
		err = fmt.Errorf("no document was updated, please send order again")
		return err
	}

	return nil
}

func CancelOrder(matchID primitive.ObjectID, order *model.Order) error {
	fmt.Println("order", order)
	if order.Canceled == true {
		return fmt.Errorf("order already canceled")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: matchID}}
	matchToFind, err := GetMatchByID(matchID)
	if err != nil {
		return err
	} else {
		matchToFind.AvailableTicketAmount = matchToFind.AvailableTicketAmount + order.Amount
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
		return fmt.Errorf("not updated")

	}

	_, errUpdateAmount := UpdateMatch(matchToFind.ID, matchToFind)
	if errUpdateAmount != nil {
		order.Canceled = false
		updateMatchCancel, err := collection.UpdateOne(context.TODO(), filter, updaterMatchCancel, options)
		if err != nil {
			return err
		}
		if updateMatchCancel.ModifiedCount == 0 {
			return fmt.Errorf("not updated")

		}
		return fmt.Errorf("error canceling match internal, please try again %v", err)
	}

	natsServer, err := ConnectToNats()
	defer natsServer.Nc.Close()

	emailContent := model.EmialContent{Name: order.User.Name, AwayMatch: matchToFind.AwayMatch, Location: matchToFind.Location, Date: matchToFind.Date, Emailadress: order.User.Email, OrderID: order.ID.String()}
	if err := natsServer.ConfirmCancel(&emailContent); err != nil {
		return fmt.Errorf("error sending confirm email %v", err)
	} else {
		order.Cancelnotified = true
	}

	updateMatchCancelNotifi, err := collection.UpdateOne(context.TODO(), filter, updaterMatchCancel, options)
	if err != nil {
		return err
	}
	if updateMatchCancelNotifi.ModifiedCount == 0 {
		return fmt.Errorf("orderid not in system")
	}
	return nil

}

func GetOrderById(orderID primitive.ObjectID) (*model.Order, error) {
	client, err := db.GetMongoClient()
	if err != nil {
		return nil, err
	}
	collection := client.Database(db.DB).Collection(db.MATCHES)

	filter := bson.M{"orders._id": orderID}

	var result model.Match

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	for _, order := range result.Orders {
		if order.ID == orderID {
			return &order, nil
		}
	}

	return nil, errors.New("Order not found")
}

func deleteOrder(matchID primitive.ObjectID, orderID primitive.ObjectID) error {
	filter := bson.D{primitive.E{Key: "_id", Value: matchID}}
	updater := bson.D{primitive.E{Key: "$pull", Value: bson.D{
		primitive.E{Key: "orders", Value: bson.D{
			primitive.E{Key: "_id", Value: orderID},
		}},
	}}}

	client, err := db.GetMongoClient()
	if err != nil {
		return err
	}

	collection := client.Database(db.DB).Collection(db.MATCHES)

	_, err = collection.UpdateOne(context.TODO(), filter, updater)
	if err != nil {
		return err
	}

	return nil
}
