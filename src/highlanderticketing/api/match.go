package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/model"
)

func GetMatchesOfApi(apiUrl string) ([]*model.Match, error) {
	data := getData(apiUrl)
	matches, err := formatJsonCreateMatch(data)
	if err != nil {
		return []*model.Match{}, err
	}
	fmt.Println(matches)
	return matches, nil
}

func getData(apiUrl string) []byte {
	request, error := http.NewRequest("GET", apiUrl, nil)

	if error != nil {
		fmt.Println(error)
	}
	client := &http.Client{}
	response, error := client.Do(request)

	if error != nil {
		fmt.Println(error)
	}

	responseBody, error := io.ReadAll(response.Body)

	if error != nil {
		fmt.Println(error)
	}
	defer response.Body.Close()

	return responseBody
}

func formatJsonCreateMatch(jsonArray []byte) ([]*model.Match, error) {
	var matches []*model.Match
	var results []map[string]interface{}

	if err := json.Unmarshal(jsonArray, &results); err != nil {
		return nil, err
	}

	for _, result := range results {
		var match model.Match
		match.ExternalID = int64(result["matchID"].(float64))
		match.LeagueName = result["leagueName"].(string)

		matchDate, err := time.Parse("2006-01-02T15:04:05", result["matchDateTime"].(string))
		if err != nil {
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
	return matches, nil
}
