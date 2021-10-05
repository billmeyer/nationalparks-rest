package pkg

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"os"
)

// InitOTEL Initializes a Jaeger exporter, and configures the corresponding trace and metric providers.
func InitOTEL(serviceName string, environmentName string) func(context.Context) error {
	//ctx := context.Background()

	os.Setenv("OTEL_LOG_LEVEL", "debug")
	var accessToken = os.Getenv("SPLUNK_ACCESS_TOKEN")
	if accessToken == "" {
		log.Fatalf("SPLUNK_ACCESS_TOKEN environment variable not set")
	}

	var realm = os.Getenv("SPLUNK_REALM")
	if realm == "" {
		log.Fatalf("SPLUNK_REALM environment variable not set")
	}

	// Initialize Jaeger Exporter
	var collectorEndpointOption []jaeger.CollectorEndpointOption
	//collectorEndpointOption = append(collectorEndpointOption, jaeger.WithEndpoint("http://octodev02:14268/api/traces"))
	//collectorEndpointOption = append(collectorEndpointOption, jaeger.WithEndpoint("http://localhost:14268/api/traces"))
	collectorEndpointOption = append(collectorEndpointOption, jaeger.WithEndpoint(fmt.Sprintf("https://ingest.%s.signalfx.com/v2/trace", realm)))
	collectorEndpointOption = append(collectorEndpointOption, jaeger.WithUsername("auth"), jaeger.WithPassword(accessToken))

	var endpointOption = jaeger.WithCollectorEndpoint(collectorEndpointOption...)
	jaegerExporter, err := jaeger.New(endpointOption)
	handleErr(err, "failed to create jaegerExporter")

	// Resources that will be attached to our Trace Provider
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", environmentName),
			attribute.String("library.language", "go"),
		),
	)
	handleErr(err, "failed to create resource")

	bsp := sdktrace.NewBatchSpanProcessor(jaegerExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithSyncer(jaegerExporter),
	)

	// set global propagator to tracecontext (the default is no-op).
	//otel.SetTextMapPropagator(propagation.TraceContext{})

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tracerProvider)

	return jaegerExporter.Shutdown
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}
