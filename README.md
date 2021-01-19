# Otel Logr
[![Go Reference](https://pkg.go.dev/badge/github.com/dmathieu/otel-logr.svg)](https://pkg.go.dev/github.com/dmathieu/otel-logr)
[![CircleCI](https://circleci.com/gh/dmathieu/otel-logr.svg?style=svg)](https://app.circleci.com/pipelines/github/dmathieu/otel-logr)

Implementation of the [logr](https://github.com/go-logr/logr) interface with [OpenTelemetry Go](https://github.com/open-telemetry/opentelemetry-go)

## Usage

```golang
logger := otellogr.NewLogger("Tracer Name")
logger = logger.WithAttributes("my-key", "my value") // Sets attributes for all spans created afterwards with this logger

logger.Info("This is some information") // Creates and ends a span with this name
logger.Info(errors.New("An error occured"), "This is some error") // Creates and ends a span with an error Event
```
