package otellogr

import (
	"errors"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
)

func TestLoggerImplements(t *testing.T) {
	assert.Implements(t, (*logr.Logger)(nil), Logger{})
}

func TestLoggerEnabled(t *testing.T) {
	l := NewLogger("")
	assert.True(t, l.Enabled())
}

func TestLoggerInfo(t *testing.T) {
	exporter := setupTestExporter()

	l := NewLogger("")
	l.Info("Some information", "key", "value")

	assert.Len(t, exporter.GetRecordedSpans(), 1)
	span := exporter.GetRecordedSpans()[0]
	attr := extractAttributesFromSpan(span)
	assert.Equal(t, "Some information", span.Name)
	assert.Equal(t, "value", attr["key"].AsString())
}

func TestLoggerInfoLevelTooHigh(t *testing.T) {
	exporter := setupTestExporter()

	l := NewLogger("").V(Error)
	l.Info("Some information", "key", "value")
	assert.Len(t, exporter.GetRecordedSpans(), 0)
}

func TestLoggerError(t *testing.T) {
	exporter := setupTestExporter()

	l := NewLogger("")
	l.Error(errors.New("error"), "An error occured", "key", "value")

	assert.Len(t, exporter.GetRecordedSpans(), 1)
	span := exporter.GetRecordedSpans()[0]
	attr := extractAttributesFromSpan(span)
	assert.Equal(t, "An error occured", span.Name)
	assert.Equal(t, "value", attr["key"].AsString())

	assert.Len(t, span.MessageEvents, 1)
	ev := span.MessageEvents[0]
	assert.Equal(t, "error", ev.Name)
}

func TestLoggerErrorLevelTooHigh(t *testing.T) {
	exporter := setupTestExporter()

	l := NewLogger("").V(Error + 1)
	l.Error(errors.New("error"), "An error occured", "key", "value")
	assert.Len(t, exporter.GetRecordedSpans(), 0)
}

func TestLoggerV(t *testing.T) {
	l := NewLogger("")
	l = l.V(Error)
	assert.Equal(t, Error, l.(Logger).level)
}

func TestLoggerWithValues(t *testing.T) {
	exporter := setupTestExporter()
	l := NewLogger("")
	l = l.WithValues("k", "v")

	l.Info("With Values")
	attr := extractAttributesFromSpan(exporter.GetRecordedSpans()[0])
	assert.Equal(t, "v", attr["k"].AsString())
}

func TestLoggerWithName(t *testing.T) {
	exporter := setupTestExporter()
	l := NewLogger("l")
	l = l.WithName("nl")

	l.Info("With Name")
	span := exporter.GetRecordedSpans()[0]
	assert.Equal(t, "l/nl", span.InstrumentationLibrary.Name)
}
