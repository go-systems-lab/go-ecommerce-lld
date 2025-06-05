package main

import "context"

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) CreateAccount(ctx context.Context, account AccountInput) (*Account, error) {
	// TODO: Implement
	return nil, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, product ProductInput) (*Product, error) {
	// TODO: Implement
	return nil, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, order OrderInput) (*Order, error) {
	// TODO: Implement
	return nil, nil
}
