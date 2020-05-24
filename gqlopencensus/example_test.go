package gqlopencensus_test

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen-contrib/gqlopencensus"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
)

var es graphql.ExecutableSchema

func Example() {
	// NOTE: requires setting of Exporter
	//   trace.RegisterExporter(exporter)

	srv := handler.NewDefaultServer(es)
	srv.Use(gqlopencensus.Tracer{})

	http.Handle("/query", srv)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
