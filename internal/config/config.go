package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	Port    string
	BaseURL string
	PG      *PGConfig
	Email   *EmailConfig
	Auth    *AuthConfig
}

type PGConfig struct {
	User       string
	Password   string
	Host       string
	Port       string
	DBName     string
	SSLMode    string
	DriverName string
}

type EmailConfig struct {
	FromAddress      string
	SMTPPassword     string
	SMTPHost         string
	SMTPPort         string
	ResetPasswordURL string
	ConfirmEmailURL  string
}

type AuthConfig struct {
	AccessTokenKey       []byte
	ActivationTokenKey   []byte
	OAuthClientID        string
	OAuthClientSecret    string
	OAuthUserInfoURL     string
	PostOAuthRedirectURL string
	OAuthConfig          oauth2.Config
}

var App *Config

func Load() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	baseURL := getEnv("BASE_URL")

	App = &Config{
		Port:    getEnv("PORT"),
		BaseURL: baseURL,

		PG: &PGConfig{
			User:       getEnv("PG_USER"),
			Password:   getEnv("PG_PASSWORD"),
			Host:       getEnv("PG_HOST"),
			Port:       getEnv("PG_PORT"),
			DBName:     getEnv("PG_DBNAME"),
			SSLMode:    getEnv("PG_SSL_MODE"),
			DriverName: getEnv("PG_DRIVER_NAME"),
		},
		Email: &EmailConfig{
			FromAddress:      getEnv("SMTP_FROM"),
			SMTPPassword:     getEnv("SMTP_PASSWORD"),
			SMTPHost:         getEnv("SMTP_HOST"),
			SMTPPort:         getEnv("SMTP_PORT"),
			ResetPasswordURL: baseURL + "/reset-password",
			ConfirmEmailURL:  baseURL + "/auth/confirm",
		},
		Auth: &AuthConfig{
			AccessTokenKey:       []byte(getEnv("ACCESS_TOKEN_SECRET")),
			ActivationTokenKey:   []byte(getEnv("ACTIVATION_TOKEN_SECRET")),
			OAuthClientID:        getEnv("OAUTH_CLIENT_ID"),
			OAuthClientSecret:    getEnv("OAUTH_CLIENT_SECRET"),
			OAuthUserInfoURL:     getEnv("OAUTH_USER_INFO_URL"),
			PostOAuthRedirectURL: baseURL + "/lobby",
			OAuthConfig: oauth2.Config{
				ClientID:     getEnv("OAUTH_CLIENT_ID"),
				ClientSecret: getEnv("OAUTH_CLIENT_SECRET"),
				RedirectURL:  baseURL + "/auth/google",
				Scopes:       []string{"openid", "profile", "email"},
				Endpoint:     google.Endpoint,
			},
		},
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return val
}

func (pg *PGConfig) PgConnString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		pg.User, pg.Password, pg.Host, pg.Port, pg.DBName, pg.SSLMode)
}
