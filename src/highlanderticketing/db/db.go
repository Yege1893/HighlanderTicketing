package db

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance      *mongo.Client
	clientInstanceError error
	clientOnce          sync.Once
)

const (
	CONNECTIONSTRING = "mongodb://mongo:27017"
	DB               = "db_issue_manager"
	DBUSER           = "db_user"
	MATCHES          = "col_matches"
	USERS            = "col_users"
	POOL_SIZE        = 10000 // Anzahl der Verbindungen im Pool
)

func GetMongoClient() (*mongo.Client, error) {
	clientOnce.Do(func() {
		// Erstelle den Verbindungspool
		clientOptions := options.Client().ApplyURI(CONNECTIONSTRING)
		clientOptions.SetMaxPoolSize(POOL_SIZE)
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			fmt.Println("hier liegt der fehler")
			clientInstanceError = err
		}
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
		}
		clientInstance = client
	})

	return clientInstance, clientInstanceError
}
func CloseMongoClient() error {
	if clientInstance != nil {
		err := clientInstance.Disconnect(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}
