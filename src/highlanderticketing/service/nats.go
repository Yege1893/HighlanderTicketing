package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/model"
)

type NatsServer struct {
	Nc *nats.Conn
}

func ConnectToNats() (NatsServer, error) {
	var natsServer NatsServer
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}
	uri := os.Getenv("NATS_URI")
	nc, err := nats.Connect(uri)
	if err != nil {
		log.Fatal("Error establishing connection to NATS:", err)
		return natsServer, err
	}
	natsServer.Nc = nc
	fmt.Println("Connected to NATS at:", natsServer.Nc.ConnectedUrl())
	return natsServer, nil

}

func (s NatsServer) ConfirmOrder(e *model.EmialContent) {
	var res *model.Response
	emailContenct, errMarshal := json.Marshal(e)
	if errMarshal != nil {
		fmt.Println(errMarshal)
		return
	}
	response, err := s.Nc.Request("confirmOrder."+string(e.OrderID), []byte(emailContenct), 2*time.Second)
	if err != nil {
		log.Println("Error making NATS request:", err)
	}

	if err := json.Unmarshal(response.Data, &res); err != nil {
		panic(err)
	}
	fmt.Println("hier die nats response", *res)
}

func (s NatsServer) confirmCancel(e *model.EmialContent) {
	var res *model.Response
	emailContenct, errMarshal := json.Marshal(e)
	if errMarshal != nil {
		fmt.Println(errMarshal)
		return
	}
	response, err := s.Nc.Request("confirmOrder."+string(e.Emailadress), []byte(emailContenct), 2*time.Second)
	if err != nil {
		log.Println("Error making NATS request:", err)
	}

	if err := json.Unmarshal(response.Data, &res); err != nil {
		panic(err)
	}
	fmt.Println("hier die nats response", &res)
}
