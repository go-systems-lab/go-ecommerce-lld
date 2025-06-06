package order

import (
	"context"
	"fmt"
	"log"
	"net"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/go-systems-lab/go-ecommerce-lld/account"
	"github.com/go-systems-lab/go-ecommerce-lld/order/pb"
	"github.com/go-systems-lab/go-ecommerce-lld/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedOrderServiceServer
	service       Service
	accountClient *account.Client
	productClient *product.Client
}

func ListenGRPC(s Service, accountURL, productURL string, port int) error {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return err
	}
	defer accountClient.Close()

	productClient, err := product.NewClient(productURL)
	if err != nil {
		return err
	}
	defer productClient.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	srv := grpc.NewServer()
	pb.RegisterOrderServiceServer(srv, &grpcServer{
		service:       s,
		accountClient: accountClient,
		productClient: productClient,
	})
	reflection.Register(srv)
	return srv.Serve(lis)
}

func (s *grpcServer) PostOrder(ctx context.Context, request *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, request.AccountId)
	if err != nil {
		log.Println("Error getting account", err)
		return nil, err
	}
	var productIDs []string
	for _, p := range request.Products {
		productIDs = append(productIDs, p.ProductId)
	}
	orderedProducts, err := s.productClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("Error getting ordered products", err)
		return nil, err
	}

	var products []OrderedProduct

	for _, p := range orderedProducts {
		productObj := OrderedProduct{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    0,
		}
		for _, requestProduct := range request.Products {
			if requestProduct.ProductId == p.ID {
				productObj.Quantity = requestProduct.Quantity
				break
			}
		}

		if productObj.Quantity != 0 {
			products = append(products, productObj)
		}
	}

	order, err := s.service.PostOrder(ctx, request.AccountId, products)
	if err != nil {
		log.Println("Error posting order", err)
		return nil, err
	}

	orderProto := &pb.Order{
		Id:         order.ID,
		AccountId:  order.AccountID,
		TotalPrice: order.TotalPrice,
		Products:   []*pb.Order_OrderProduct{},
	}
	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()

	// Use the original products with full details instead of order.Products
	for _, p := range products {
		orderProto.Products = append(orderProto.Products, &pb.Order_OrderProduct{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
		})
	}
	return &pb.PostOrderResponse{
		Order: orderProto,
	}, nil
}

func (s *grpcServer) GetOrdersForAccount(ctx context.Context, request *pb.GetOrdersForAccountRequest) (*pb.GetOrdersForAccountResponse, error) {
	accountOrders, err := s.service.GetOrdersForAccount(ctx, request.AccountId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Taking unique products. We use set to avoid repeating
	productIDsSet := mapset.NewSet[string]()
	for _, o := range accountOrders {
		for _, p := range o.Products {
			productIDsSet.Add(p.ID)
		}
	}

	productIDs := productIDsSet.ToSlice()

	products, err := s.productClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("Error getting account products: ", err)
		return nil, err
	}

	// Collecting orders

	var orders []*pb.Order
	for _, o := range accountOrders {
		// Encode order
		op := &pb.Order{
			AccountId:  o.AccountID,
			Id:         o.ID,
			TotalPrice: 0, // We'll calculate this correctly
			Products:   []*pb.Order_OrderProduct{},
		}
		op.CreatedAt, _ = o.CreatedAt.MarshalBinary()

		// Decorate orders with products and calculate correct total price
		var calculatedTotalPrice float64
		for _, orderedProduct := range o.Products {
			// Populate product fields
			for _, p := range products {
				if p.ID == orderedProduct.ID {
					orderedProduct.Name = p.Name
					orderedProduct.Description = p.Description
					orderedProduct.Price = p.Price
					// Calculate total price correctly: price * quantity
					calculatedTotalPrice += p.Price * float64(orderedProduct.Quantity)
					break
				}
			}

			op.Products = append(op.Products, &pb.Order_OrderProduct{
				Id:          orderedProduct.ID,
				Name:        orderedProduct.Name,
				Description: orderedProduct.Description,
				Price:       orderedProduct.Price,
				Quantity:    orderedProduct.Quantity,
			})
		}

		// Set the correctly calculated total price
		op.TotalPrice = calculatedTotalPrice
		orders = append(orders, op)
	}
	return &pb.GetOrdersForAccountResponse{Orders: orders}, nil
}
