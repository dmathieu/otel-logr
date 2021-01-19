package otellogr

import (
	"context"
	"sync"

	sdktrace "go.opentelemetry.io/otel/sdk/export/trace"
)

// OpenTelemetry exporter for collecting spans in memory for testing
type testExporter struct {
	mu    sync.Mutex
	spans []*sdktrace.SpanSnapshot
}

func (t *testExporter) ExportSpans(ctx context.Context, s []*sdktrace.SpanSnapshot) error {
	t.mu.Lock()
	t.spans = append(t.spans, s...)
	t.mu.Unlock()

	return nil
}

func (t *testExporter) GetRecordedSpans() []*sdktrace.SpanSnapshot {
	return t.spans
}

func (t *testExporter) ClearRecordedSpans() {
	t.spans = []*sdktrace.SpanSnapshot{}
}

func (t *testExporter) Shutdown(ctx context.Context) error {
	return nil
}
