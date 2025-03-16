# Zohan

Zohan is a production-ready Go package for internal error reporting. It captures errors (and panics) with full stack traces, allows you to attach labels (similar to Sentry), and sends Slack notifications asynchronously.

## Features

- **Automatic Error Capture:**  
  Captures errors with full stack traces.
  
- **Automatic Panic Handling:**  
  Use `defer CapturePanic()` in your main function and goroutines to automatically capture and report panics.

- **Labeling:**  
  Attach labels (e.g., module, severity, environment) to errors for categorization.

- **Asynchronous Slack Notifications:**  
  Sends error notifications to Slack using Slack Block Kit formatting without blocking application flow.

- **Modularity and Extensibility:**  
  The project is split into two packages:  
  - `config`: Centralized configuration management (parameters are provided during Init).  
  - `reporter`: Core error reporting functionality. Easily extendable (e.g., adding Sentry integration).

## Production-Ready Considerations

- **Required Configuration:**  
  The `Init` function in the config package requires `Environment` and `SlackWebhookURL`. Other parameters have defaults or are optional.
  
- **Centralized Configuration:**  
  Configuration is managed in the `config` package. In production, you may extend it to load from files or secret managers.

- **Asynchronous Operations:**  
  Slack notifications are sent asynchronously, ensuring that error reporting does not block your application.

- **Automatic Panic Handling:**  
  Use `defer CapturePanic()` to automatically log and report panics immediately.

## Usage

1. **Import Zohan in your project:**

   ```go
   import (
       "github.com/arwenizer/zohan/config"
       "github.com/arwenizer/zohan/reporter"
   )
