package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func StartSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	tr := otel.Tracer("")
	return tr.Start(ctx, spanName)
}
