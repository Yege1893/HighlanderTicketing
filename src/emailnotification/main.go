package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	fmt.Println("aufgerufen")

	//testemailcontent := model.EmialContenct{Name: "Yannick", Location: "Stuttgart"}
	//receiver, body, subject := service.CreateEmail(testemailcontent, "confirmOrder")
	//service.SendEmail(receiver, body, subject)

	natsUrl := os.Getenv("NATS_URL")
	if nc, err := nats.Connect(natsUrl); err != nil {
		fmt.Println(err, nc)
	}

	go func() {
		fmt.Println("wird aufgerufen")
		nc, err := nats.Connect(natsUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer nc.Close()
		sub, err := nc.SubscribeSync("Order")
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Service hört auf die Subscription:", sub.Subject)

		for {
			msg, err := sub.NextMsg(0)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Nachricht erhalten: %s", msg.Data)

			// Führe hier deine spezifische Logik aus, um die Nachricht zu verarbeiten

			// Bestätige die Verarbeitung der Nachricht
			msg.Ack()
		}

	}()

	go func() {
		nc, err := nats.Connect(natsUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer nc.Close()
		sub, err := nc.SubscribeSync("Cancelation")
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Service hört auf die Subscription:", sub.Subject)

		for {
			msg, err := sub.NextMsg(0)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Nachricht erhalten: %s", msg.Data)

			// Führe hier deine spezifische Logik aus, um die Nachricht zu verarbeiten

			// Bestätige die Verarbeitung der Nachricht
			msg.Ack()
		}

	}()

}
