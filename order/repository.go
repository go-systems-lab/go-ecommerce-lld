package order

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Close()
	PutOrder(ctx context.Context, order Order) error
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
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

func (r postgresRepository) PutOrder(ctx context.Context, order Order) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO orders (id, created_at, account_id, total_price)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			created_at = $2,
			account_id = $3,
			total_price = $4
	`

	_, err = tx.Exec(ctx, query, order.ID, order.CreatedAt, order.AccountID, order.TotalPrice)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO order_products (order_id, product_id, quantity)
		VALUES ($1, $2, $3)
	`

	for _, product := range order.Products {
		_, err = tx.Exec(ctx, query, order.ID, product.ID, product.Quantity)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r postgresRepository) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	query := `
		SELECT o.id, o.created_at, o.account_id, o.total_price, op.product_id, op.quantity
		FROM orders o
		JOIN order_products op ON o.id = op.order_id
		WHERE account_id = $1
		ORDER BY o.id
	`

	rows, err := r.db.Query(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	var products []OrderedProduct
	var lastOrderID string
	order := &Order{}
	orderedProduct := &OrderedProduct{}

	for rows.Next() {
		if err = rows.Scan(
			&order.ID,
			&order.CreatedAt,
			&order.AccountID,
			&order.TotalPrice,
			&orderedProduct.ID,
			&orderedProduct.Quantity,
		); err != nil {
			return nil, err
		}

		if lastOrderID != "" && lastOrderID != order.ID {
			order.Products = products
			orders = append(orders, *order)
			products = []OrderedProduct{}
		}

		products = append(products, OrderedProduct{
			ID:       orderedProduct.ID,
			Quantity: orderedProduct.Quantity,
		})

		lastOrderID = order.ID
	}

	if lastOrderID != "" {
		order.Products = products
		orders = append(orders, *order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
