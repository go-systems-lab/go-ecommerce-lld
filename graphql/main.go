package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
	"github.com/vektah/gqlparser/v2/ast"
)

type AppConfig struct {
	AccountServiceURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	ProductServiceURL string `envconfig:"PRODUCT_SERVICE_URL"`
	OrderServiceURL   string `envconfig:"ORDER_SERVICE_URL"`
	Port              string `envconfig:"PORT"`
}

func main() {
	var cfg AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	s, err := NewGraphQLServer(cfg.AccountServiceURL, cfg.ProductServiceURL, cfg.OrderServiceURL)
	if err != nil {
		log.Fatal(err)
	}

	srv := handler.New(s.toExecutableSchema())

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
