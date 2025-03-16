package main

import (
	"errors"
	"log"
	"time"

	"github.com/arwenizer/zohan/config"
	"github.com/arwenizer/zohan/pkg"
)

func main() {
	// Initialize configuration with required parameters.
	err := config.Init(&config.Config{
		Environment:     "production", // Required
		LogLevel:        "DEBUG",      // Optional; defaults to INFO if empty.
		SlackWebhookURL: "https://hooks.slack.com/services/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", // Required for Slack notifications.
		SentryDSN:       "",           // Optional (for future Sentry integration).
	})
	if err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	// Set up automatic panic capture.
	defer reporter.CapturePanic()

	// Create a Reporter.
	rep := reporter.NewReporter()

	// Simulate an error.
	err = errors.New("simulated error for demonstration")
	labels := map[string]string{
		"module":   "payment",
		"severity": "critical",
		"env":      config.GlobalConfig.Environment,
	}

	// Report the error.
	rep.Report(err, labels)

	// Simulate a panic in a goroutine.
	go func() {
		defer reporter.CapturePanic()
		panic("simulated panic in goroutine")
	}()

	// Allow time for asynchronous operations.
	time.Sleep(3 * time.Second)
}
