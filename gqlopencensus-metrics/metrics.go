// Package metrics collects opencensus metrics
// for a GraphQL server.
package metrics

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// Register views
func Register() error {
	return view.Register(GQLViews...)
}

// Unregister views
func Unregister() {
	view.Unregister(GQLViews...)
}

var (
	// GQLViews contains all opencensus stats views declared by the GraphQL stats collector
	GQLViews = []*view.View{
		OperationCountView,
		FieldCountView,
		OperationErrorsView,
		OperationLatencyView,
		FieldLatencyView,
		OperationParsingView,
	}

	// measurements

	ServerRequestCount = stats.Int64(
		"gql/server/request_count",
		"Number of GraphQL requests started",
		stats.UnitDimensionless)

	ServerFieldCount = stats.Int64(
		"gql/server/field_count",
		"Number of GraphQL field resolutions, per field and query path",
		stats.UnitDimensionless)

	ServerErrorCount = stats.Int64(
		"gql/server/request_count",
		"Number of GraphQL requests started",
		stats.UnitDimensionless)

	ServerLatency = stats.Float64(
		"gql/server/latency",
		"Execution latency",
		stats.UnitMilliseconds)

	ServerFieldLatency = stats.Float64(
		"gql/server/field_latency",
		"Single field execution latency",
		stats.UnitMilliseconds)

	ServerParsing = stats.Float64(
		"gql/server/parsing_validation",
		"Parsing & validation latency",
		stats.UnitMilliseconds)

	// views

	OperationCountView = &view.View{
		Name:        "gql/server/operation_count",
		Description: "Count of GraphQL requests started by operation",
		Measure:     ServerRequestCount,
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{TagHost, TagOperation},
	}

	FieldCountView = &view.View{
		Name:        "gql/server/field_count",
		Description: "Count of GraphQL fields requests by field and by query path",
		Measure:     ServerFieldCount,
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{TagHost, TagField, TagPath},
	}

	OperationErrorsView = &view.View{
		Name:        "gql/server/operation_error_count",
		Description: "Count of GraphQL requests returning an error by operation",
		Measure:     ServerErrorCount,
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{TagHost, TagOperation},
	}

	OperationLatencyView = &view.View{
		Name:        "gql/server/latency",
		Description: "Execution time distribution of GraphQL requests by operation, excluding parsing and validation",
		Measure:     ServerLatency,
		Aggregation: DefaultLatencyDistribution,
		TagKeys:     []tag.Key{TagHost, TagOperation},
	}

	FieldLatencyView = &view.View{
		Name:        "gql/server/field_latency",
		Description: "Execution time distribution of GraphQL requests by operation, excluding parsing and validation",
		Measure:     ServerFieldLatency,
		Aggregation: DefaultLatencyDistribution,
		TagKeys:     []tag.Key{TagHost, TagField, TagPath},
	}

	OperationParsingView = &view.View{
		Name:        "gql/server/parsing_validation",
		Description: "Parsing  and validation time distribution of GraphQL requests by operation",
		Measure:     ServerParsing,
		Aggregation: DefaultLatencyDistribution,
		TagKeys:     []tag.Key{TagHost, TagOperation},
	}

	// TagHost is the name of the graphQL server
	TagHost = tag.MustNewKey("gql.host")

	// TagOperation is the query operation name
	TagOperation = tag.MustNewKey("gql.operation")

	// TagField is an individual GraphQL field requested
	TagField = tag.MustNewKey("gql.field")

	// TagPath is an individual GraphQL path to a field requested
	TagPath = tag.MustNewKey("gql.path")

	DefaultSizeDistribution    = view.Distribution(1024, 2048, 4096, 16384, 65536, 262144, 1048576, 4194304, 16777216, 67108864, 268435456, 1073741824, 4294967296)
	DefaultLatencyDistribution = view.Distribution(1, 2, 3, 4, 5, 6, 8, 10, 13, 16, 20, 25, 30, 40, 50, 65, 80, 100, 130, 160, 200, 250, 300, 400, 500, 650, 800, 1000, 2000, 5000, 10000, 20000, 50000, 100000)
)
