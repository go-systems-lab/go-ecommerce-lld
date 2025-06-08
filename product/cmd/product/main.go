package main

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/go-systems-lab/go-ecommerce-lld/product"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ElasticsearchURL      string `envconfig:"ELASTICSEARCH_URL" default:"http://localhost:9200"`
	Port                  int    `envconfig:"PORT" default:"8080"`
	KafkaBootstrapServers string `envconfig:"KAFKA_BOOTSTRAP_SERVERS" default:"kafka:9092"`
}

func main() {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("failed to process envconfig: %v", err)
	}

	var r product.Repository
	var producer sarama.AsyncProducer
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	producer, err := sarama.NewAsyncProducer([]string{cfg.KafkaBootstrapServers}, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer producer.Close()

	r, err = product.NewElasticRepository(cfg.ElasticsearchURL)
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}
	defer r.Close()

	log.Printf("starting product service on port %d", cfg.Port)
	s := product.NewProductService(r, producer)
	log.Fatal(product.ListenGRPC(s, cfg.Port))
}
