package prometheus

import (
	"context"
	"time"

	prometheusclient "github.com/prometheus/client_golang/prometheus"

	"github.com/99designs/gqlgen/graphql"
)

const (
	existStatusFailure = "failure"
	exitStatusSuccess  = "success"
)

var (
	requestStartedCounter    prometheusclient.Counter
	requestCompletedCounter  prometheusclient.Counter
	resolverStartedCounter   *prometheusclient.CounterVec
	resolverCompletedCounter *prometheusclient.CounterVec
	timeToResolveField       *prometheusclient.HistogramVec
	timeToHandleRequest      *prometheusclient.HistogramVec
)

type (
	Tracer struct{}
)

var _ interface {
	graphql.HandlerExtension
	graphql.OperationInterceptor
	graphql.ResponseInterceptor
	graphql.FieldInterceptor
} = Tracer{}

func Register() {
	RegisterOn(prometheusclient.DefaultRegisterer)
}

func RegisterOn(registerer prometheusclient.Registerer) {
	requestStartedCounter = prometheusclient.NewCounter(
		prometheusclient.CounterOpts{
			Name: "graphql_request_started_total",
			Help: "Total number of requests started on the graphql server.",
		},
	)

	requestCompletedCounter = prometheusclient.NewCounter(
		prometheusclient.CounterOpts{
			Name: "graphql_request_completed_total",
			Help: "Total number of requests completed on the graphql server.",
		},
	)

	resolverStartedCounter = prometheusclient.NewCounterVec(
		prometheusclient.CounterOpts{
			Name: "graphql_resolver_started_total",
			Help: "Total number of resolver started on the graphql server.",
		},
		[]string{"object", "field"},
	)

	resolverCompletedCounter = prometheusclient.NewCounterVec(
		prometheusclient.CounterOpts{
			Name: "graphql_resolver_completed_total",
			Help: "Total number of resolver completed on the graphql server.",
		},
		[]string{"object", "field"},
	)

	const start, factor, count = 1, 2, 11
	//nolint:promlinter // change the label name will break the existing metrics
	timeToResolveField = prometheusclient.NewHistogramVec(prometheusclient.HistogramOpts{
		Name:    "graphql_resolver_duration_ms",
		Help:    "The time taken to resolve a field by graphql server.",
		Buckets: prometheusclient.ExponentialBuckets(start, factor, count),
	}, []string{"exitStatus", "object", "field"})

	//nolint:promlinter // change the label name will break the existing metrics
	timeToHandleRequest = prometheusclient.NewHistogramVec(prometheusclient.HistogramOpts{
		Name:    "graphql_request_duration_ms",
		Help:    "The time taken to handle a request by graphql server.",
		Buckets: prometheusclient.ExponentialBuckets(start, factor, count),
	}, []string{"exitStatus"})

	registerer.MustRegister(
		requestStartedCounter,
		requestCompletedCounter,
		resolverStartedCounter,
		resolverCompletedCounter,
		timeToResolveField,
		timeToHandleRequest,
	)
}

func UnRegister() {
	UnRegisterFrom(prometheusclient.DefaultRegisterer)
}

func UnRegisterFrom(registerer prometheusclient.Registerer) {
	registerer.Unregister(requestStartedCounter)
	registerer.Unregister(requestCompletedCounter)
	registerer.Unregister(resolverStartedCounter)
	registerer.Unregister(resolverCompletedCounter)
	registerer.Unregister(timeToResolveField)
	registerer.Unregister(timeToHandleRequest)
}

func (a Tracer) ExtensionName() string {
	return "Prometheus"
}

func (a Tracer) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func (a Tracer) InterceptOperation(
	ctx context.Context,
	next graphql.OperationHandler,
) graphql.ResponseHandler {
	requestStartedCounter.Inc()
	return next(ctx)
}

func (a Tracer) InterceptResponse(
	ctx context.Context,
	next graphql.ResponseHandler,
) *graphql.Response {
	errList := graphql.GetErrors(ctx)

	var exitStatus string
	if len(errList) > 0 {
		exitStatus = existStatusFailure
	} else {
		exitStatus = exitStatusSuccess
	}

	oc := graphql.GetOperationContext(ctx)
	observerStart := oc.Stats.OperationStart

	timeToHandleRequest.With(prometheusclient.Labels{"exitStatus": exitStatus}).
		Observe(float64(time.Since(observerStart).Nanoseconds() / int64(time.Millisecond)))

	requestCompletedCounter.Inc()

	return next(ctx)
}

func (a Tracer) InterceptField(ctx context.Context, next graphql.Resolver) (any, error) {
	fc := graphql.GetFieldContext(ctx)

	resolverStartedCounter.WithLabelValues(fc.Object, fc.Field.Name).Inc()

	observerStart := time.Now()

	res, err := next(ctx)

	var exitStatus string
	if err != nil {
		exitStatus = existStatusFailure
	} else {
		exitStatus = exitStatusSuccess
	}

	timeToResolveField.WithLabelValues(exitStatus, fc.Object, fc.Field.Name).
		Observe(float64(time.Since(observerStart).Nanoseconds() / int64(time.Millisecond)))

	resolverCompletedCounter.WithLabelValues(fc.Object, fc.Field.Name).Inc()

	return res, err
}
