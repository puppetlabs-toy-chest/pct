//go:build !telemetry
// +build !telemetry

package telemetry

import (
	"context"
)

func Start(ctx context.Context, honeycomb_api_key string, honeycomb_dataset string, rootSpanName string) (context.Context, string, string) {
	// deliberately does nothing
	return ctx, "", ""
}

func EndSpan(span string) {
	// deliberately does nothing
}

func GetSpanFromContext(ctx context.Context) string {
	// deliberately does nothing
	return ""
}

func NewSpan(ctx context.Context, name string) (context.Context, string) {
	// deliberately does nothing
	return ctx, ""
}

func AddStringSpanAttribute(span string, key string, value string) {
	// deliberately does nothing
}

func ShutDown(ctx context.Context, provider string, span string) {
	// deliberately does nothing
}
