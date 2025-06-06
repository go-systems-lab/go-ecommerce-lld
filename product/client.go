package product

import (
	"context"

	"github.com/go-systems-lab/go-ecommerce-lld/product/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.ProductServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := pb.NewProductServiceClient(conn)

	return &Client{
		conn:    conn,
		service: c,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostProduct(ctx context.Context, name, description string, price float64, accountId string) (*Product, error) {
	r, err := c.service.PostProduct(ctx, &pb.CreateProductRequest{
		Name:        name,
		Description: description,
		Price:       price,
		AccountId:   accountId,
	})

	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          r.Product.Id,
		Name:        r.Product.Name,
		Description: r.Product.Description,
		Price:       r.Product.Price,
		AccountID:   r.Product.GetAccountId(),
	}, nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	r, err := c.service.GetProduct(ctx, &pb.ProductByIdRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          r.Product.Id,
		Name:        r.Product.Name,
		Description: r.Product.Description,
		Price:       r.Product.Price,
		AccountID:   r.Product.GetAccountId(),
	}, nil
}

func (c *Client) GetProducts(ctx context.Context, skip, take uint64, ids []string, query string) ([]Product, error) {
	r, err := c.service.GetProducts(ctx, &pb.GetProductsRequest{
		Skip:  skip,
		Take:  take,
		Ids:   ids,
		Query: query,
	})

	if err != nil {
		return nil, err
	}

	var products []Product
	for _, p := range r.Products {
		products = append(products, Product{
			ID:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			AccountID:   p.GetAccountId(),
		})
	}

	return products, nil
}

func (c *Client) UpdateProduct(ctx context.Context, id, name, description string, price float64, accountId string) (*Product, error) {
	res, err := c.service.UpdateProduct(ctx, &pb.UpdateProductRequest{
		Id:          id,
		Name:        name,
		Description: description,
		Price:       price,
		AccountId:   accountId,
	})
	if err != nil {
		return nil, err
	}
	return &Product{
		res.Product.Id,
		res.Product.Name,
		res.Product.Description,
		res.Product.Price,
		res.Product.GetAccountId(),
	}, nil
}

func (c *Client) DeleteProduct(ctx context.Context, id string, accountId string) error {
	_, err := c.service.DeleteProduct(ctx, &pb.DeleteProductRequest{
		ProductId: id,
		AccountId: accountId,
	})
	return err
}
