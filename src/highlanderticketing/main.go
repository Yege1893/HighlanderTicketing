package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/db"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/handler"
)

func main() {
	//init db
	_, err := db.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Starting Highlander Ticketing server")
	router := mux.NewRouter()
	router.HandleFunc("/health", handler.Health).Methods("GET")
	router.HandleFunc("/match", handler.CreateMatch).Methods("POST")
	router.HandleFunc("/matches", handler.GetAllMatches).Methods("GET")
	router.HandleFunc("/match/{id}", handler.GetMatchByID).Methods("GET")
	router.HandleFunc("/match/{id}", handler.UpdateMatch).Methods("PUT")
	router.HandleFunc("/match/{id}", handler.DeleteMatch).Methods("DELETE")
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
