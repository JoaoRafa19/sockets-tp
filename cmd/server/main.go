package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sockets-tc/config"
	"sockets-tc/internal/adapters/repository/postgres/db"
	tcp "sockets-tc/internal/adapters/tcp"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	env, err := config.NewEnv()
	if err != nil {
		fmt.Println("error on env definition")
		os.Exit(1)
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, env.DatabaseURL)
	if err != nil {
		log.Fatalf("db config: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("db ping: %v", err)
	}

	_ = db.New(pool)

	err = pool.Ping(ctx)
	fmt.Println("ping error", err)
	// repo     := repository.New(pool, queries)
	// services := service.New...(repo)
	// handler  := tcp.NewHandler(services...)
	// server   := tcp.NewServer(env, handler)
	// log.Fatal(server.Start())

	server := tcp.NewServer(env)

	if err := server.Start(); err != nil {
		panic(err)
	}
}
