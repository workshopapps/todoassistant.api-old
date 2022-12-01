package utils

import (
	"log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Config file stores configuration of application

// Values read by viper viper from a config file or enviroment variable


type GoogleConfig struct{
	GoogleLoginConfig oauth2.Config
}

var LoginConfig GoogleConfig
const GoogleAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="


func LoadGoogleConfig(){

	config, err := LoadConfig("./")

	if err != nil {
		log.Fatal("cannot load config", err)
	}

	googleSecret := config.GoogleSecret
	googleClient := config.GoogleClient
	googleCallBack := config.GoogleCallBack

	LoginConfig.GoogleLoginConfig = oauth2.Config{
		ClientID: googleSecret,
		ClientSecret: googleClient,
		Endpoint:     google.Endpoint,
		RedirectURL:  googleCallBack,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}
}


