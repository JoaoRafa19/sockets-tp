package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"sockets-tc/config"
	"sockets-tc/internal/adapters/postgres"
	"sockets-tc/internal/adapters/tcp"
	"sockets-tc/internal/service"
)

func main() {
	env, err := config.NewEnv()
	if err != nil {
		log.Fatalf("env: %v", err)
	}

	ctx := context.Background()

	// banco
	pool, err := pgxpool.New(ctx, env.DatabaseURL)
	if err != nil {
		log.Fatalf("db config: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("db ping: %v", err)
	}

	// repository -> services -> handler -> server
	repo := postgres.New(pool)

	eventService := service.NewEventService(repo)
	reservationService := service.NewReservationService(repo)

	handler := tcp.NewHandler(eventService, reservationService)
	server := tcp.NewServer(env, handler)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
