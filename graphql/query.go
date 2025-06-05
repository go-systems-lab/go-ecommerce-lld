package main

import "context"

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	// TODO: Implement
	return nil, nil
}

func (r *queryResolver) Product(ctx context.Context, pagination *PaginationInput, query, id *string) ([]*Product, error) {
	// TODO: Implement
	return nil, nil
}
