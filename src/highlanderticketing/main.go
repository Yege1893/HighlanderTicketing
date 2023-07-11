package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/api"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/db"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/handler"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/model"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/service"
)

func main() {
	updateChan := make(chan *model.Match)
	service.DeleteAllMatches()
	service.DeleteAllUsers()
	/* var userArray []model.User
	userArray, _ = service.GetAllUsers()
	fmt.Println(userArray)*/

	go func() {
		for {
			err := api.GetlatestMatchesOfApi("https://api.openligadb.de/getmatchesbyteamid/16/10/0", updateChan)
			if err != nil {
				log.Println("Fehler beim Abrufen der Matches:", err)
			}
			time.Sleep(3 * time.Minute)
		}
	}()

	go func() {
		for {
			match := <-updateChan
			fmt.Println(match, "match")
			service.InserExternalMatch(match)
		}
	}()

	/*
		matches, errMatches := api.GetMatchesOfApi("https://api.openligadb.de/getmatchesbyteamid/16/10/0")
		if errMatches != nil {
			return
		}
		for _, match := range matches {
			service.CreateMatch(match)
		}
	*/

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	var natsServer service.NatsServer

	uri := os.Getenv("NATS_URI")

	nc, err := nats.Connect(uri)
	if err == nil {
		natsServer.Nc = nc
	}
	if err != nil {
		log.Fatal("Error establishing connection to NATS:", err)
	}

	log.Println("Starting Highlander Ticketing server")
	router := mux.NewRouter()
	router.HandleFunc("/register", handler.HandleRegister).Methods("GET")
	router.HandleFunc("/callback/register", handler.HandleCallbackRegister).Methods("GET")
	router.HandleFunc("/login", handler.HandleLogin).Methods("GET")
	router.HandleFunc("/callback/login", handler.HandleCallbackLogin).Methods("GET")
	router.HandleFunc("/health", handler.Health).Methods("GET")
	router.HandleFunc("/match", handler.CreateMatch).Methods("POST")
	router.HandleFunc("/matches", handler.GetAllMatches).Methods("GET")
	router.HandleFunc("/match/{id}", handler.GetMatchByID).Methods("GET")
	router.HandleFunc("/match/{id}", handler.UpdateMatch).Methods("PUT")
	router.HandleFunc("/match/{id}", handler.DeleteMatch).Methods("DELETE")
	router.HandleFunc("/match/{id}/updatetickets", handler.UpdateTickets).Methods("PUT")
	router.HandleFunc("/match/{id}/matchorder", handler.AddMatchOrder).Methods("POST")
	router.HandleFunc("/match/{id}/cancelorder/{orderid}", handler.CancelOrder).Methods("PUT")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatal(err)
	}

	err = db.CloseMongoClient()
	if err != nil {
		log.Fatal(err)
	}

}

func init() {
	//init db
	_, err := db.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}
	// init logger
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)

	os.Setenv("LOG_LEVEL", "INFO")

	level, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Info("Log level not specified, set default to: INFO")
		log.SetLevel(log.InfoLevel)
		return
	}
	log.SetLevel(level)
}
