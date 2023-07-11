package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/model"
)

func GetMatchesOfApi(apiUrl string) ([]*model.Match, error) {
	data, err := getData(apiUrl)
	if err != nil {
		log.Errorf("Failure loading data of match %v", err)
		return []*model.Match{}, err
	}
	matches, err := formatJsonToMatch(data)
	if err != nil {
		log.Errorf("Failure formating match %v", err)
		return []*model.Match{}, err
	}
	log.Info("Got matches of API")
	return matches, nil
}

func GetlatestMatchesOfApi(url string, updateChan chan<- *model.Match) error {
	data, err := getData(url)
	if err != nil {
		log.Errorf("Failure loading data of match %v", err)
		return err
	}
	matches, err := formatJsonToMatch(data)
	if err != nil {
		log.Errorf("Failure formating match %v", err)
		return err
	}
	for _, match := range matches {
		updateChan <- match
	}
	log.Info("added matches to chain")
	return nil
}

func getData(apiUrl string) ([]byte, error) {
	request, err := http.NewRequest("GET", apiUrl, nil)

	if err != nil {
		log.Errorf("Failure creating request %v", err)
		return []byte{}, err
	}
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		log.Errorf("Failure making request %v", err)
		return []byte{}, err
	}

	responseBody, err := io.ReadAll(response.Body)

	if err != nil {
		log.Errorf("Failure searalizing request %v", err)
		return []byte{}, err
	}
	defer response.Body.Close()
	log.Info("Got data of API")
	return responseBody, nil
}

func formatJsonToMatch(jsonArray []byte) ([]*model.Match, error) {
	var matches []*model.Match
	var results []map[string]interface{}

	if err := json.Unmarshal(jsonArray, &results); err != nil {
		log.Errorf("Failure unmarshaling data %v", err)
		return nil, err
	}

	for _, result := range results {
		var match model.Match
		match.ExternalID = int64(result["matchID"].(float64))
		match.LeagueName = result["leagueName"].(string)

		matchDate, err := time.Parse("2006-01-02T15:04:05", result["matchDateTime"].(string))
		if err != nil {
			log.Errorf("Failure parsing date %v", err)
			return nil, err
		}
		match.Date = matchDate

		if team1, ok := result["team1"].(map[string]interface{}); ok {
			if name, ok := team1["shortName"].(string); ok {
				match.Location = name
			}
		}

		if team1, ok := result["team1"].(map[string]interface{}); ok {
			if name, ok := team1["teamName"].(string); ok {
				if name != "VfB Stuttgart" {
					match.Opponenent = name
				}
			}
		}

		if team2, ok := result["team2"].(map[string]interface{}); ok {
			if name, ok := team2["teamName"].(string); ok {
				if name == "VfB Stuttgart" {
					match.AwayMatch = true
				} else {
					match.AwayMatch = false
					match.Opponenent = name
				}
			}
		}
		matches = append(matches, &match)
	}
	log.Info("formated matches in json")
	return matches, nil
}
