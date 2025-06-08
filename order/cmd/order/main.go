package main

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/go-systems-lab/go-ecommerce-lld/order"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseURL           string `envconfig:"DATABASE_URL"`
	AccountURL            string `envconfig:"ACCOUNT_URL"`
	ProductURL            string `envconfig:"PRODUCT_URL"`
	Port                  int    `envconfig:"PORT"`
	KafkaBootstrapServers string `envconfig:"KAFKA_BOOTSTRAP_SERVERS" default:"kafka:9092"`
}

func main() {
	var cfg Config
	var producer sarama.AsyncProducer
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	producer, err := sarama.NewAsyncProducer([]string{cfg.KafkaBootstrapServers}, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer producer.Close()

	r, err := order.NewPostgresRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer r.Close()

	s := order.NewOrderService(r, producer)

	log.Fatal(order.ListenGRPC(s, cfg.AccountURL, cfg.ProductURL, cfg.Port))
}
