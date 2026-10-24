package config

import (
	"fmt"
	"os"
)

type Env struct {
	Port string
}

func NewEnv() *Env {
	port, has := os.LookupEnv("port")
	if !has {
		fmt.Println("Port not found")
	}

	return &Env{
		Port: port,
	}

}
