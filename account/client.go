package account

import (
	"context"

	"github.com/go-systems-lab/go-ecommerce-lld/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	C := pb.NewAccountServiceClient(conn)

	return &Client{
		conn:    conn,
		service: C,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	r, err := c.service.PostAccount(ctx, &pb.PostAccountRequest{Name: name})

	if err != nil {
		return nil, err
	}

	return &Account{
		ID:   r.Account.GetId(),
		Name: r.Account.GetName(),
	}, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	r, err := c.service.GetAccount(ctx, &pb.GetAccountRequest{Id: id})

	if err != nil {
		return nil, err
	}

	return &Account{
		ID:   r.Account.GetId(),
		Name: r.Account.GetName(),
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip, take uint64) ([]Account, error) {
	r, err := c.service.GetAccounts(ctx, &pb.GetAccountsRequest{Skip: skip, Take: take})

	if err != nil {
		return nil, err
	}

	var accounts []Account

	for _, a := range r.Accounts {
		accounts = append(accounts, Account{
			ID:   a.GetId(),
			Name: a.GetName(),
		})
	}

	return accounts, nil
}
