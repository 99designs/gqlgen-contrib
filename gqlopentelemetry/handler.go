package gqlopentelemetry

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const instrumentationName = "github.com/99designs/gqlgen-contrib/gqlopentelemetry"

// Handler is the extension to add OpenTelemetry Tracing
type Handler struct {
	ComplexityExtensionName string

	Tracer trace.Tracer
}

var _ interface {
	graphql.HandlerExtension
	graphql.OperationInterceptor
	graphql.FieldInterceptor
} = Handler{}

// ExtensionName is the name of the extension
func (a Handler) ExtensionName() string {
	return "OpenTelemetry"
}

// Validate validates the graphql schema
func (a Handler) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

// InterceptOperation runs tracing on every graphql operation
func (a Handler) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	ctx, span := a.getTracer().Start(ctx, operationName(ctx))
	defer span.End()

	oc := graphql.GetOperationContext(ctx)
	span.SetAttributes(
		attribute.String("request.query", oc.RawQuery),
	)
	complexityExtension := a.ComplexityExtensionName
	if complexityExtension == "" {
		complexityExtension = "ComplexityLimit" // from extention package
	}
	complexityStats, ok := oc.Stats.GetExtension(complexityExtension).(*extension.ComplexityStats)
	if !ok {
		// complexity extension is not used
		complexityStats = &extension.ComplexityStats{}
	}

	if complexityStats.ComplexityLimit > 0 {
		span.SetAttributes(
			attribute.Int("request.complexityLimit", complexityStats.ComplexityLimit),
			attribute.Int("request.operationComplexity", complexityStats.Complexity),
		)
	}

	for key, val := range oc.Variables {
		span.SetAttributes(
			attribute.String(fmt.Sprintf("request.variables.%s", key), fmt.Sprintf("%+v", val)),
		)
	}

	return next(ctx)
}

// InterceptField intercepts every field resolver
func (a Handler) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	fc := graphql.GetFieldContext(ctx)
	ctx, span := a.getTracer().Start(ctx, fc.Field.ObjectDefinition.Name+"/"+fc.Field.Name)
	defer span.End()

	span.SetAttributes(
		attribute.String("resolver.path", fc.Path().String()),
		attribute.String("resolver.object", fc.Field.ObjectDefinition.Name),
		attribute.String("resolver.field", fc.Field.Name),
		attribute.String("resolver.alias", fc.Field.Alias),
	)
	for _, arg := range fc.Field.Arguments {
		if arg.Value != nil {
			span.SetAttributes(
				attribute.String(fmt.Sprintf("resolver.args.%s", arg.Name), arg.Value.String()),
			)
		}
	}

	resp, err := next(ctx)

	errList := graphql.GetFieldErrors(ctx, fc)
	if len(errList) != 0 {
		span.RecordError(errList)
	}

	return resp, err
}

func operationName(ctx context.Context) string {
	opContext := graphql.GetOperationContext(ctx)
	requestName := "nameless-operation"
	if opContext.Doc != nil && len(opContext.Doc.Operations) != 0 {
		op := opContext.Doc.Operations[0]
		if op.Name != "" {
			requestName = op.Name
		}
	}

	return requestName
}

func (a Handler) getTracer() trace.Tracer {
	if a.Tracer == nil {
		a.Tracer = otel.Tracer(instrumentationName)
	}

	return a.Tracer
}
