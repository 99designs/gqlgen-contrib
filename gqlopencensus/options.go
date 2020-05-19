package gqlopencensus

import (
	"github.com/99designs/gqlgen/graphql"
	"go.opencensus.io/trace"
)

// Option for an opencensus tracer. At this moment, it is possible to configure span attributes retrieved from the GraphQL contexts.
type Option func(*config)

// FieldAttributer is a functor producing trace attributes from the GraphL field context
type FieldAttributer func(*graphql.FieldContext) []trace.Attribute

// OperationAttributer is a functor producing trace attributes from the GraphL operation context
type OperationAttributer func(*graphql.OperationContext) []trace.Attribute

type config struct {
	fieldAttributers     []FieldAttributer
	operationAttributers []OperationAttributer
}

func (c config) fieldAttributes(ctx *graphql.FieldContext) []trace.Attribute {
	attrs := make([]trace.Attribute, 0, 10)
	for _, apply := range c.fieldAttributers {
		attrs = append(attrs, apply(ctx)...)
	}
	return attrs
}

func (c config) operationAttributes(ctx *graphql.OperationContext) []trace.Attribute {
	attrs := make([]trace.Attribute, 0, 10)
	for _, apply := range c.operationAttributers {
		attrs = append(attrs, apply(ctx)...)
	}
	return attrs
}

func defaultTracer() *Tracer {
	return &Tracer{
		config: config{
			fieldAttributers: []FieldAttributer{func(fc *graphql.FieldContext) []trace.Attribute {
				return []trace.Attribute{
					trace.StringAttribute("server", "gqlgen"),
					trace.StringAttribute("field", fc.Field.Name),
				}
			},
			},
			operationAttributers: []OperationAttributer{func(oc *graphql.OperationContext) []trace.Attribute {
				return []trace.Attribute{
					trace.StringAttribute("server", "gqlgen"),
					trace.StringAttribute("operation", operationName(oc)),
				}
			},
			},
		},
	}
}

// WithFieldAttributes adds some extra attributes from the graphQL field context to the span
func WithFieldAttributes(attributers ...FieldAttributer) Option {
	return func(c *config) {
		c.fieldAttributers = append(c.fieldAttributers, attributers...)
	}
}

// WithOperationAttributes adds some extra attributes from the graphQL operation context to the span
func WithOperationAttributes(attributers ...OperationAttributer) Option {
	return func(c *config) {
		c.operationAttributers = append(c.operationAttributers, attributers...)
	}
}

// WithDataDog provides DataDog specific span attrs.
// see github.com/DataDog/opencensus-go-exporter-datadog
func WithDataDog() Option {
	return func(c *config) {
		c.operationAttributers = append(c.operationAttributers, func(oc *graphql.OperationContext) []trace.Attribute {
			return []trace.Attribute{
				trace.StringAttribute("resource.name", operationName(oc)),
			}
		})
	}
}

func operationName(ctx *graphql.OperationContext) (opName string) {
	if ctx.Operation != nil {
		opName = ctx.Operation.Name
	}
	if opName == "" && ctx.Operation != nil {
		//parent response case
		opName = string(ctx.Operation.Operation)
	}
	if opName == "" {
		opName = ctx.OperationName
	}
	return
}
