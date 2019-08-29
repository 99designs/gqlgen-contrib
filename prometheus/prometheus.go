package prometheus

import (
	"context"
	"time"

	"github.com/99designs/gqlgen/graphql"
	prometheusclient "github.com/prometheus/client_golang/prometheus"
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

func Register() {
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

	timeToResolveField = prometheusclient.NewHistogramVec(prometheusclient.HistogramOpts{
		Name:    "graphql_resolver_duration_ms",
		Help:    "The time taken to resolve a field by graphql server.",
		Buckets: prometheusclient.ExponentialBuckets(1, 2, 11),
	}, []string{"exitStatus", "object", "field"})

	timeToHandleRequest = prometheusclient.NewHistogramVec(prometheusclient.HistogramOpts{
		Name:    "graphql_request_duration_ms",
		Help:    "The time taken to handle a request by graphql server.",
		Buckets: prometheusclient.ExponentialBuckets(1, 2, 11),
	}, []string{"exitStatus"})

	prometheusclient.MustRegister(
		requestStartedCounter,
		requestCompletedCounter,
		resolverStartedCounter,
		resolverCompletedCounter,
		timeToResolveField,
		timeToHandleRequest,
	)
}

func UnRegister() {
	prometheusclient.Unregister(requestStartedCounter)
	prometheusclient.Unregister(requestCompletedCounter)
	prometheusclient.Unregister(resolverStartedCounter)
	prometheusclient.Unregister(resolverCompletedCounter)
	prometheusclient.Unregister(timeToResolveField)
	prometheusclient.Unregister(timeToHandleRequest)
}

func ResolverMiddleware() graphql.FieldMiddleware {
	return func(ctx context.Context, next graphql.Resolver) (interface{}, error) {
		rctx := graphql.GetResolverContext(ctx)

		resolverStartedCounter.WithLabelValues(rctx.Object, rctx.Field.Name).Inc()

		observerStart := time.Now()

		res, err := next(ctx)

		var exitStatus string
		if err != nil {
			exitStatus = existStatusFailure
		} else {
			exitStatus = exitStatusSuccess
		}

		timeToResolveField.WithLabelValues(exitStatus, rctx.Object, rctx.Field.Name).
			Observe(float64(time.Since(observerStart).Nanoseconds() / int64(time.Millisecond)))

		resolverCompletedCounter.WithLabelValues(rctx.Object, rctx.Field.Name).Inc()

		return res, err
	}
}

func RequestMiddleware() graphql.RequestMiddleware {
	return func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
		requestStartedCounter.Inc()

		observerStart := time.Now()

		res := next(ctx)

		rctx := graphql.GetResolverContext(ctx)
		reqCtx := graphql.GetRequestContext(ctx)
		errList := reqCtx.GetErrors(rctx)

		var exitStatus string
		if len(errList) > 0 {
			exitStatus = existStatusFailure
		} else {
			exitStatus = exitStatusSuccess
		}

		timeToHandleRequest.With(prometheusclient.Labels{"exitStatus": exitStatus}).
			Observe(float64(time.Since(observerStart).Nanoseconds() / int64(time.Millisecond)))

		requestCompletedCounter.Inc()

		return res
	}
}
