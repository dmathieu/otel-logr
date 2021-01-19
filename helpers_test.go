package otellogr

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func extractAttributesFromSpan(s *export.SpanSnapshot) map[string]label.Value {
	data := map[string]label.Value{}
	for _, kv := range s.Attributes {
		data[string(kv.Key)] = kv.Value
	}

	return data
}

func setupTestExporter() *TestExporter {
	tp := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(tp)
	exporter := new(TestExporter)
	tp.RegisterSpanProcessor(sdktrace.NewSimpleSpanProcessor(exporter))

	return exporter
}
