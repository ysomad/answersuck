package app

import (
	"time"

	"github.com/exaring/otelpgx"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/ysomad/answersuck/internal/config"
)

type openTelemetry struct {
	appTracer      trace.Tracer
	pgxTracer      *otelpgx.Tracer
	tracerProvider *sdktrace.TracerProvider
	meterProvider  *metric.MeterProvider
}

func newOpenTelemetry(conf *config.Config) (ot openTelemetry, err error) {
	res := newResource(conf.App)

	jaegerExp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(conf.Jaeger.Endpoint)))
	if err != nil {
		return openTelemetry{}, err
	}

	ot.tracerProvider = newTracerProvider(res, jaegerExp)

	prometheusExp, err := prometheus.New()
	if err != nil {
		return openTelemetry{}, err
	}

	ot.meterProvider, err = newMeterProvider(res, prometheusExp)
	if err != nil {
		return openTelemetry{}, err
	}

	ot.appTracer = otel.GetTracerProvider().Tracer(conf.App.Name)
	ot.pgxTracer = otelpgx.NewTracer(otelpgx.WithTrimSQLInSpanName())

	return ot, nil
}

func newResource(conf config.App) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(conf.Name),
		semconv.ServiceVersionKey.String(conf.Ver),
		attribute.String("environment", conf.Environment),
	)
}

// newTracerProvider returns and registers configured tracer provider.
//
// docs: https://opentelemetry.io/docs/instrumentation/go/exporting_data/
func newTracerProvider(r *resource.Resource, exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)

	otel.SetTracerProvider(tp)

	return tp
}

// newMeterProvider sets global meter provider and returns it.
func newMeterProvider(res *resource.Resource, exp *prometheus.Exporter) (*metric.MeterProvider, error) {
	mp := metric.NewMeterProvider(
		metric.WithReader(exp),
		metric.WithResource(res),
	)

	global.SetMeterProvider(mp)

	if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second)); err != nil {
		return nil, err
	}

	return mp, nil
}
