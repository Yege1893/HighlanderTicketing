package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type result struct {
	Success string `json:"success"`
}

func sendJson(w http.ResponseWriter, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Errorf("Failure encoding value to JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func getID(r *http.Request) (primitive.ObjectID, error) {
	vars := mux.Vars(r)
	id := vars["id"]

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("Can't get ObjectID from request: %v", err)
		return primitive.NilObjectID, err
	}

	return objectID, nil
}
func getBearerToken(r *http.Request) (string, error) {
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		log.Error("no Bearer Token in Request")
		return "", fmt.Errorf("Please parse in Bearer Token")

	}
	splitToken := strings.Split(reqToken, "Bearer")
	if len(splitToken) != 2 {
		log.Error("Beaerer Token could not be extracted")
		return "", fmt.Errorf("Can not extract Token")
	}

	reqToken = strings.TrimSpace(splitToken[1])
	return reqToken, nil
}
