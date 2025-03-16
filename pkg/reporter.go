package reporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/arwenizer/zohan/config"
)

// ErrorEvent represents a captured error with additional context.
type ErrorEvent struct {
	Message    string            `json:"message"`
	StackTrace string            `json:"stack_trace"`
	Timestamp  time.Time         `json:"timestamp"`
	Labels     map[string]string `json:"labels,omitempty"`
}

// Reporter handles error reporting.
type Reporter struct {
	slackWebhookURL string
	sentryDSN       string // Reserved for future Sentry integration.
}

// NewReporter creates a new Reporter using the global configuration.
func NewReporter() *Reporter {
	cfg := config.GlobalConfig
	return &Reporter{
		slackWebhookURL: cfg.SlackWebhookURL,
		sentryDSN:       cfg.SentryDSN,
	}
}

// Report captures an error (with full stack trace), attaches labels,
// logs it to the console, and sends a notification to Slack if configured.
func (r *Reporter) Report(err error, labels map[string]string) {
	// Capture the full stack trace.
	stack := string(debug.Stack())
	event := ErrorEvent{
		Message:    err.Error(),
		StackTrace: stack,
		Timestamp:  time.Now(),
		Labels:     labels,
	}

	// Log the error locally.
	fmt.Printf("Error reported at %s:\nMessage: %s\nLabels: %v\nStack Trace:\n%s\n\n",
		event.Timestamp.Format(time.RFC3339),
		event.Message,
		event.Labels,
		event.StackTrace)

	// Send the error to Slack asynchronously, if a webhook URL is configured.
	if r.slackWebhookURL != "" {
		go r.sendToSlack(event)
	}

	// Optionally, add additional integrations (e.g., Sentry) here.
}

// sendToSlack sends the error event to Slack using Slack Block Kit formatting.
func (r *Reporter) sendToSlack(event ErrorEvent) {
	// Build a Slack payload using blocks.
	payload := map[string]interface{}{
		"blocks": []interface{}{
			map[string]interface{}{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("*Error Reported at %s*", event.Timestamp.Format(time.RFC3339)),
				},
			},
			map[string]interface{}{
				"type": "section",
				"fields": []interface{}{
					map[string]interface{}{
						"type": "mrkdwn",
						"text": fmt.Sprintf("*Message:*\n%s", event.Message),
					},
					map[string]interface{}{
						"type": "mrkdwn",
						"text": fmt.Sprintf("*Labels:*\n%s", formatLabels(event.Labels)),
					},
				},
			},
			map[string]interface{}{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("*Stack Trace:*\n```%s```", event.StackTrace),
				},
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Failed to marshal Slack payload: %v\n", err)
		return
	}

	resp, err := http.Post(r.slackWebhookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Failed to send error to Slack: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Non-OK response from Slack: %d\n", resp.StatusCode)
	}
}

// formatLabels converts a map of labels to a formatted string.
func formatLabels(labels map[string]string) string {
	if len(labels) == 0 {
		return "None"
	}
	var result string
	for key, value := range labels {
		result += fmt.Sprintf("*%s:* %s\n", key, value)
	}
	return result
}

// CapturePanic is a helper function that recovers from a panic and reports it.
// Use it with defer in your main function or goroutines.
func CapturePanic() {
	if rec := recover(); rec != nil {
		var err error
		switch v := rec.(type) {
		case error:
			err = v
		default:
			err = fmt.Errorf("%v", v)
		}
		labels := map[string]string{
			"panic": "true",
		}
		NewReporter().Report(err, labels)
	}
}
