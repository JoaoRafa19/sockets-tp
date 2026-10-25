package main

import (
	"sockets-tc/config"
	tcp "sockets-tc/internal/adapters/tcp"
)

func main() {
	env := config.NewEnv()

	server := tcp.NewServer(env)

	if err := server.Start(); err != nil {
		panic(err)
	}
}
