package reporter

import (
	"errors"
	"testing"
	"time"

	"github.com/arwenizer/zohan/config"
)

func initTestConfig() {
	// Initialize configuration for tests.
	config.GlobalConfig = &config.Config{
		Environment:     "test",
		LogLevel:        "DEBUG",
		SentryDSN:       "",
		SlackWebhookURL: "http://invalid", // Dummy URL for testing asynchronous function.
	}
}

func TestReportLogsError(t *testing.T) {
	initTestConfig()
	rep := NewReporter()
	err := errors.New("test error")
	// Call Report and ensure it doesn't panic.
	rep.Report(err, map[string]string{"module": "unittest"})
}

func TestCapturePanic(t *testing.T) {
	initTestConfig()

	// Use a flag to check whether a panic escaped the inner function.
	panicEscaped := false

	func() {
		defer func() {
			if r := recover(); r != nil {
				// If we reach here, it means the panic was not recovered by CapturePanic.
				panicEscaped = true
			}
		}()

		// This inner function should recover the panic via CapturePanic.
		func() {
			defer CapturePanic()
			panic("simulated panic for testing")
		}()
	}()

	if panicEscaped {
		t.Errorf("CapturePanic did not recover the panic as expected")
	} else {
		t.Log("CapturePanic successfully recovered the panic")
	}

	// Also test with an error type panic.
	panicEscaped = false
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicEscaped = true
			}
		}()
		func() {
			defer CapturePanic()
			panic(errors.New("simulated error panic"))
		}()
	}()

	if panicEscaped {
		t.Errorf("CapturePanic did not recover the error panic as expected")
	} else {
		t.Log("CapturePanic successfully recovered the error panic")
	}

	// Allow some time for asynchronous operations (if any) to complete.
	time.Sleep(100 * time.Millisecond)
}

