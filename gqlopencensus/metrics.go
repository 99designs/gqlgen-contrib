package gqlopencensus

import (
	"context"
	"log"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

const (
	existStatusFailure = "failure"
	exitStatusSuccess  = "success"
)

var _ graphql.Tracer = (tracerImpl)(0)

var (
	requestStartedCounter    = stats.Int64("graphql_request_started_total", "Total number of requests started on the graphql server.", stats.UnitDimensionless)
	requestCompletedCounter  = stats.Int64("graphql_request_completed_total", "Total number of requests started on the graphql server.", stats.UnitDimensionless)
	resolverStartedCounter   = stats.Int64("graphql_resolver_started_total", "Total number of resolver started on the graphql server.", stats.UnitDimensionless)
	resolverCompletedCounter = stats.Int64("graphql_resolver_completed_total", "Total number of resolver completed on the graphql server.", stats.UnitDimensionless)
	// TODO: remove ms
	timeToResolveField = stats.Int64("graphql_resolver_duration_ms", "The time taken to resolve a field by graphql server.", stats.UnitMilliseconds)
	// TODO: remove ms
	timeToHandleRequest = stats.Int64("graphql_request_duration_ms", "The time taken to handle a request by graphql server.", stats.UnitMilliseconds)

	keyObject, _     = tag.NewKey("object")
	keyField, _      = tag.NewKey("field")
	keyExitStatus, _ = tag.NewKey("exitStatus")

	requestStartedCounterView = &view.View{
		Name:        "graphql_request_started_total",
		Description: "Total number of requests started on the graphql server.",
		Measure:     requestStartedCounter,
		Aggregation: view.Count(),
	}
	requestCompletedCounterView = &view.View{
		Name:        "graphql_request_completed_total",
		Description: "Total number of requests completed on the graphql server.",
		Measure:     requestCompletedCounter,
		Aggregation: view.Count(),
	}
	resolverStartedCounterView = &view.View{
		Name:        "graphql_resolver_started_total",
		Description: "Total number of resolver started on the graphql server.",
		Measure:     resolverStartedCounter,
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{keyObject, keyField},
	}
	resolverCompletedCounterView = &view.View{
		Name:        "graphql_resolver_completed_total",
		Description: "Total number of resolver completed on the graphql server.",
		Measure:     resolverCompletedCounter,
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{keyObject, keyField},
	}
	timeToResolveFieldView = &view.View{
		Name:        "graphql_resolver_duration_ms",
		Description: "The time taken to resolve a field by graphql server.",
		Measure:     timeToResolveField,
		// TODO: adjust buckets
		Aggregation: view.Distribution(1<<32, 2<<32, 3<<32),
		TagKeys:     []tag.Key{keyExitStatus},
	}
	timeToHandleRequestView = &view.View{
		Name:        "graphql_request_duration_ms",
		Description: "The time taken to handle a request by graphql server.",
		Measure:     timeToHandleRequest,
		// TODO: adjust buckets
		Aggregation: view.Distribution(1<<32, 2<<32, 3<<32),
		TagKeys:     []tag.Key{keyExitStatus},
	}
)

func Register() {
	RegisterOn()
}

func RegisterOn() {
	if err := view.Register(requestStartedCounterView); err != nil {
		log.Fatalf("Failed to register view: %v", err)
	}

	if err := view.Register(requestCompletedCounterView); err != nil {
		log.Fatalf("Failed to register view: %v", err)
	}

	if err := view.Register(resolverStartedCounterView); err != nil {
		log.Fatalf("Failed to register view: %v", err)
	}

	if err := view.Register(resolverCompletedCounterView); err != nil {
		log.Fatalf("Failed to register view: %v", err)
	}

	if err := view.Register(timeToResolveFieldView); err != nil {
		log.Fatalf("Failed to register view: %v", err)
	}

	if err := view.Register(timeToHandleRequestView); err != nil {
		log.Fatalf("Failed to register view: %v", err)
	}
}

func UnRegister() {
	UnRegisterFrom()
}

func UnRegisterFrom() {
	view.Unregister(requestStartedCounterView)
	view.Unregister(requestCompletedCounterView)
	view.Unregister(resolverStartedCounterView)
	view.Unregister(resolverCompletedCounterView)
	view.Unregister(timeToResolveFieldView)
	view.Unregister(timeToHandleRequestView)
}

func ResolverMiddleware() graphql.FieldMiddleware {
	return func(ctx context.Context, next graphql.Resolver) (interface{}, error) {
		rctx := graphql.GetResolverContext(ctx)

		ctx, err := tag.New(context.Background(),
			tag.Insert(keyObject, rctx.Object),
			tag.Insert(keyField, rctx.Field.Name),
		)

		if err != nil {
			return nil, err
		}

		stats.Record(ctx, resolverStartedCounter.M(1))
		// resolverStartedCounter.WithLabelValues(rctx.Object, rctx.Field.Name).Inc()

		observerStart := time.Now()

		res, err := next(ctx)

		var exitStatus string
		if err != nil {
			exitStatus = existStatusFailure
		} else {
			exitStatus = exitStatusSuccess
		}

		ctx, err = tag.New(context.Background(),
			tag.Insert(keyExitStatus, exitStatus),
		)

		if err != nil {
			return nil, err
		}

		stats.Record(ctx, timeToResolveField.M(time.Since(observerStart).Nanoseconds()/int64(time.Millisecond)))
		//timeToResolveField.With(prometheusclient.Labels{"exitStatus": exitStatus}).
		//	Observe(float64(time.Since(observerStart).Nanoseconds() / int64(time.Millisecond)))

		stats.Record(ctx, resolverCompletedCounter.M(1))
		// resolverCompletedCounter.WithLabelValues(rctx.Object, rctx.Field.Name).Inc()

		return res, err
	}
}

func RequestMiddleware() graphql.RequestMiddleware {
	return func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
		// requestStartedCounter.Inc()
		stats.Record(ctx, requestStartedCounter.M(1))

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

		ctx, _ = tag.New(context.Background(),
			tag.Insert(keyExitStatus, exitStatus),
		)

		stats.Record(ctx, timeToHandleRequest.M(time.Since(observerStart).Nanoseconds()/int64(time.Millisecond)))
		//timeToHandleRequest.With(prometheusclient.Labels{"exitStatus": exitStatus}).
		//	Observe(float64(time.Since(observerStart).Nanoseconds() / int64(time.Millisecond)))

		stats.Record(ctx, requestCompletedCounter.M(1))
		// requestCompletedCounter.Inc()

		return res
	}
}
