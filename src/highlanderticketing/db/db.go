package db

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"

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
	USERS_TEST       = "col_users_test"
	POOL_SIZE        = 10000 // Anzahl der Verbindungen im Pool
)

func GetMongoClient() (*mongo.Client, error) {
	clientOnce.Do(func() {
		// Erstelle den Verbindungspool
		clientOptions := options.Client().ApplyURI(CONNECTIONSTRING)
		clientOptions.SetMaxPoolSize(POOL_SIZE)
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {

			clientInstanceError = err
			log.Errorf("Failure instancing client %v", err)
		}
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
			log.Errorf("Failure pinging client %v", err)
		}
		clientInstance = client
	})
	log.Info("client returned")
	return clientInstance, clientInstanceError
}
func CloseMongoClient() error {
	if clientInstance != nil {
		err := clientInstance.Disconnect(context.Background())
		if err != nil {
			log.Errorf("Failure disconnecting client %v", err)
			return err
		}
	}
	log.Info("client disconnected")
	return nil
}
