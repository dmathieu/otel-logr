package otellogr

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestTestExporter(t *testing.T) {
	tp := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(tp)
	exporter := new(testExporter)
	tp.RegisterSpanProcessor(sdktrace.NewSimpleSpanProcessor(exporter))

	assert.Empty(t, exporter.GetRecordedSpans())
	_, span := tp.Tracer("").Start(context.Background(), "test")
	span.End()
	assert.Len(t, exporter.GetRecordedSpans(), 1)
}
