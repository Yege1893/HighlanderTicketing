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
	/*service.DeleteAllUsers()
	var userArray []model.User
	userArray, _ = service.GetAllUsers()
	fmt.Println(userArray)
	service.DeleteAllMatches()
	api.GetMatchesOfApiToDb("https://api.openligadb.de/getmatchesbyteamid/16/5/0")*/
	//init db
	_, err := db.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Starting Highlander Ticketing server")
	router := mux.NewRouter()
	router.HandleFunc("/login", handler.HandleLogin).Methods("GET")
	router.HandleFunc("/callback", handler.HandleCallback).Methods("GET")
	router.HandleFunc("/health", handler.Health).Methods("GET")
	router.HandleFunc("/match", handler.CreateMatch).Methods("POST")
	router.HandleFunc("/matches", handler.GetAllMatches).Methods("GET")
	router.HandleFunc("/match/{id}", handler.GetMatchByID).Methods("GET")
	router.HandleFunc("/match/{id}", handler.UpdateMatch).Methods("PUT")
	router.HandleFunc("/match/{id}", handler.DeleteMatch).Methods("DELETE")
	router.HandleFunc("/match/{id}/matchorder", handler.AddMatchOrder).Methods("POST")
	router.HandleFunc("/match/{id}/travelorder", handler.AddTravelOrder).Methods("POST")
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
