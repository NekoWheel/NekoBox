// Copyright 2025 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// config is a group of options for this instrumentation.
type config struct {
	TracerProvider trace.TracerProvider
	Propagators    propagation.TextMapPropagator
}

// Option applies an option value for a config.
type Option interface {
	apply(*config)
}

// newConfig returns a config configured with all the passed Options.
func newConfig(opts []Option) *config {
	c := &config{
		Propagators:    otel.GetTextMapPropagator(),
		TracerProvider: otel.GetTracerProvider(),
	}
	for _, o := range opts {
		o.apply(c)
	}
	return c
}

type propagatorsOption struct{ p propagation.TextMapPropagator }

func (o propagatorsOption) apply(c *config) {
	if o.p != nil {
		c.Propagators = o.p
	}
}

// WithPropagators returns an Option to use the Propagators when extracting
// and injecting trace context from HTTP requests.
func WithPropagators(p propagation.TextMapPropagator) Option {
	return propagatorsOption{p: p}
}

type tracerProviderOption struct{ tp trace.TracerProvider }

func (o tracerProviderOption) apply(c *config) {
	if o.tp != nil {
		c.TracerProvider = o.tp
	}
}

// WithTracerProvider returns an Option to use the TracerProvider when
// creating a Tracer.
func WithTracerProvider(tp trace.TracerProvider) Option {
	return tracerProviderOption{tp: tp}
}
