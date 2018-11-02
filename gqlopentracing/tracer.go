package gqlopencensus

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

var _ graphql.Tracer = (tracerImpl)(0)

// New returns Tracer for OpenTracing.
// see https://opentracing.io/
func New() graphql.Tracer {
	return tracerImpl(0)
}

type tracerImpl int

func (tracerImpl) StartOperationParsing(ctx context.Context) context.Context {
	return ctx
}

func (tracerImpl) EndOperationParsing(ctx context.Context) {
}

func (tracerImpl) StartOperationValidation(ctx context.Context) context.Context {
	return ctx
}

func (tracerImpl) EndOperationValidation(ctx context.Context) {
}

func (tracerImpl) StartOperationExecution(ctx context.Context) context.Context {
	requestContext := graphql.GetRequestContext(ctx)
	span, ctx := opentracing.StartSpanFromContext(ctx, requestContext.RawQuery)
	ext.SpanKind.Set(span, "server")
	ext.Component.Set(span, "gqlgen")

	return ctx
}

func (tracerImpl) StartFieldExecution(ctx context.Context, field graphql.CollectedField) context.Context {
	span, ctx := opentracing.StartSpanFromContext(ctx, "unnamed")
	ext.SpanKind.Set(span, "server")
	ext.Component.Set(span, "gqlgen")

	return ctx
}

func (tracerImpl) StartFieldResolverExecution(ctx context.Context, rc *graphql.ResolverContext) context.Context {
	span := opentracing.SpanFromContext(ctx)
	span.SetOperationName(rc.Object + "_" + rc.Field.Name)
	span.SetTag("resolver.object", rc.Object)
	span.SetTag("resolver.field", rc.Field.Name)

	return ctx
}

func (tracerImpl) StartFieldChildExecution(ctx context.Context) context.Context {
	return ctx
}

func (tracerImpl) EndFieldExecution(ctx context.Context) {
	span := opentracing.SpanFromContext(ctx)
	rc := graphql.GetResolverContext(ctx)
	reqCtx := graphql.GetRequestContext(ctx)
	if reqCtx.HasError(rc) {
		ext.Error.Set(span, true)
		span.LogFields(
			log.String("event", "error"),
		)
		// TODO We need something like reqCtx.GetErrors(rc)
		// log.String("message", err.Error()),
		// log.String("error.kind", fmt.Sprintf("%T", err)),
	}

	defer span.Finish()
}

func (tracerImpl) EndOperationExecution(ctx context.Context) {
	span := opentracing.SpanFromContext(ctx)
	defer span.Finish()
}
