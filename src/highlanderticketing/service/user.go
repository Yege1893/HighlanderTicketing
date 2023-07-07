package service

import (
	"context"
	"fmt"

	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/db"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var isFirstCall = true

func CreateUser(user *model.User) error {
	if isFirstCall == false {
		user.IsAdmin = false
	} else if isFirstCall == true {
		_, err := GetAllUsers()
		if err == mongo.ErrNoDocuments {
			user.IsAdmin = true
		} else if err != nil {
			return err
		} else {
			user.IsAdmin = false
		}
		isFirstCall = false
	}
	user.ID = primitive.NewObjectID()
	client, err := db.GetMongoClient()
	if err != nil {
		return err
	}

	filter := bson.M{"email": user.Email}
	update := bson.M{
		"$setOnInsert": bson.M{
			"_id":         user.ID,
			"email":       user.Email,
			"google_id":   user.GoogleID,
			"name":        user.Name,
			"family_name": user.FamilyName,
			"is_admin":    user.IsAdmin,
		},
	}

	collection := client.Database(db.DB).Collection(db.USERS)
	options := options.FindOneAndUpdate().SetUpsert(true)

	result := collection.FindOneAndUpdate(context.TODO(), filter, update, options)

	if result.Err() == mongo.ErrNoDocuments {
		fmt.Println(3)
		return nil // dokument wurd erstellt
	} else if result.Err() != nil {
		fmt.Println(2)
		return result.Err() // fehler beim process an sich
	} else {
		fmt.Println(1)
		return fmt.Errorf("Der Benutzer existiert bereits")
	}
}

func UpdateUser(userID primitive.ObjectID, user *model.User) (*model.User, error) { //darf nur ein admin machen
	result := model.User{}
	existingUser, err := GetUserByID(userID)
	if existingUser == nil || err != nil {
		return existingUser, err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: userID}}

	updater := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "email", Value: user.Email},
		primitive.E{Key: "google_id", Value: user.GoogleID},
		primitive.E{Key: "name", Value: user.Name},
		primitive.E{Key: "family_name", Value: user.FamilyName},
		primitive.E{Key: "is_admin", Value: user.IsAdmin},
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

func GetAllUsers() ([]model.User, error) {
	filter := bson.D{{}}
	users := []model.User{}

	client, err := db.GetMongoClient()
	if err != nil {
		return users, err
	}

	collection := client.Database(db.DB).Collection(db.USERS)
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return users, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var user model.User
		if err := cur.Decode(&user); err != nil {
			return users, err
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		return users, mongo.ErrNoDocuments
	}

	return users, nil
}

func GetUserByID(userID primitive.ObjectID) (*model.User, error) {
	result := model.User{}
	filter := bson.D{primitive.E{Key: "_id", Value: userID}}

	client, err := db.GetMongoClient()
	if err != nil {
		return &result, err
	}
	collection := client.Database(db.DB).Collection(db.USERS)

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}
func GetUserByEmail(email string) (*model.User, error) {
	result := model.User{}
	filter := bson.D{primitive.E{Key: "email", Value: email}}

	client, err := db.GetMongoClient()
	if err != nil {
		return &result, err
	}
	collection := client.Database(db.DB).Collection(db.USERS)

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func DeleteUser(UserID primitive.ObjectID) error {
	filter := bson.D{primitive.E{Key: "_id", Value: UserID}}
	client, err := db.GetMongoClient()
	if err != nil {
		return err
	}
	collection := client.Database(db.DB).Collection(db.USERS)
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAllUsers() error {
	selector := bson.D{{}}
	client, err := db.GetMongoClient()
	if err != nil {
		return err
	}
	collection := client.Database(db.DB).Collection(db.USERS)
	_, err = collection.DeleteMany(context.TODO(), selector)
	if err != nil {
		return err
	}
	return nil
}
