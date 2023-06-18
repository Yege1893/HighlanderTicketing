package handler

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/model"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/service"
)

func CreateMatch(w http.ResponseWriter, r *http.Request) {
	var match *model.Match
	match, err := getMatch(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := service.CreateMatch(match); err != nil {
		log.Errorf("Error calling service CreateMatch: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJson(w, match)
}

// noch testen
/*func CreateMatches(w http.ResponseWriter, r *http.Request) {
	var match *model.Match
	match, err := getMatch(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := service.CreateMatch(match); err != nil {
		log.Errorf("Error calling service CreateMatch: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJson(w, match)
}
*/
func UpdateMatch(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		log.Errorf("Please parse in ID at the url %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	MatchToUpdate, err := getMatch(r)
	if err != nil {
		log.Errorf("Match not found %v", err)
		return
	}
	MatchUpdated, err := service.UpdateMatch(id, MatchToUpdate)
	if err != nil {
		log.Errorf("Campaign could not be updated %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sendJson(w, MatchUpdated)
}

func GetAllMatches(w http.ResponseWriter, r *http.Request) {
	matches, err := service.GetAllMatches()
	if err != nil {
		log.Errorf("Error calling service GetAllMatches: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJson(w, matches)
}

func GetMatchByID(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		log.Errorf("Please parse in ID at the url %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	campaign, err := service.GetMatchByID(id)
	if err != nil {
		log.Errorf("No Match with this ID %v", err)
		return
	}
	sendJson(w, campaign)
}

func DeleteMatch(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		log.Errorf("Please parse in ID at the url %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		log.Infof("ID to delete was found in struct")
	}
	err1 := service.DeleteMatch(id)
	if err1 != nil {
		log.Errorf("Match could not be deleted %v", err1)
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	} else {
		log.Infof("ID deleted")
		log.Tracef("ID: %v deleted", id)
	}
	sendJson(w, result{Success: "OK"})
}

// nur intern
/*func DeleteAllMatches(w http.ResponseWriter, r *http.Request) {
	err := service.DeleteAllMatches()
	if err != nil {
		log.Errorf("Match could not be deleted %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		log.Infof("Matches deleted")
	}
	sendJson(w, result{Success: "OK"})
}*/
func getMatch(r *http.Request) (*model.Match, error) {
	var match model.Match
	err := json.NewDecoder(r.Body).Decode(&match)
	if err != nil {
		log.Errorf("Can't serialize request body to campaign struct: %v", err)
		return nil, err
	} else {
		log.Infof("request body seralized to campaign struct")
		log.Tracef("body seralized in struct campaign: %v", match)
	}
	return &match, nil
}

/*func getMatches(r *http.Request) (*[]model.Match, error){
	var match model.Match
	err := json.NewDecoder(r.Body).Decode(&match)
	if err != nil {
		log.Errorf("Can't serialize request body to campaign struct: %v", err)
		return nil, err
	} else {
		log.Infof("request body seralized to campaign struct")
		log.Tracef("body seralized in struct campaign: %v", match)
	}
	return &matches, nil
}*/
