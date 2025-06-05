package main

import "github.com/99designs/gqlgen/graphql"

type Server struct {
}

func NewGraphQLServer() (*Server, error) {
	return &Server{}, nil
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
