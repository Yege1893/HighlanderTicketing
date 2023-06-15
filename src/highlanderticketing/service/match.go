package service

import (
	"context"

	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/db"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateMatch(match *model.Match) error {
	client, err := db.GetMongoClient()
	if err != nil {
		return err
	}
	collection := client.Database(db.DB).Collection(db.ISSUES)

	_, err = collection.InsertOne(context.TODO(), match)
	if err != nil {
		return nil
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

	collection := client.Database(db.DB).Collection(db.ISSUES)
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
