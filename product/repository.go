package product

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v9"
)

var (
	ErrNotFound = errors.New("entity not found")
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, p Product) error
	GetProductById(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip, take uint64) ([]Product, error)
	ListProductsWithIds(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error)
	UpdateProduct(ctx context.Context, updatedProduct Product) error
	DeleteProduct(ctx context.Context, productId string) error
}

type elasticRepository struct {
	client *elasticsearch.Client
}

type ProductDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return nil, err
	}

	return &elasticRepository{client: client}, nil
}

func (r *elasticRepository) Close() {
	// Elasticsearch client doesn't require explicit closing
	// The underlying HTTP client will be garbage collected
}

func (r *elasticRepository) PutProduct(ctx context.Context, p Product) error {
	doc := ProductDocument{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}

	docBytes, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	_, err = r.client.Index(
		"catalog",
		bytes.NewReader(docBytes),
		r.client.Index.WithDocumentID(p.ID),
		r.client.Index.WithContext(ctx),
	)
	return err
}

func (r *elasticRepository) GetProductById(ctx context.Context, id string) (*Product, error) {
	res, err := r.client.Get(
		"catalog",
		id,
		r.client.Get.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, ErrNotFound
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	found, ok := result["found"].(bool)
	if !ok || !found {
		return nil, ErrNotFound
	}

	source, ok := result["_source"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid source format")
	}

	sourceBytes, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}

	var product ProductDocument
	if err := json.Unmarshal(sourceBytes, &product); err != nil {
		return nil, err
	}

	return &Product{
		ID:          id,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}, nil
}

func (r *elasticRepository) ListProducts(ctx context.Context, skip, take uint64) ([]Product, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"from": skip,
		"size": take,
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithIndex("catalog"),
		r.client.Search.WithBody(bytes.NewReader(queryBytes)),
		r.client.Search.WithContext(ctx),
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid hits format")
	}

	hitsArray, ok := hits["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid hits array format")
	}

	var products []Product
	for _, hit := range hitsArray {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}

		id, _ := hitMap["_id"].(string)
		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			continue
		}

		sourceBytes, err := json.Marshal(source)
		if err != nil {
			continue
		}

		var product ProductDocument
		if err := json.Unmarshal(sourceBytes, &product); err == nil {
			products = append(products, Product{
				ID:          id,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			})
		}
	}
	return products, nil
}

func (r *elasticRepository) ListProductsWithIds(ctx context.Context, ids []string) ([]Product, error) {
	docs := make([]map[string]interface{}, len(ids))
	for i, id := range ids {
		docs[i] = map[string]interface{}{
			"_index": "catalog",
			"_id":    id,
		}
	}

	mgetQuery := map[string]interface{}{
		"docs": docs,
	}

	queryBytes, err := json.Marshal(mgetQuery)
	if err != nil {
		return nil, err
	}

	res, err := r.client.Mget(
		bytes.NewReader(queryBytes),
		r.client.Mget.WithContext(ctx),
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	docsArray, ok := result["docs"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid docs format")
	}

	var products []Product
	for _, doc := range docsArray {
		docMap, ok := doc.(map[string]interface{})
		if !ok {
			continue
		}

		found, ok := docMap["found"].(bool)
		if !ok || !found {
			continue
		}

		id, _ := docMap["_id"].(string)
		source, ok := docMap["_source"].(map[string]interface{})
		if !ok {
			continue
		}

		sourceBytes, err := json.Marshal(source)
		if err != nil {
			continue
		}

		var product ProductDocument
		if err := json.Unmarshal(sourceBytes, &product); err == nil {
			products = append(products, Product{
				ID:          id,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			})
		}
	}
	return products, nil
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error) {
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"name", "description"},
			},
		},
		"from": skip,
		"size": take,
	}

	queryBytes, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithIndex("catalog"),
		r.client.Search.WithBody(bytes.NewReader(queryBytes)),
		r.client.Search.WithContext(ctx),
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid hits format")
	}

	hitsArray, ok := hits["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid hits array format")
	}

	var products []Product
	for _, hit := range hitsArray {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}

		id, _ := hitMap["_id"].(string)
		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			continue
		}

		sourceBytes, err := json.Marshal(source)
		if err != nil {
			continue
		}

		var product ProductDocument
		if err := json.Unmarshal(sourceBytes, &product); err == nil {
			products = append(products, Product{
				ID:          id,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			})
		}
	}
	return products, nil
}

func (r *elasticRepository) UpdateProduct(ctx context.Context, updatedProduct Product) error {
	doc := map[string]interface{}{
		"doc": ProductDocument{
			Name:        updatedProduct.Name,
			Description: updatedProduct.Description,
			Price:       updatedProduct.Price,
		},
	}

	docBytes, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	_, err = r.client.Update(
		"catalog",
		updatedProduct.ID,
		bytes.NewReader(docBytes),
		r.client.Update.WithContext(ctx),
	)
	return err
}

func (r *elasticRepository) DeleteProduct(ctx context.Context, productId string) error {
	_, err := r.client.Delete(
		"catalog",
		productId,
		r.client.Delete.WithContext(ctx),
	)
	return err
}
