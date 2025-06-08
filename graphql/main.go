package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/go-systems-lab/go-ecommerce-lld/account"
	"github.com/go-systems-lab/go-ecommerce-lld/pkg/middleware"
	"github.com/kelseyhightower/envconfig"
	"github.com/vektah/gqlparser/v2/ast"
)

type AppConfig struct {
	AccountServiceURL        string `envconfig:"ACCOUNT_SERVICE_URL"`
	ProductServiceURL        string `envconfig:"PRODUCT_SERVICE_URL"`
	OrderServiceURL          string `envconfig:"ORDER_SERVICE_URL"`
	RecommendationServiceURL string `envconfig:"RECOMMENDATION_SERVICE_URL"`
	Port                     string `envconfig:"PORT"`
	SecretKey                string `envconfig:"SECRET_KEY"`
	Issuer                   string `envconfig:"ISSUER"`
}

func main() {
	var cfg AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	s, err := NewGraphQLServer(cfg.AccountServiceURL, cfg.ProductServiceURL, cfg.OrderServiceURL, cfg.RecommendationServiceURL)
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

	engine := gin.Default()

	engine.Use(middleware.GinContextToContextMiddleware())

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	engine.POST("/graphql", AuthorizeJWT(account.NewJwtService(cfg.SecretKey, cfg.Issuer)), gin.WrapH(srv))
	engine.GET("/playground", gin.WrapH(playground.Handler("GraphQL playground", "/graphql")))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.Port)
	log.Fatal(engine.Run(":" + cfg.Port))
}
