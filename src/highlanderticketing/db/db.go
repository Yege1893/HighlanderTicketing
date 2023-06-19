package db

import (
	"context"
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
	CONNECTIONSTRING = "mongodb://localhost:27017"
	DB               = "db_issue_manager"
	MATCHES          = "col_matches"
	POOL_SIZE        = 10 // Anzahl der Verbindungen im Pool
)

func GetMongoClient() (*mongo.Client, error) {
	clientOnce.Do(func() {
		// Erstelle den Verbindungspool
		clientOptions := options.Client().ApplyURI(CONNECTIONSTRING)
		clientOptions.SetMaxPoolSize(POOL_SIZE)
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
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
