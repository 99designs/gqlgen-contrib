package prometheus_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/99designs/gqlgen-contrib/prometheus"
	"github.com/99designs/gqlgen-contrib/prometheus/internal/graph"
)

func TestPrometheus_ResolverMiddleware_RequestMiddleware(t *testing.T) {

	prometheus.Register()

	mux := http.NewServeMux()
	gqlHandler := handler.NewDefaultServer(
		graph.NewExecutableSchema(graph.Config{
			Resolvers: &graph.Resolver{},
		}),
	)
	gqlHandler.Use(&prometheus.Metrics{})
	mux.Handle("/query", gqlHandler)

	for i := 0; i < 100; i++ {
		resp := doRequest(mux, http.MethodPost, "/query", `{"query":"{ todos { id text } }"}`)
		require.Equal(t, http.StatusOK, resp.Code)
	}

	resp := doRequest(promhttp.Handler(), http.MethodGet, "/", "")
	require.Equal(t, http.StatusOK, resp.Code)

	prometheus.UnRegister()

	body := resp.Body.String()

	assert.Contains(t, body, "graphql_request_duration_ms_bucket")
	assert.Contains(t, body, "graphql_request_duration_ms_count")
	assert.Contains(t, body, "graphql_request_duration_ms_sum")
	assert.Contains(t, body, "graphql_resolver_duration_ms_bucket")
	assert.Contains(t, body, "graphql_resolver_duration_ms_count")
	assert.Contains(t, body, "graphql_resolver_duration_ms_sum")
}

func doRequest(handler http.Handler, method string, target string, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	return w
}
