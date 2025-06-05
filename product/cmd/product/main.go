package main

import (
	"log"

	"github.com/go-systems-lab/go-ecommerce-lld/product"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ElasticsearchURL string `envconfig:"ELASTICSEARCH_URL" default:"http://localhost:9200"`
	Port             int    `envconfig:"PORT" default:"8080"`
}

func main() {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("failed to process envconfig: %v", err)
	}

	var r product.Repository
	r, err := product.NewElasticRepository(cfg.ElasticsearchURL)
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}
	defer r.Close()

	log.Printf("starting product service on port %d", cfg.Port)
	s := product.NewProductService(r)
	log.Fatal(product.ListenGRPC(s, cfg.Port))
}
