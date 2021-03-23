package gqlopentracing

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

type OpenTracingTracer struct{}

var _ interface {
	graphql.HandlerExtension
	graphql.FieldInterceptor
} = OpenTracingTracer{}

func (OpenTracingTracer) ExtensionName() string {
	return "Opentracing"
}

func (OpenTracingTracer) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func (OpenTracingTracer) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	opCtx := graphql.GetOperationContext(ctx)

	opName := opCtx.OperationName
	if opName == "" {
		opName = "unnamed_operation"
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, opName)
	ext.SpanKind.Set(span, "server")
	ext.Component.Set(span, "gqlgen")
	defer span.Finish()

	return next(ctx)
}

func (OpenTracingTracer) InterceptField(ctx context.Context, next graphql.Resolver) ( interface{},  error) {
	fieldCtx := graphql.GetFieldContext(ctx)

	// dont start spans for a plain field (not a resolver method)
	if !fieldCtx.IsMethod {
		return next(ctx)
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, fieldCtx.Object + "." + fieldCtx.Field.Name)
	defer span.Finish()
	span.SetTag("resolver.object", fieldCtx.Object)
	span.SetTag("resolver.method", fieldCtx.Field.Name)
	span.SetTag("resolver.path", fieldCtx.Path())

	res, err := next(ctx)

	errList := graphql.GetFieldErrors(ctx, fieldCtx)
	if len(errList) != 0 {
		ext.Error.Set(span, true)
		span.LogFields(
			log.String("event", "error"),
		)

		for idx, err := range errList {
			span.LogFields(log.String(fmt.Sprintf("error.%d.message", idx), err.Error()))
		}
	}

	return res, err
}
