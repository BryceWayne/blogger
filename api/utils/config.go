package utils

import (
	"fmt"
	"os"
)

type Config struct {
	JwtSecret string
	Env       string
	GCPCreds  string
}

func LoadConfig() (*Config, error) {
	jwtSecret, exists := os.LookupEnv("JWT_SECRET")
	if !exists {
		return nil, fmt.Errorf("JWT_SECRET not set")
	}

	env, exists := os.LookupEnv("ENV")
	if !exists {
		return nil, fmt.Errorf("ENV not set")
	}

	gcpCreds, exists := os.LookupEnv("GCP_CREDS")
	if !exists {
		return nil, fmt.Errorf("GCP_CREDS not set")
	}

	return &Config{
		JwtSecret: jwtSecret,
		Env:       env,
		GCPCreds:  gcpCreds,
	}, nil
}
