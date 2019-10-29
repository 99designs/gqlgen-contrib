package gqlopencensus_test

import (
	"github.com/99designs/gqlgen-contrib/gqlopencensus"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/99designs/gqlgen-contrib/gqlopencensus/internal/graph"
	"github.com/99designs/gqlgen/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpencensus_ResolverMiddleware_RequestMiddleware(t *testing.T) {

	gqlopencensus.Register()

	pe, err := prometheus.NewExporter(prometheus.Options{})

	if err != nil {
		log.Fatalf("Failed to create the Prometheus exporter: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/query", handler.GraphQL(
		graph.NewExecutableSchema(graph.Config{
			Resolvers: &graph.Resolver{},
		}),
		handler.RequestMiddleware(gqlopencensus.RequestMiddleware()),
		handler.ResolverMiddleware(gqlopencensus.ResolverMiddleware()),
	))

	for i := 0; i < 100; i++ {
		resp := doRequest(mux, http.MethodPost, "/query", `{"query":"{ todos { id text } }"}`)
		require.Equal(t, http.StatusOK, resp.Code)
	}

	resp := doRequest(pe, http.MethodGet, "/", "")
	require.Equal(t, http.StatusOK, resp.Code)

	gqlopencensus.UnRegister()

	body := resp.Body.String()

	assert.Contains(t, body, "graphql_request_duration_ms_bucket")
	assert.Contains(t, body, "graphql_resolver_duration_ms_bucket")
	assert.Contains(t, body, "graphql_request_started_total")
	assert.Contains(t, body, "graphql_request_completed_total")
	assert.Contains(t, body, "graphql_resolver_started_total")
	assert.Contains(t, body, "graphql_resolver_completed_total")
}

func doRequest(handler http.Handler, method string, target string, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	return w
}
