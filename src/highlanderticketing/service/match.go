package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/db"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateMatch(match *model.Match) error {
	match.ID = primitive.NewObjectID()
	match.Orders = []model.Order{}
	match.Travel.ID = primitive.NewObjectID()
	match.Travel.Orders = []model.Order{}
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

// noch testen nur intern für anbindung an die api
/*func CreateMatches(list *[]model.Match) error {
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
}*/

func UpdateMatch(matchID primitive.ObjectID, match *model.Match) (*model.Match, error) {
	result := model.Match{}
	existingMatch, err := GetMatchByID(matchID)
	if existingMatch == nil || err != nil {
		return existingMatch, err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: matchID}}

	updater := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "initial_ticket_amount", Value: match.InitialTicketAmount},
		primitive.E{Key: "available_ticket_amount", Value: match.AvailableTicketAmount},
		primitive.E{Key: "away_match", Value: match.AwayMatch},
		primitive.E{Key: "location", Value: match.Location},
		//primitive.E{Key: "date", Value: match.Date},
	}}}

	client, err := db.GetMongoClient()
	if err != nil {
		return nil, err
	}
	collection := client.Database(db.DB).Collection(db.MATCHES)

	updateResult, err := collection.UpdateOne(context.TODO(), filter, updater)
	if err != nil {
		return nil, err
	}

	if updateResult.ModifiedCount == 0 {
		return nil, fmt.Errorf("no document was updated")
	}

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

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
func GetMatchesOfApiToDb(apiUrl string) {
	data := getData(apiUrl)
	formatJsonCreateMatch(data)
}

func getData(apiUrl string) []byte {
	request, error := http.NewRequest("GET", apiUrl, nil)

	if error != nil {
		fmt.Println(error)
	}
	client := &http.Client{}
	response, error := client.Do(request)

	if error != nil {
		fmt.Println(error)
	}

	responseBody, error := io.ReadAll(response.Body)

	if error != nil {
		fmt.Println(error)
	}
	defer response.Body.Close()

	return responseBody
}

func formatJsonCreateMatch(jsonArray []byte) {
	var match model.Match
	var results []map[string]interface{}

	json.Unmarshal([]byte(jsonArray), &results)

	for _, result := range results {
		match.Date = result["matchDateTime"].(string)
		if team1, ok := result["team1"].(map[string]interface{}); ok {
			if name, ok := team1["teamName"].(string); ok {
				match.Location = name
			}
		}
		if team2, ok := result["team2"].(map[string]interface{}); ok {
			if name, ok := team2["teamName"].(string); ok {
				if name == "VfB Stuttgart" {
					match.AwayMatch = true
				}
			}
		}
		CreateMatch(&match)
	}
}
