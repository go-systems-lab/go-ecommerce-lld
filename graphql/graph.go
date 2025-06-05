package main

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/go-systems-lab/go-ecommerce-lld/account"
)

type Server struct {
	accountClient *account.Client
}

func NewGraphQLServer() (*Server, error) {
	// Connect to account service
	// TODO: Use env variable
	accountClient, err := account.NewClient("localhost:50051")
	if err != nil {
		return nil, err
	}

	return &Server{
		accountClient: accountClient,
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
