package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	fmt.Println("aufgerufen")

	time.Sleep(2 * time.Second)

	natsUrl := os.Getenv("NATS_URL")
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		log.Fatalf("Error connecting to NATS server: %v", err)
	}
	defer nc.Close()

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		fmt.Println("wird aufgerufen")
		orderSub, err := nc.SubscribeSync("Order")
		if err != nil {
			log.Fatalf("Error subscribing to 'Order': %v", err)
		}
		log.Println("Service hört auf die Subscription:", orderSub.Subject)

		for {
			msg, err := orderSub.NextMsg(-1)
			if err != nil {
				log.Println("Error receiving message:", err)
				continue
			}
			if msg != nil {
				log.Printf("Nachricht erhalten: %s", msg.Data)
				msg.Ack()
			}
		}

	}()

	go func() {
		defer wg.Done()
		cancelationSub, err := nc.SubscribeSync("Cancelation")
		if err != nil {
			log.Fatalf("Error subscribing to 'Cancelation': %v", err)
		}
		log.Println("Service hört auf die Subscription:", cancelationSub.Subject)

		for {
			msg, err := cancelationSub.NextMsg(-1)
			if err != nil {
				log.Println("Error receiving message:", err)
				continue
			}
			if msg != nil {
				log.Printf("Nachricht erhalten: %s", msg.Data)
				msg.Ack()
			}
		}

	}()
	wg.Wait()

	fmt.Println("Alle Go-Routinen sind abgeschlossen. Hauptprogramm wird beendet.")
}
