package main

import (
	"log"

	"github.com/go-systems-lab/go-ecommerce-lld/order"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
	AccountURL  string `envconfig:"ACCOUNT_URL"`
	ProductURL  string `envconfig:"PRODUCT_URL"`
	Port        int    `envconfig:"PORT"`
}

func main() {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	var r order.Repository
	r, err := order.NewPostgresRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer r.Close()

	s := order.NewOrderService(r)

	log.Fatal(order.ListenGRPC(s, cfg.AccountURL, cfg.ProductURL, cfg.Port))
}
