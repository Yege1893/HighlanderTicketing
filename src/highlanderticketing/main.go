package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/db"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/handler"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/model"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	router := mux.NewRouter()
	_, err := db.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}

	// Testen
	objectID := primitive.NewObjectID()
	var initialMatch = model.Match{ID: objectID, InitialTicketAmount: 1, AvailableTicketAmount: 1, AwayMatch: true, Location: "Stuttgart"}
	err1 := service.CreateMatch(&initialMatch)
	if err1 != nil {
		fmt.Println(err)
	}
	matches, err := service.GetAllMatches()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(matches)

	var updateModel = model.Match{ID: objectID, InitialTicketAmount: 1, AvailableTicketAmount: 1, AwayMatch: true, Location: "Schalke"}
	updatedmatch, err := service.UpdateMatch(objectID, &updateModel)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(updatedmatch)
	match, err := service.GetMatchByID(objectID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(match)
	deleted := service.DeleteAllMatches()
	if err != nil {
		fmt.Println(deleted)
	}

	// ende tests

	router.HandleFunc("/health", handler.Health).Methods("GET")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatal(err)
	}

}

func init() {
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
