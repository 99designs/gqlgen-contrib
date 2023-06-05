# Awesome gqlgen community contributions
[99designs/gqlgen](https://github.com/99designs/gqlgen) is both a library and a tool for generating Go GraphQL servers from GrahpQL schema.

The wider community has created many others tools and libraries to enhance and extend gqlgen.

This README contains a curated list of awesome [99designs/gqlgen](https://github.com/99designs/gqlgen) community contributions, only some of which are in this repository.

## Command line

### gqlgen Plugins
- [99designs/gqlgen/plugin/modelgen](https://github.com/99designs/gqlgen/tree/master/plugin/modelgen) - Generate go model type from GraphQL schema
- [99designs/gqlgen/plugin/resolvergen](https://github.com/99designs/gqlgen/tree/master/plugin/resolvergen) - Generate go resolver from GraphQL schema
- [99designs/gqlgen/plugin/federation](https://github.com/99designs/gqlgen/tree/master/plugin/federation) - Generate go Apollo Federation resolver from GraphQL schema
- [ravilushqa/otelgqlgen](https://github.com/ravilushqa/otelgqlgen) - OpenTelemetry instrumentation for 99designs/gqlgen
- [gqlapollotracing](https://github.com/99designs/gqlgen/tree/master/graphql/handler/apollotracing) - Apollo tracing instrumentation for 99designs/gqlgen
- [gqlopencensus](./gqlopencensus) - OpenCensus instrumentation for 99designs/gqlgen
- [gqlopentracing](./gqlopentracing) - OpenTracing instrumentation for 99designs/gqlgen
- [prometheus](./prometheus) - Prometheus instrumentation for 99designs/gqlgen
- [fasibio/autogql](https://github.com/fasibio/autogql) - CRUD SQL generator plugin
- [Warashi/compgen](https://github.com/Warashi/compgen) -  gqlgen plugin designed to simplify the generation of ComplexityRoot.
- [StevenACoffman/gqlgen_plugins](https://github.com/StevenACoffman/gqlgen-plugins) - some Internal Khan Academy plugins that might be useful to others.

### Client and Doc Generators
- [Khan/genqlient](https://github.com/Khan/genqlient) - Generate go GraphQL client from GraphQL query
- [infiotinc/gqlgenc](https://github.com/infiotinc/gqlgenc) - Generate go GraphQL client from GraphQL query
- [Yamashou/gqlgenc](https://github.com/Yamashou/gqlgenc) - Generate go GraphQL client from GraphQL query
- [Code-Hex/gqldoc](https://github.com/Code-Hex/gqldoc) - Generate Markdown GraphQL documents from GraphQL schema

###  Formatter
- [gqlgo/gqlfmt](https://github.com/gqlgo/gqlfmt) - Format GraphQL file
- [gqlgo/querystring](https://github.com/gqlgo/querystring) -
	`querystring` finds a GraphQL query in your files.

### Linter
- [gqlgo/lackid](https://github.com/gqlgo/lackid) - Detect lack of id in GrahpQL query
- [gqlgo/deprecatedquery](https://github.com/gqlgo/deprecatedquery) - About
	`deprecatedquery` finds a deprecated query in your GraphQL query files.
- [gqlgo/optionalschema](https://github.com/gqlgo/optionalschema) - About
	`optionalschema` finds optional fields and arguments in your GraphQL schema

## Library

### GraphQL

- [vektah/gqlparser](https://github.com/vektah/gqlparser) - GrahpQL schema and query parser
- [Yamashou/gqlgenc/introspection](https://github.com/Yamashou/gqlgenc/tree/master/introspection) - GraphQL instrospection query parser
- [Yamashou/gqlgenc/graphqljson](https://github.com/Yamashou/gqlgenc/tree/master/graphqljson) - GrahpQL JSON encoder and decoder

### Static analysis
- [gqlgo/gqlanaylsis](https://github.com/gqlgo/gqlanalysis) - Static analysis Go library for GraphQL, inspired by [go/analysis](https://pkg.go.dev/golang.org/x/tools/go/analysis)


