package account

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Close()
	PutAccount(ctx context.Context, a Account) (*Account, error)
	GetAccountByEmail(ctx context.Context, email string) (*Account, error)
	GetAccountByID(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &postgresRepository{db: db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) Ping() error {
	return r.db.Ping(context.Background())
}

func (r *postgresRepository) PutAccount(ctx context.Context, a Account) (*Account, error) {
	query := `
		INSERT INTO accounts (id, name, email, password)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET name = $2, email = $3, password = $4
	`

	_, err := r.db.Exec(ctx, query, a.ID, a.Name, a.Email, a.Password)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (r *postgresRepository) GetAccountByEmail(ctx context.Context, email string) (*Account, error) {
	query := `
		SELECT id, name, email, password
		FROM accounts
		WHERE email = $1
	`

	row := r.db.QueryRow(ctx, query, email)

	var a Account
	if err := row.Scan(&a.ID, &a.Name, &a.Email, &a.Password); err != nil {
		return nil, err
	}

	return &a, nil
}

func (r *postgresRepository) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	query := `
		SELECT id, name, email, password
		FROM accounts
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var a Account
	if err := row.Scan(&a.ID, &a.Name, &a.Email, &a.Password); err != nil {
		return nil, err
	}

	return &a, nil
}

func (r *postgresRepository) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	query := `
		SELECT id, name, email, password
		FROM accounts
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, take, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var a Account
		if err := rows.Scan(&a.ID, &a.Name, &a.Email, &a.Password); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}

	return accounts, nil
}
