package config

import (
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2"
	"os"
)

var GoogleOAuthConfig *oauth2.Config 

func InitGoogleOAuth() {
	GoogleOAuthConfig = &oauth2.Config{
        ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
        ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        RedirectURL: "https://finbangla-voice-backend-production.up.railway.app/auth/google/callback",
        Scopes: []string{
            "https://www.googleapis.com/auth/userinfo.email",
            "https://www.googleapis.com/auth/userinfo.profile",
        },
        Endpoint: google.Endpoint,
	}
}