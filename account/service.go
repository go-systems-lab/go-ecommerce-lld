package account

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Register(ctx context.Context, name, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	GetAccountByID(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type Account struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type accountService struct {
	repository  Repository
	authService JwtService
}

func NewService(repository Repository, authService JwtService) Service {
	return &accountService{repository: repository, authService: authService}
}

func (s accountService) Register(ctx context.Context, name, email, password string) (string, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return "", err
	}

	a := Account{
		ID:       uuid.New().String(),
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	account, err := s.repository.PutAccount(ctx, a)
	if err != nil {
		return "", err
	}

	token, err := s.authService.GenerateToken(account.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s accountService) Login(ctx context.Context, email, password string) (string, error) {
	account, err := s.repository.GetAccountByEmail(ctx, email)
	if err == nil && VerifyPassword(account.Password, password) {
		token, err := s.authService.GenerateToken(account.ID)
		if err != nil {
			return "", err
		}
		return token, nil
	}

	return "", err
}

func (s accountService) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	return s.repository.GetAccountByID(ctx, id)
}

func (s accountService) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	return s.repository.ListAccounts(ctx, skip, take)
}
