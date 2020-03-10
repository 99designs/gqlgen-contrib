package gqlopentracing

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

type OpenTracingTracer struct{}

var _ interface {
	graphql.HandlerExtension
	graphql.ResponseInterceptor
	graphql.FieldInterceptor
} = OpenTracingTracer{}

func (OpenTracingTracer) ExtensionName() string {
	return "Opentracing"
}

func (OpenTracingTracer) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func (OpenTracingTracer) InterceptField(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	fieldCtx := graphql.GetFieldContext(ctx)
	span, ctx := opentracing.StartSpanFromContext(ctx, fieldCtx.Path().String())
	defer span.Finish()
	ext.SpanKind.Set(span, "server")
	ext.Component.Set(span, "gqlgen")

	return next(ctx)
}

func (OpenTracingTracer) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	opCtx := graphql.GetOperationContext(ctx)
	opName := ""
	if opCtx.Operation != nil {
		opName = opCtx.Operation.Name
	}
	if opName == "" && opCtx.Operation != nil {
		//parent response case
		opName = string(opCtx.Operation.Operation)
	}
	if opName == "" {
		opName = opCtx.OperationName
	}
	span, ctx := opentracing.StartSpanFromContext(ctx, opName)
	defer span.Finish()
	ext.SpanKind.Set(span, "server")
	ext.Component.Set(span, "gqlgen")

	resp := next(ctx)
	if resp == nil {
		return nil
	}

	if err := resp.Errors.Error(); err != "" {
		ext.Error.Set(span, true)
		span.LogFields(log.String("error", err))
	}

	return resp
}
