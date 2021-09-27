//go:build telemetry
// +build telemetry

package telemetry

import (
	"context"
	"runtime"

	"github.com/denisbrodbeck/machineid"
	"github.com/rs/zerolog/log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel/sdk/resource"
)

func Start(ctx context.Context, honeycomb_api_key string, honeycomb_dataset string) {
	// if telemetry is turned on and honeycomb is configured, hook it up
	if honeycomb_api_key != "not_set" && honeycomb_dataset != "not_set" {
		exp, err := otlp.NewExporter(
			ctx,
			otlpgrpc.NewDriver(
				otlpgrpc.WithEndpoint("api.honeycomb.io:443"),
				otlpgrpc.WithHeaders(map[string]string{
					"x-honeycomb-team":    honeycomb_api_key,
					"x-honeycomb-dataset": honeycomb_dataset,
				}),
				otlpgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
			),
		)
		if err != nil {
			log.Fatal().Msgf("failed to initialize exporter: %v", err)
		}

		// Create a new tracer provider with a batch span processor and the otlp exporter.
		// Add a resource attribute service.name that identifies the service in the Honeycomb UI.
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exp),
			sdktrace.WithResource(resource.NewWithAttributes(semconv.ServiceNameKey.String("ExampleService"))),
		)

		// Handle this error in a sensible manner where possible
		defer func() { _ = tp.Shutdown(ctx) }()

		// Set the Tracer Provider and the W3C Trace Context propagator as globals
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(
			propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
		)
	}

	tracer := otel.Tracer("pct")

	uuid := attribute.Key("uuid")
	osKey := attribute.Key("osinfo/os")
	osArch := attribute.Key("osinfo/arch")

	var span trace.Span
	_, span = tracer.Start(ctx, "execution")
	defer span.End()

	// The Protected ID is hashed base on application name to prevent any
	// accidental leakage of a reversable ID.
	machineUUID, _ := machineid.ProtectedID("pdk")

	span.SetAttributes(uuid.String(machineUUID))
	span.SetAttributes(osKey.String(runtime.GOOS))
	span.SetAttributes(osArch.String(runtime.GOARCH))
}
