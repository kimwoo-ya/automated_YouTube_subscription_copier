package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

type Configuration struct {
	google_api_key       string
	target_channel_id    string
	google_client_id     string
	google_client_secret string
	redirect_url         string
	oauth_config         oauth2.Config
}

var instance *Configuration

func init() {
	instance = new()
	fmt.Printf("config.init():%+v\n", instance)
}

func new() *Configuration {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	google_client_id := os.Getenv("GOOGLE_CLIENT_ID")
	google_client_secret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirect_url := os.Getenv("TARGET_CHANNEL_ID")

	return &Configuration{
		google_api_key:       os.Getenv("GOOGLE_API_KEY"),
		target_channel_id:    os.Getenv("TARGET_CHANNEL_ID"),
		google_client_id:     google_client_id,
		google_client_secret: google_client_secret,
		redirect_url:         redirect_url,
		oauth_config: oauth2.Config{
			ClientID:     google_client_id,
			ClientSecret: google_client_secret,
			RedirectURL:  redirect_url,
			Scopes:       []string{youtube.YoutubeScope},
			Endpoint:     google.Endpoint,
		},
	}
}

func GetInstance() *Configuration {
	return instance
}

func (conf *Configuration) GetAPIKey() string {
	return conf.google_api_key
}
func (conf *Configuration) GetTargetChannelId() string {
	return conf.target_channel_id
}
func (conf *Configuration) GetGoogleClientId() string {
	return conf.google_client_id
}
func (conf *Configuration) GetClientSecret() string {
	return conf.google_client_secret
}
func (conf *Configuration) GetRedirectUrl() string {
	return conf.redirect_url
}

func (conf *Configuration) GetOauthConfig() oauth2.Config {
	return conf.oauth_config
}
