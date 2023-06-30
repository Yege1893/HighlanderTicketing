package handler

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/config"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/service"

	"golang.org/x/oauth2"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	oauthConfig := config.GetOAuthConfig()
	url := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleCallback(w http.ResponseWriter, r *http.Request) {
	oauthConfig := config.GetOAuthConfig()
	code := r.URL.Query().Get("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Println("Fehler beim Austausch des Autorisierungscodes:", err)
		http.Error(w, "Fehler beim Authentifizieren", http.StatusInternalServerError)
		return
	}
	service.Register(token.AccessToken)
	sendJson(w, token.AccessToken)
}
