package service

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/db"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateMatch(match *model.Match) error {
	match.ID = primitive.NewObjectID()
	match.Orders = []model.Order{}
	client, err := db.GetMongoClient()
	if err != nil {
		log.Errorf("Eror gettin db client: %v", err)
		return err
	}
	collection := client.Database(db.DB).Collection(db.MATCHES)

	_, err = collection.InsertOne(context.TODO(), match)
	if err != nil {
		log.Errorf("Erorr inserting match: %v", err)
		return nil
	}
	log.Info("match inserted: %v", match)
	return nil
}

func InserExternalMatch(match *model.Match) error {
	match.ID = primitive.NewObjectID()
	match.Orders = []model.Order{}
	existingMatch := &model.Match{}
	filter := bson.M{"externalID": match.ExternalID}

	client, err := db.GetMongoClient()
	if err != nil {
		log.Errorf("Eror gettin db client: %v", err)
		return err
	}

	collection := client.Database(db.DB).Collection(db.MATCHES)
	err = collection.FindOne(context.TODO(), filter).Decode(existingMatch)
	if err == nil {
		log.Info("match already exists")
		return fmt.Errorf("Match mit ExternalID %d exists already", match.ExternalID)
	} else if err != mongo.ErrNoDocuments {
		return err
	}
	_, err = collection.InsertOne(context.TODO(), match)
	if err != nil {
		return nil
	}
	log.Info("match inserted: %v", match)
	return nil
}

func UpdateMatch(matchID primitive.ObjectID, match *model.Match) (*model.Match, error) {
	result := model.Match{}

	filter := bson.D{primitive.E{Key: "_id", Value: matchID}}

	updater := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "initial_ticket_amount", Value: match.InitialTicketAmount},
		primitive.E{Key: "external_id", Value: match.ExternalID},
		primitive.E{Key: "price", Value: match.Price},
		primitive.E{Key: "opponent", Value: match.Opponenent},
		primitive.E{Key: "league_name", Value: match.LeagueName},
		primitive.E{Key: "available_ticket_amount", Value: match.AvailableTicketAmount},
		primitive.E{Key: "away_match", Value: match.AwayMatch},
		primitive.E{Key: "location", Value: match.Location},
		primitive.E{Key: "date", Value: match.Date},
	}}}

	client, err := db.GetMongoClient()
	if err != nil {
		log.Errorf("Eror gettin db client: %v", err)
		return nil, err
	}
	collection := client.Database(db.DB).Collection(db.MATCHES)

	updateResult, err := collection.UpdateOne(context.TODO(), filter, updater)
	if err != nil {
		log.Errorf("Erorr getting collection of db: %v", err)
		return nil, err
	}

	if updateResult.ModifiedCount == 0 {
		return nil, fmt.Errorf("no document was updated")
	}

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Errorf("Erorr finding match in db: %v", err)
		return nil, err
	}
	log.Info("match updated: %v", match)
	return &result, nil
}

func UpdateTickets(matchID primitive.ObjectID, match *model.Match) (*model.Match, error) {
	result := model.Match{}
	filter := bson.D{primitive.E{Key: "_id", Value: matchID}}

	existingmatch, err := GetMatchByID(matchID)
	if err != nil {
		log.Errorf("Erorr getting match by id: %v", err)
		return &result, err
	}

	match.AvailableTicketAmount = existingmatch.AvailableTicketAmount + match.InitialTicketAmount
	match.InitialTicketAmount = existingmatch.InitialTicketAmount + match.InitialTicketAmount

	updater := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "initial_ticket_amount", Value: match.InitialTicketAmount},
		primitive.E{Key: "price", Value: match.Price},
		primitive.E{Key: "available_ticket_amount", Value: match.AvailableTicketAmount},
	}}}

	client, err := db.GetMongoClient()
	if err != nil {
		log.Errorf("Eror gettin db client: %v", err)
		return nil, err
	}
	collection := client.Database(db.DB).Collection(db.MATCHES)

	updateResult, err := collection.UpdateOne(context.TODO(), filter, updater)
	if err != nil {
		log.Errorf("Erorr getting collection of db: %v", err)
		return nil, err
	}

	if updateResult.ModifiedCount == 0 {
		return nil, fmt.Errorf("no document was updated")
	}

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	log.Info("tickets updated: %v", match)
	return &result, nil
}

func GetAllMatches() ([]model.Match, error) {
	filter := bson.D{{}}
	matches := []model.Match{}

	client, err := db.GetMongoClient()
	if err != nil {
		return matches, err
	}

	collection := client.Database(db.DB).Collection(db.MATCHES)
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return matches, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var match model.Match
		if err := cur.Decode(&match); err != nil {
			return matches, err
		}
		matches = append(matches, match)
	}

	if len(matches) == 0 {
		return matches, mongo.ErrNoDocuments
	}

	return matches, nil
}

func GetMatchByID(matchID primitive.ObjectID) (*model.Match, error) {
	result := model.Match{}
	filter := bson.D{primitive.E{Key: "_id", Value: matchID}}

	client, err := db.GetMongoClient()
	if err != nil {
		return &result, err
	}
	collection := client.Database(db.DB).Collection(db.MATCHES)

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func DeleteMatch(matchID primitive.ObjectID) error {
	filter := bson.D{primitive.E{Key: "_id", Value: matchID}}
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

func DeleteAllMatches() error {
	selector := bson.D{{}}
	client, err := db.GetMongoClient()
	if err != nil {
		return err
	}
	collection := client.Database(db.DB).Collection(db.MATCHES)
	_, err = collection.DeleteMany(context.TODO(), selector)
	if err != nil {
		return err
	}
	return nil
}
