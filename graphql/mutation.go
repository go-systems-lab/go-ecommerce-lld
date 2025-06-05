package main

import (
	"context"
)

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) CreateAccount(ctx context.Context, account AccountInput) (*Account, error) {
	// Create account via microservice
	createdAccount, err := r.server.accountClient.PostAccount(ctx, account.Name)
	if err != nil {
		return nil, err
	}

	return &Account{
		ID:   createdAccount.ID,
		Name: createdAccount.Name,
	}, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, product ProductInput) (*Product, error) {
	// TODO: Implement
	return nil, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, order OrderInput) (*Order, error) {
	// TODO: Implement
	return nil, nil
}
