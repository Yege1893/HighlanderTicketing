package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/model"
)

type NatsServer struct {
	Nc *nats.Conn
}

func (s NatsServer) ConfirmOrder(e *model.EmialContent) {
	emailContenct, errMarshal := json.Marshal(e)
	if errMarshal != nil {
		fmt.Println(errMarshal)
		return
	}
	response, err := s.Nc.Request("confirmOrder", []byte(emailContenct), 2*time.Second)
	if err != nil {
		log.Println("Error making NATS request:", err)
	}
	fmt.Println("hier die nats response", string(response.Data))
}

// hier dann confirm cancel
