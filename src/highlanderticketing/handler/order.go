package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/model"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/service"
)

func AddMatchOrder(w http.ResponseWriter, r *http.Request) {
	if err := CheckAccessToken(w, r, false); err != nil {
		log.Errorf("Eror checking AccessToken: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
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
	reqToken, err := getBearerToken(r)
	if err != nil {
		log.Errorf("Failure loading bearer token  %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fmt.Println(reqToken, "reqToken")
	userOfOrder, err := service.GetUserInfoByToken(reqToken)
	if err != nil {
		log.Errorf("Failure loading user Info %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println(userOfOrder, "userOfOrder")
	internalUser, err := service.GetUserByEmail(userOfOrder.Email)
	if err != nil {
		log.Errorf("Failure loading internal user Info %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	order.User = *internalUser
	err = service.AddMatchOrder(id, order)
	if err != nil {
		log.Errorf("Failure adding order to match with ID %v: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJson(w, order)

}
func AddTravelOrder(w http.ResponseWriter, r *http.Request) {
	if err := CheckAccessToken(w, r, false); err != nil {
		log.Errorf("Eror checking AccessToken: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
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
	reqToken, err := getBearerToken(r)
	if err != nil {
		log.Errorf("Failure loading bearer token  %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	userOfOrder, err := service.GetUserInfoByToken(reqToken)
	if err != nil {
		log.Errorf("Failure loading user Info %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	internalUser, err := service.GetUserByEmail(userOfOrder.Email)
	if err != nil {
		log.Errorf("Failure loading internal user Info %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	order.User = *internalUser

	err = service.AddTravelOrder(id, order)
	if err != nil {
		log.Errorf("Failure adding donation to campaign with ID %v: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJson(w, order)

}
func CancelOrder(w http.ResponseWriter, r *http.Request) {
	if err := CheckAccessToken(w, r, false); err != nil {
		log.Errorf("Eror checking AccessToken: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	id, err := getID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	orderId, err := getOrderID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	order, err := getOrder(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	reqToken, err := getBearerToken(r)
	if err != nil {
		log.Errorf("Failure loading bearer token  %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	userOfOrder, err := service.GetUserInfoByToken(reqToken)
	if err != nil {
		log.Errorf("Failure loading user Info %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	internalUser, err := service.GetUserByEmail(userOfOrder.Email)
	if err != nil {
		log.Errorf("Failure loading internal user Info %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	order.User = *internalUser
	order.ID = orderId
	fmt.Println(orderId, "orderId")

	err = service.CancelOrder(id, order)
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
