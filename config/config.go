package config

import (
	"errors"
)

// Config holds configuration settings for Zohan.
type Config struct {
	// Required: The environment in which the app is running (e.g., "development", "production").
	Environment string 
	// Optional: The log level (e.g., "DEBUG", "INFO", "WARN", "ERROR"). Defaults to "INFO" if not set.
	LogLevel string 
	// Optional: Sentry DSN for future integration.
	SentryDSN string 
	// Required: Slack Incoming Webhook URL for sending notifications.
	SlackWebhookURL string 
}

// GlobalConfig holds the application-wide configuration.
var GlobalConfig *Config

// Init initializes the global configuration using the provided parameters.
// It returns an error if required parameters are missing.
func Init(cfg *Config) error {
	if cfg.Environment == "" {
		return errors.New("Environment must be provided")
	}
	if cfg.SlackWebhookURL == "" {
		return errors.New("SlackWebhookURL must be provided")
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "INFO"
	}
	GlobalConfig = cfg
	return nil
}
