package gqlopencensus

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"go.opencensus.io/trace"
)

type (
	Tracer struct {
		ComplexityExtensionName string
		DataDog                 bool
	}
)

var _ interface {
	graphql.HandlerExtension
	graphql.OperationInterceptor
	graphql.FieldInterceptor
} = Tracer{}

func (a Tracer) ExtensionName() string {
	return "OpenCensus"
}

func (a Tracer) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func (a Tracer) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	ctx, span := trace.StartSpan(ctx, operationName(ctx))
	if !span.IsRecordingEvents() {
		return next(ctx)
	}
	defer span.End()

	oc := graphql.GetOperationContext(ctx)
	span.AddAttributes(
		trace.StringAttribute("request.query", oc.RawQuery),
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
		span.AddAttributes(
			trace.Int64Attribute("request.complexityLimit", int64(complexityStats.ComplexityLimit)),
			trace.Int64Attribute("request.operationComplexity", int64(complexityStats.Complexity)),
		)
	}

	for key, val := range oc.Variables {
		span.AddAttributes(
			trace.StringAttribute(fmt.Sprintf("request.variables.%s", key), fmt.Sprintf("%+v", val)),
		)
	}

	return next(ctx)
}

func (a Tracer) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	fc := graphql.GetFieldContext(ctx)
	ctx, span := trace.StartSpan(ctx, fc.Field.ObjectDefinition.Name+"/"+fc.Field.Name)
	defer span.End()
	if !span.IsRecordingEvents() {
		return next(ctx)
	}
	span.AddAttributes(
		trace.StringAttribute("resolver.path", fc.Path().String()),
		trace.StringAttribute("resolver.object", fc.Field.ObjectDefinition.Name),
		trace.StringAttribute("resolver.field", fc.Field.Name),
		trace.StringAttribute("resolver.alias", fc.Field.Alias),
	)
	if a.DataDog {
		span.AddAttributes(
			// key from gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext#ResourceName
			trace.StringAttribute("resource.name", operationName(ctx)),
		)
	}
	for _, arg := range fc.Field.Arguments {
		if arg.Value != nil {
			span.AddAttributes(
				trace.StringAttribute(fmt.Sprintf("resolver.args.%s", arg.Name), arg.Value.String()),
			)
		}
	}

	resp, err := next(ctx)

	errList := graphql.GetFieldErrors(ctx, fc)
	if len(errList) != 0 {
		span.SetStatus(trace.Status{
			Code:    2, // UNKNOWN, HTTP Mapping: 500 Internal Server Error
			Message: errList.Error(),
		})
		span.AddAttributes(
			trace.BoolAttribute("resolver.hasError", true),
			trace.Int64Attribute("resolver.errorCount", int64(len(errList))),
		)
		for idx, err := range errList {
			span.AddAttributes(
				trace.StringAttribute(fmt.Sprintf("resolver.error.%d.message", idx), err.Error()),
				trace.StringAttribute(fmt.Sprintf("resolver.error.%d.kind", idx), fmt.Sprintf("%T", err)),
			)
		}
	}

	return resp, err
}

func operationName(ctx context.Context) string {
	requestContext := graphql.GetRequestContext(ctx)
	requestName := "nameless-operation"
	if requestContext.Doc != nil && len(requestContext.Doc.Operations) != 0 {
		op := requestContext.Doc.Operations[0]
		if op.Name != "" {
			requestName = op.Name
		}
	}

	return requestName
}
