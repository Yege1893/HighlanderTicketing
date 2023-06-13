package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/handler"
)

func main() {
	router := mux.NewRouter()
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
