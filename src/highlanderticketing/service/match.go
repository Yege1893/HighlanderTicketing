package service

import (
	"context"

	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/db"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateMatch(match *model.Match) error {
	client, err := db.GetMongoClient()
	if err != nil {
		return err
	}
	collection := client.Database(db.DB).Collection(db.MATCHES)

	_, err = collection.InsertOne(context.TODO(), match)
	if err != nil {
		return nil
	}
	return nil
}

// noch testen
func CreateMatches(list *[]model.Match) error {
	insertableList := make([]interface{}, len(*list))
	for i, v := range *list {
		insertableList[i] = v
	}

	client, err := db.GetMongoClient()
	if err != nil {
		return err
	}

	collection := client.Database(db.DB).Collection(db.MATCHES)

	_, err = collection.InsertMany(context.TODO(), insertableList)
	if err != nil {
		return err
	}

	return nil
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
