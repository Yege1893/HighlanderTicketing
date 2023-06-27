package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/config"

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

	// Verwende das Token, um auf die Google API zuzugreifen oder speichere es für spätere Verwendung
	// token.AccessToken enthält den Zugriffstoken
	// token.RefreshToken enthält den Aktualisierungstoken
	fmt.Fprintf(w, "Token %s", token.AccessToken)
	fmt.Fprintf(w, "Token %s", token.RefreshToken)

	// Beispiel: Drucke den Namen des authentifizierten Benutzers
	client := oauthConfig.Client(context.Background(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Println("Fehler beim Abrufen der Benutzerinfo:", err)
		http.Error(w, "Fehler beim Abrufen der Benutzerinfo", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// Lese die Benutzerinfo als JSON-Daten
	var userinfo struct {
		Email string `json:"email"`
	}
	err = json.NewDecoder(response.Body).Decode(&userinfo)
	if err != nil {
		log.Println("Fehler beim Lesen der Benutzerinfo:", err)
		http.Error(w, "Fehler beim Lesen der Benutzerinfo", http.StatusInternalServerError)
		return
	}

	// Hier kannst du die E-Mail-Adresse des Benutzers verwenden
	fmt.Fprintf(w, "E-Mail-Adresse: %s", userinfo.Email)
}
