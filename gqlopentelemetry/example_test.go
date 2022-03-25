package gqlopentelemetry_test

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen-contrib/gqlopentelemetry"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
)

var es graphql.ExecutableSchema

func Example() {
	srv := handler.NewDefaultServer(es)
	srv.Use(gqlopentelemetry.Handler{})

	http.Handle("/query", srv)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
