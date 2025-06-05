package main

import (
	"log"

	"github.com/go-systems-lab/go-ecommerce-lld/account"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	DatabaseURL string `envconfig:"DATABASE_URL" default:"postgres://postgres:postgres@localhost:5432/ecommerce_account?sslmode=disable"`
	Port        int    `envconfig:"PORT" default:"50051"`
}

func main() {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("failed to process envconfig: %v", err)
	}

	log.Printf("connecting to database: %s", cfg.DatabaseURL)

	r, err := account.NewPostgresRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}
	defer r.Close()

	log.Printf("starting account service on port %d", cfg.Port)
	s := account.NewService(r)
	log.Fatal(account.ListenGRPC(s, cfg.Port))
}
