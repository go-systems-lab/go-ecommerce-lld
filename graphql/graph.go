package main

import (
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-systems-lab/go-ecommerce-lld/account"
	"github.com/go-systems-lab/go-ecommerce-lld/order"
	"github.com/go-systems-lab/go-ecommerce-lld/product"
	"github.com/go-systems-lab/go-ecommerce-lld/recommender"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

type Server struct {
	accountClient     *account.Client
	productClient     *product.Client
	orderClient       *order.Client
	recommenderClient *recommender.Client
}

func NewGraphQLServer(
	accountServiceURL,
	productServiceURL,
	orderServiceURL,
	recommenderServiceURL string,
) (*Server, error) {
	// Connect to account service
	accountClient, err := account.NewClient(accountServiceURL)
	if err != nil {
		return nil, err
	}

	// Connect to product service
	productClient, err := product.NewClient(productServiceURL)
	if err != nil {
		return nil, err
	}

	// Connect to order service
	orderClient, err := order.NewClient(orderServiceURL)
	if err != nil {
		return nil, err
	}

	// Connect to recommender service
	recommenderClient, err := recommender.NewClient(recommenderServiceURL)
	if err != nil {
		return nil, err
	}

	return &Server{
		accountClient:     accountClient,
		productClient:     productClient,
		orderClient:       orderClient,
		recommenderClient: recommenderClient,
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
