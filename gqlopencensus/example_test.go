package gqlopencensus_test //nolint:testableexamples // allows the example only without any tests

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/99designs/gqlgen-contrib/gqlopencensus"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
)

var es graphql.ExecutableSchema

func Example() {
	// NOTE: requires setting of Exporter
	//   trace.RegisterExporter(exporter)

	logger := slog.Default()
	srv := handler.NewDefaultServer(es)
	srv.Use(gqlopencensus.Tracer{})

	httpServer := &http.Server{
		Addr:              ":8080",
		Handler:           srv,
		ReadHeaderTimeout: 5 * time.Second,
	}
	http.Handle("/query", httpServer.Handler)

	if err := httpServer.ListenAndServe(); err != nil {
		logger.Error(err.Error())
	}
}
