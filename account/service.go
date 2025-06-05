package account

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	PostAccount(ctx context.Context, name string) (*Account, error)
	GetAccountByID(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type accountService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &accountService{repository: repository}
}

func (s accountService) PostAccount(ctx context.Context, name string) (*Account, error) {
	a := Account{
		ID:   uuid.New().String(),
		Name: name,
	}

	if err := s.repository.PutAccount(ctx, a); err != nil {
		return nil, err
	}

	return &a, nil
}

func (s accountService) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	return s.repository.GetAccountByID(ctx, id)
}

func (s accountService) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	return s.repository.ListAccounts(ctx, skip, take)
}
