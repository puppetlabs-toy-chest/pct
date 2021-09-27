//go:build !telemetry
// +build !telemetry

package telemetry

import (
	"context"
)

func Start(ctx context.Context, honeycomb_api_key string, honeycomb_dataset string) {
	// deliberately does nothing
}
