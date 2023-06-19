package handler

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/model"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/service"
)

func AddMatchOrder(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	order, err := getOrder(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = service.AddMatchOrder(id, order)
	if err != nil {
		log.Errorf("Failure adding donation to campaign with ID %v: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJson(w, order)

}
func AddTravelOrder(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	order, err := getOrder(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = service.AddTravelOrder(id, order)
	if err != nil {
		log.Errorf("Failure adding donation to campaign with ID %v: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJson(w, order)

}

func getOrder(r *http.Request) (*model.Order, error) {
	var order model.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		log.Errorf("Can't serialize request body to order struct: %v", err)
		return nil, err
	} else {
		log.Infof("request body seralized to order struct")
		log.Tracef("body seralized in struct order: %v", order)
	}
	return &order, nil
}
