package config

import (
	"encoding/json"
	"os"
	"sync"
)

// Database connection config
type Database struct {
	User     string `json:"User"`
	Password string `json:"Password"`
	Hostname string `json:"Hostname"`
	Port     int    `json:"Port"`
	Name     string `json:"Name"`
	Type     string `json:"Type"`
}

// Application config
type Config struct {
	JwtSecret               string   `json:"jwtSecret"`
	JwtCookieKey            string   `json:"jwtCookieKey"`
	GoogleClientID          string   `json:"googleClientID"`
	GoogleClientSecret      string   `json:"googleClientSecret"`
	GoogleAuthRedirectURL   string   `json:"googleAuthRedirectURL"`
	GoogleOauthStateString  string   `json:"googleOauthStateString"`
	Database                Database `json:"Database"`
	PORT                    int      `json:"PORT"`
	ClientProfilePage       string   `json:"clientProfilePage"`
	ClientAuthErrorCallback string   `json:"clientAuthErrorCallback"`
	SendGridApiKey          string   `json:"SENDGRID_API_KEY"`
}

var (
	once     sync.Once
	instance *Config
)

// LoadConfig loads the configuration from a JSON file
func LoadConfig() error {
	var err error
	once.Do(func() {
		instance, err = loadConfig()
	})
	return err
}

func loadConfig() (*Config, error) {
	config := &Config{}

	// Read the configuration file
	data, err := os.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON data into Config struct
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// GetConfig returns the loaded configuration
func GetConfig() *Config {
	return instance
}
