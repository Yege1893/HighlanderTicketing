package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// in env varialben packen
func GetOAuthConfigLogin() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     "20993146996-uvmr7479e3qmhiu9gasko3lr163ll76j.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8000/callback/login",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func GetOAuthConfigRegister() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     "20993146996-uvmr7479e3qmhiu9gasko3lr163ll76j.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8000/callback/register",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
