package main

import (
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-systems-lab/go-ecommerce-lld/account"
	"github.com/go-systems-lab/go-ecommerce-lld/product"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

type Server struct {
	accountClient *account.Client
	productClient *product.Client
}

func NewGraphQLServer() (*Server, error) {
	// Connect to account service
	accountServiceURL := getEnv("ACCOUNT_SERVICE_URL", "localhost:8080")
	accountClient, err := account.NewClient(accountServiceURL)
	if err != nil {
		return nil, err
	}

	// Connect to product service
	productServiceURL := getEnv("PRODUCT_SERVICE_URL", "localhost:8080")
	productClient, err := product.NewClient(productServiceURL)
	if err != nil {
		return nil, err
	}

	return &Server{
		accountClient: accountClient,
		productClient: productClient,
	}, nil
}

func (s *Server) Mutation() MutationResolver {
	return &mutationResolver{server: s}
}

func (s *Server) Query() QueryResolver {
	return &queryResolver{server: s}
}

func (s *Server) Account() AccountResolver {
	return &accountResolver{server: s}
}

func (s *Server) toExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: s,
	})
}
