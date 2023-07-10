package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/model"
)

func GetMatchesOfApi(apiUrl string) (error, []*model.Match) {
	data := getData(apiUrl)
	err, matches := formatJsonCreateMatch(data)
	if err != nil {
		return err, make([]*model.Match, 0)
	}
	fmt.Println(matches)
	return nil, matches
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

func formatJsonCreateMatch(jsonArray []byte) (error, []*model.Match) {
	var match model.Match
	var matches []*model.Match
	var results []map[string]interface{}

	json.Unmarshal([]byte(jsonArray), &results)

	for _, result := range results {
		match.ExternalID = int64(result["matchID"].(float64))
		match.LeagueName = result["leagueName"].(string)
		matchDate, err := time.Parse("2006-01-02T15:04:05", result["matchDateTime"].(string))
		if err != nil {
			return err, matches
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
	return nil, matches
}
