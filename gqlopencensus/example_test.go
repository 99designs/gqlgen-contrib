package gqlopencensus_test

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen-contrib/gqlopencensus"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
)

var es graphql.ExecutableSchema

func Example() {
	// NOTE: requires setting of Exporter
	//   trace.RegisterExporter(exporter)

	handler := handler.GraphQL(
		es,
		handler.Tracer(gqlopencensus.New()),
	)
	http.Handle("/query", handler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
