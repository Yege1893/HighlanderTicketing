package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/model"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/service"
)

func GetMatchesOfApiToDb(apiUrl string) {
	data := getData(apiUrl)
	formatJsonCreateMatch(data)
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

func formatJsonCreateMatch(jsonArray []byte) {
	var match model.Match
	var results []map[string]interface{}

	json.Unmarshal([]byte(jsonArray), &results)

	for _, result := range results {
		match.Date = result["matchDateTime"].(string)
		if team1, ok := result["team1"].(map[string]interface{}); ok {
			if name, ok := team1["teamName"].(string); ok {
				match.Location = name
			}
		}
		if team2, ok := result["team2"].(map[string]interface{}); ok {
			if name, ok := team2["teamName"].(string); ok {
				if name == "VfB Stuttgart" {
					match.AwayMatch = true
				}
			}
		}
		service.CreateMatch(&match)
	}
}
