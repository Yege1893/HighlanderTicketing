package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/config"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/service"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
)

var secretKey = []byte("mysecretkey")

// aufgerufen, wenn der benutzer sich ein loggt, ruft oauth conf auf, leitet zu authorisisuerngsseite um
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	oauthConfig := config.GetOAuthConfigLogin()
	url := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// aufgerufen, wenn der benutzer sich registriert, ruft oauth conf auf, leitet zu authorisisuerngsseite um
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	oauthConfig := config.GetOAuthConfigRegister()
	url := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// aufgerufen, wenn autorisierungscode von google kommt, tauscht code gegen access token ein, und registriert den benutzer
func HandleCallbackRegister(w http.ResponseWriter, r *http.Request) {
	oauthConfig := config.GetOAuthConfigRegister()
	code := r.URL.Query().Get("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Errorf("Failure exchanging authorisingcode: %v", err)
		http.Error(w, "Fehler beim Authentifizieren", http.StatusInternalServerError)
		return
	}
	err = service.Register(token.AccessToken)
	if err != nil {
		log.Errorf("Failure registrating: %v", err)
		io.WriteString(w, `user besteht bereits`)
	} else {
		sendJson(w, "user erfolgreich angelegt")
	}

}

// aufgerufen, wenn autorisierungscode von google kommt, tauscht code gegen Zugriffstoken ein, generiert jwt token
func HandleCallbackLogin(w http.ResponseWriter, r *http.Request) {

	oauthConfig := config.GetOAuthConfigLogin()
	code := r.URL.Query().Get("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Errorf("Failure exchanging authorisingcode: %v", err)
		http.Error(w, "Fehler beim Authentifizieren", http.StatusInternalServerError)
		return
	}

	user, err := service.GetUserInfoByToken(token.AccessToken)
	if err != nil {
		sendJson(w, err)
		return
	}
	userfound, errUser := service.GetUserByEmail(user.Email)
	if errUser != nil {
		sendJson(w, err)
		sendJson(w, "user nicht registriert")
		return
	}

	tokenJwt := jwt.New(jwt.SigningMethodHS256)
	claims := tokenJwt.Claims.(jwt.MapClaims)
	claims["username"] = userfound.Email
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()

	tokenString, err := tokenJwt.SignedString(secretKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Info("Login erfolgreich")
	sendJson(w, tokenString)
}

// überprüft den jwt token
func CheckAccessToken(w http.ResponseWriter, r *http.Request, needAdmin bool) (error, string) {
	tokenString, err := getBearerToken(r)
	if err != nil {
		log.Errorf("no bearer token handed over: %v", err)
		return err, ""
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Ungültiges Authorization-Token")
		return err, ""
	}
	var username string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username = claims["username"].(string)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Ungültiges Authorization-Token")
	}
	if needAdmin {
		err := checkAdmin(username)
		if err != nil {
			return err, ""
		}
	}
	log.Info("token checked succesfully")
	return nil, username
}

func checkAdmin(userEmail string) error {
	user, err := service.GetUserByEmail(userEmail)
	if err != nil {
		log.Errorf("failure getting user by email: %v", err)
		return err
	}
	if user.IsAdmin {
		return nil
	} else {
		return fmt.Errorf("User has not adminrights")
	}
}
