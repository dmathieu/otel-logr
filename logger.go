package otellogr

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

const (
	All int = iota
	Info
	Error
)

type Logger struct {
	Name   string
	Tracer trace.Tracer

	values []label.KeyValue
	level  int
}

func NewLogger(n string) logr.Logger {
	return Logger{
		Name:   n,
		Tracer: otel.Tracer(n),
	}
}

func (l Logger) Enabled() bool {
	return true
}

func (l Logger) Info(msg string, keysAndValues ...interface{}) {
	if !l.Enabled() || l.level > Info {
		return
	}

	_, span := l.Tracer.Start(context.Background(), msg)
	defer span.End()

	kv, err := toKv(keysAndValues...)
	if err != nil {
		span.RecordError(err)
	}
	span.SetAttributes(l.values...)
	span.SetAttributes(kv...)
}

func (l Logger) Error(err error, msg string, keysAndValues ...interface{}) {
	if !l.Enabled() || l.level > Error {
		return
	}

	_, span := l.Tracer.Start(context.Background(), msg)
	defer span.End()
	span.RecordError(err)

	kv, err := toKv(keysAndValues...)
	if err != nil {
		span.RecordError(err)
	}
	span.SetAttributes(l.values...)
	span.SetAttributes(kv...)
}

func (l Logger) V(level int) logr.Logger {
	nl := l.clone()
	nl.level += level
	return nl
}

func (l Logger) WithValues(kvList ...interface{}) logr.Logger {
	nl := l.clone()
	kv, err := toKv(kvList...)
	if err != nil {
		l.Error(err, "WithValues")
	}

	nl.values = append(nl.values, kv...)
	return nl
}

func (l Logger) WithName(name string) logr.Logger {
	nl := l.clone()
	nl.Name = l.Name + "/" + name
	nl.Tracer = otel.Tracer(nl.Name)

	return nl
}

func (l Logger) clone() Logger {
	out := l
	l.values = copyKeyValues(l.values)
	return out
}

func copyKeyValues(in []label.KeyValue) []label.KeyValue {
	out := make([]label.KeyValue, len(in))
	copy(out, in)
	return out
}

func toKv(kvList ...interface{}) ([]label.KeyValue, error) {
	kv := []label.KeyValue{}

	for i := 0; i < len(kvList); i += 2 {
		k, ok := kvList[i].(string)
		if !ok {
			return kv, fmt.Errorf("key is not a string: %q", kvList[i])
		}
		var v interface{}
		if i+1 < len(kvList) {
			v = kvList[i+1]
		}
		kv = append(kv, label.Any(k, v))
	}

	return kv, nil
}
