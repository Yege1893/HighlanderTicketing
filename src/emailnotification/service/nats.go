package service

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
)

func ConnectToNats() (*nats.Conn, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}
	uri := os.Getenv("NATS_URI")
	var err error
	var nc *nats.Conn
	nc, err = nats.Connect(uri)
	if err != nil {
		log.Fatal("Error establishing connection to NATS:", err)
	}
	return nc, nil
}
