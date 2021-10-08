package pkg

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"net/http"

	"os"
)

// InitOTEL Initializes a Jaeger exporter, and configures the corresponding trace and metric providers.
func InitOTEL(serviceName string, environmentName string) func(context.Context) error {
	os.Setenv("OTEL_LOG_LEVEL", "debug")

	const useSplunk = true
	var collectorUrl string

	// Initialize Jaeger Exporter
	var collectorEndpointOptions []jaeger.CollectorEndpointOption

	if useSplunk == true {
		// Send Traces directly to Splunk Observability Cloud for ingestion.
		// To do this, we need to use the correct ingest URL (see below).
		// Additionally, we need to pass our Splunk Realm and Access Token as well.  These are read from the
		// SPLUNK_REALM and SPLUNK_ACCESS_TOKEN environment variables respectively.

		var accessToken = os.Getenv("SPLUNK_ACCESS_TOKEN")
		if accessToken == "" {
			log.Fatalf("SPLUNK_ACCESS_TOKEN environment variable not set")
		}

		var realm = os.Getenv("SPLUNK_REALM")
		if realm == "" {
			log.Fatalf("SPLUNK_REALM environment variable not set")
		}

		collectorUrl = fmt.Sprintf("https://ingest.%s.signalfx.com/v2/trace", realm)
		collectorEndpointOptions = append(collectorEndpointOptions, jaeger.WithUsername("auth"), jaeger.WithPassword(accessToken))
	} else {
		// Send Traces to a local OpenTelemetry Collector instance.  For this, we just need to specify the URL.
		collectorUrl = "http://octodev02:14268/api/traces"
	}

	collectorEndpointOptions = append(collectorEndpointOptions, jaeger.WithEndpoint(collectorUrl))

	var endpointOption = jaeger.WithCollectorEndpoint(collectorEndpointOptions...)
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

	// Set the Tracer Provider and the W3C Trace Context propagator as globals
	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			//jaegerPropagator.Jaeger{},
			//propagation.Baggage{},
			propagation.TraceContext{}),
	)

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tracerProvider)

	return jaegerExporter.Shutdown
}

// AddTraceParentToResponse Adds the required HTTP Headers to the HTTP Response such that trace context is propagated
// between services.  This is only necessary for the parent-most Span in a Trace.  All child spans will utilize
// the same HTTP Response so it's not necessary to call for every child span.
func AddTraceParentToResponse(span trace.Span, w http.ResponseWriter) {
	traceParent, _ := formatAsTraceParent(span.SpanContext())
	fmt.Printf("traceparent: %s\n", traceParent)
	w.Header().Add("Access-Control-Expose-Headers", "Server-Timing")
	w.Header().Add("Server-Timing", traceParent)
}

// Converts the Trace ID and Span ID in the supplied context into a `traceparent` entry in the standard W3C format.
// See https://www.w3.org/TR/trace-context/#version-format for specific formatting.
func formatAsTraceParent(ctx trace.SpanContext) (string, bool) {
	traceID := ctx.TraceID()
	spanID := ctx.SpanID()

	answer := fmt.Sprintf("traceparent;desc=\"00-%s-%s-01\"", traceID, spanID)
	return answer, true
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}
