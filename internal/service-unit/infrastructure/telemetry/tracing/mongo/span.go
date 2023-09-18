package mongo

import (
	"context"
	"fmt"

	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/telemetry/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func CreateSpanFactory(operation string) func(ctx context.Context, name, db, collection string) trace.Span {
	return func(ctx context.Context, name, db, collection string) trace.Span {
		tracer := tracing.GetTracer()
		_, span := tracer.Start(ctx, fmt.Sprintf("%s %s %s.%s", name, operation, db, collection))
		span.SetAttributes(
			attribute.String("db.system", "mongodb"),
			attribute.String("db.user", config.GetEnvs().MONGO_USER),
			attribute.String("server.address", name),
			attribute.String("server.port", config.GetEnvs().MONGO_PORT),
			attribute.String("network.transport", "IP.TCP"),
			attribute.String("db.name", db),
			attribute.String("db.operation", operation),
			attribute.String("db.mongodb.collection", collection),
		)
		return span
	}
}

type CreateReadSpanFactory func (ctx context.Context, name, db, collection string) trace.Span
type CreateWriteSpanFactory func (ctx context.Context, name, db, collection string) trace.Span

// Package-level variable holding the factory function
var CreateReadSpan CreateReadSpanFactory
var CreateWriteSpan CreateWriteSpanFactory

func init() {
    // Initialize the factory function
    CreateReadSpan = CreateSpanFactory("read")
    CreateWriteSpan = CreateSpanFactory("write")
}
