package otel

import (
	oteltrace "go.opentelemetry.io/otel/trace"
)

type config struct {
	tracerProvider oteltrace.TracerProvider
	serviceName    string
}

// Option is used to configure the client.
type Option func(*config)

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
// If none is specified, the global provider is used.
func WithTracerProvider(provider oteltrace.TracerProvider) Option {
	return func(cfg *config) {
		cfg.tracerProvider = provider
	}
}

// WithServiceName sets the service name.
func WithServiceName(serviceName string) Option {
	return func(cfg *config) {
		cfg.serviceName = serviceName
	}
}
