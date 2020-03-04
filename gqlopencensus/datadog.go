package gqlopencensus

// WithDataDog provides DataDog specific span attrs.
// see github.com/DataDog/opencensus-go-exporter-datadog
//func WithDataDog() Option {
//	return func(cfg *config) {
//		cfg.tracer = &datadogTracerImpl{cfg.tracer}
//	}
//}
//
//type datadogTracerImpl struct {
//	graphql.Tracer
//}
//
//func (dt *datadogTracerImpl) StartFieldResolverExecution(ctx context.Context, rc *graphql.ResolverContext) context.Context {
//	ctx = dt.Tracer.StartFieldResolverExecution(ctx, rc)
//	span := trace.FromContext(ctx)
//	if !span.IsRecordingEvents() {
//		return ctx
//	}
//	span.AddAttributes(
//		// key from gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext#ResourceName
//		trace.StringAttribute("resource.name", operationName(ctx)),
//	)
//
//	return ctx
//}

/////////////////// Refactored part //////////////////

// WithDataDog provides DataDog specific span attrs.
// see github.com/DataDog/opencensus-go-exporter-datadog
//func WithDataDog() {
//}
//
//type datadogTracerImpl struct{}
//
//func (datadogTracerImpl) InterceptField(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
//	span := trace.FromContext(ctx)
//	res, err = next(ctx)
//
//	if !span.IsRecordingEvents() {
//		return res, err
//	}
//	span.AddAttributes(
//		// key from gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext#ResourceName
//		trace.StringAttribute("resource.name", operationName(ctx)),
//	)
//
//	return res, err
//}
