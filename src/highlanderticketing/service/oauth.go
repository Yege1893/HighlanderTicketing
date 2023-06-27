package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func ValidateGoogleAccessToken(accessToken string) (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/tokeninfo?access_token="+accessToken, nil)
	if err != nil {
		return false, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var tokenInfo struct {
		ExpiresIn int    `json:"expires_in"`
		Error     string `json:"error"`
	}

	err = json.Unmarshal(body, &tokenInfo)
	if err != nil {
		return false, err
	}

	if tokenInfo.Error != "" {
		return false, fmt.Errorf("Fehler bei der Überprüfung des Tokens: %s", tokenInfo.Error)
	}

	if tokenInfo.ExpiresIn > 0 {
		return true, nil
	}

	return false, nil
}
