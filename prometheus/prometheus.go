package prometheus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
	prometheusclient "github.com/prometheus/client_golang/prometheus"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const (
	existStatusFailure = "failure"
	exitStatusSuccess  = "success"
	errCodeExt         = "error_code"
)

var (
	timeToResolveField  *prometheusclient.HistogramVec
	timeToHandleRequest *prometheusclient.HistogramVec
)

func Register() {
	RegisterOn(prometheusclient.DefaultRegisterer)
}

func RegisterOn(registerer prometheusclient.Registerer) {
	timeToResolveField = prometheusclient.NewHistogramVec(prometheusclient.HistogramOpts{
		Name: "graphql_resolver_duration_ms",
		Help: "The time taken to resolve a field by graphql server.",
	}, []string{"err_code", "exit_status", "object", "field"})

	timeToHandleRequest = prometheusclient.NewHistogramVec(prometheusclient.HistogramOpts{
		Name: "graphql_request_duration_ms",
		Help: "The time taken to handle a request by graphql server.",
	}, []string{"exit_status", "operation"})

	registerer.MustRegister(
		timeToResolveField,
		timeToHandleRequest,
	)
}

func UnRegister() {
	UnRegisterFrom(prometheusclient.DefaultRegisterer)
}

func UnRegisterFrom(registerer prometheusclient.Registerer) {
	registerer.Unregister(timeToResolveField)
	registerer.Unregister(timeToHandleRequest)
}

type Metrics struct{}

var _ interface {
	graphql.HandlerExtension
	graphql.ResponseInterceptor
	graphql.FieldInterceptor
} = Metrics{}

func (Metrics) ExtensionName() string {
	return "PrometheusMetrics"
}

func (Metrics) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func (Metrics) InterceptField(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	fieldCtx := graphql.GetFieldContext(ctx)

	defer func(start time.Time) {
		var exitStatus string
		if err != nil {
			exitStatus = existStatusFailure
		} else {
			exitStatus = exitStatusSuccess
		}

		var errCode string
		var gqlErr *gqlerror.Error
		if err != nil && errors.As(err, &gqlErr) {
			if errCodeVal, ok := gqlErr.Extensions[errCodeExt]; ok {
				errCode = fmt.Sprint(errCodeVal)
			}
		}

		timeToResolveField.WithLabelValues(errCode, exitStatus, fieldCtx.Object, fieldCtx.Field.Name).
			Observe(float64(time.Since(start).Nanoseconds() / int64(time.Millisecond)))
	}(time.Now())

	res, err = next(ctx)
	return res, err
}

func (Metrics) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) (res *graphql.Response) {
	opCtx := graphql.GetOperationContext(ctx)

	defer func(start time.Time) {
		if res == nil {
			return
		}

		var exitStatus string
		if res.Errors.Error() != "" {
			exitStatus = existStatusFailure
		} else {
			exitStatus = exitStatusSuccess
		}

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

		timeToHandleRequest.WithLabelValues(exitStatus, opName).
			Observe(float64(time.Since(start).Nanoseconds() / int64(time.Millisecond)))

	}(time.Now())

	res = next(ctx)
	return res
}
