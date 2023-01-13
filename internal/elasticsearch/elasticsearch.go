package elasticsearch

import (
	"context"

	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

const otelName = "github.com/nei7/gls/internal/elasticsearch"

func newOtelSpan(ctx context.Context, name string) trace.Span {
	_, span := otel.Tracer(otelName).Start(ctx, name)
	span.SetAttributes(semconv.DBSystemElasticsearch)

	return span
}
