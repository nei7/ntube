package elasticsearch

import (
	"context"

	esv7 "github.com/elastic/go-elasticsearch/v7"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

const otelName = "github.com/nei7/ntube/internal/elasticsearch"

func newOtelSpan(ctx context.Context, name string) trace.Span {
	_, span := otel.Tracer(otelName).Start(ctx, name)
	span.SetAttributes(semconv.DBSystemElasticsearch)

	return span
}

func NewElasticSearch() (es *esv7.Client, err error) {
	es, err = esv7.NewDefaultClient()
	if err != nil {
		return
	}

	res, err := es.Info()
	if err != nil {
		return nil, err
	}

	defer func() {
		err = res.Body.Close()
	}()

	return
}
