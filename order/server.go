package order

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

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
	log.Printf("PostOrder gRPC handler called with %d products", len(request.Products))
	for i, p := range request.Products {
		log.Printf("Request product %d: ID=%s, Quantity=%d", i, p.Id, p.Quantity)
	}

	_, err := s.accountClient.GetAccount(ctx, request.AccountId)
	if err != nil {
		log.Println("Error getting account", err)
		return nil, err
	}

	log.Printf("Fetching product details (optimistic approach)...")
	var orderedProducts []product.Product
	var asyncProductIDs []string

	for _, p := range request.Products {
		log.Printf("Fetching product with ID: %s", p.Id)

		// Try immediate fetch - don't block if not found
		fetchedProduct, err := s.productClient.GetProduct(ctx, p.Id)
		if err == nil {
			log.Printf("✅ Product found immediately: %s", fetchedProduct.Name)
			orderedProducts = append(orderedProducts, *fetchedProduct)
		} else {
			log.Printf("⏳ Product %s not immediately available, will retry asynchronously", p.Id)
			asyncProductIDs = append(asyncProductIDs, p.Id)

			// Create placeholder product for order (optimistic)
			orderedProducts = append(orderedProducts, product.Product{
				ID:    p.Id,
				Name:  "Product (Processing...)",
				Price: 0.0, // Placeholder price
			})
		}
	}

	// Start async processing for products that weren't immediately available
	if len(asyncProductIDs) > 0 {
		log.Printf("🚀 Starting async interaction event processing for %d products", len(asyncProductIDs))
		go s.processInteractionEventsAsync(ctx, request.AccountId, asyncProductIDs)
	}

	log.Printf("Retrieved %d products from product service", len(orderedProducts))

	var products []OrderedProduct
	var calculatedTotalPrice float64

	for _, p := range orderedProducts {
		productObj := OrderedProduct{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    0,
		}
		for _, requestProduct := range request.Products {
			if requestProduct.Id == p.ID {
				productObj.Quantity = requestProduct.Quantity
				break
			}
		}

		if productObj.Quantity != 0 {
			// Calculate total price: price * quantity
			calculatedTotalPrice += p.Price * float64(productObj.Quantity)
			products = append(products, productObj)
		}
	}

	order, err := s.service.PostOrder(ctx, request.AccountId, calculatedTotalPrice, products)
	if err != nil {
		log.Println("Error posting order", err)
		return nil, err
	}

	orderProto := &pb.Order{
		Id:         order.ID,
		AccountId:  order.AccountID,
		TotalPrice: order.TotalPrice,
		Products:   []*pb.OrderedProduct{},
	}
	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()

	// Use the original products with full details instead of order.Products
	for _, p := range products {
		orderProto.Products = append(orderProto.Products, &pb.OrderedProduct{
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
			TotalPrice: o.TotalPrice,
			Products:   []*pb.OrderedProduct{},
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

			op.Products = append(op.Products, &pb.OrderedProduct{
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

// processInteractionEventsAsync handles async interaction event processing with retries
func (s *grpcServer) processInteractionEventsAsync(ctx context.Context, accountID string, productIDs []string) {
	maxRetries := 3
	baseDelay := 2 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		allFound := true
		var validProducts []product.Product

		for _, productID := range productIDs {
			product, err := s.productClient.GetProduct(ctx, productID)
			if err != nil {
				log.Printf("🔄 Async attempt %d: Product %s not found yet", attempt+1, productID)
				allFound = false
				break
			}
			validProducts = append(validProducts, *product)
		}

		if allFound {
			log.Printf("✅ All async products found on attempt %d, sending interaction events", attempt+1)
			// Send interaction events for all found products
			for _, prod := range validProducts {
				if orderSvc, ok := s.service.(*orderService); ok {
					err := orderSvc.SendMessageToRecommender(Event{
						Type: "purchase",
						EventData: EventData{
							AccountId: accountID,
							ProductId: prod.ID,
						},
					}, "interaction_events")

					if err != nil {
						log.Printf("❌ Error sending async interaction event for product %s: %v", prod.ID, err)
					} else {
						log.Printf("✅ Successfully sent async interaction event for product %s", prod.ID)
					}
				}
			}
			return
		}

		if attempt < maxRetries-1 {
			waitTime := baseDelay * time.Duration(attempt+1)
			log.Printf("⏳ Retrying async interaction events in %v (attempt %d/%d)", waitTime, attempt+1, maxRetries)
			time.Sleep(waitTime)
		}
	}

	log.Printf("❌ Failed to process async interaction events after %d attempts", maxRetries)
}
