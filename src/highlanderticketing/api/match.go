package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/model"
)

func GetMatchesOfApi(apiUrl string) []*model.Match {
	data := getData(apiUrl)
	fmt.Println(data)
	matches := formatJsonToMatches(data)
	return matches
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

func formatJsonToMatches(jsonArray []byte) []*model.Match {
	var match *model.Match
	var matches []*model.Match
	var results []map[string]interface{}

	json.Unmarshal([]byte(jsonArray), &results)

	for _, result := range results {
		fmt.Println(result, "result")
		match.ExternalID = int64(result["matchID"].(float64))
		fmt.Println(match.ExternalID)
		match.LeagueName = result["leagueName"].(string)
		match.Date = result["matchDateTime"].(time.Time)
		fmt.Println(*match)

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
		fmt.Println(&matches)
		matches = append(matches, match)
	}
	return matches
}
