package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Env struct {
	Port        string
	DatabaseURL string
}

func NewEnv() (*Env, error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("ler .env: %w", err)
	}
	viper.AutomaticEnv()
	_ = viper.BindEnv("db_url", "db_url", "DATABASE_URL") // nome usado no docker-compose
	viper.SetDefault("port", ":3000")

	env := &Env{
		Port:        viper.GetString("port"),
		DatabaseURL: viper.GetString("db_url"),
	}

	if env.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL não definida")
	}
	return env, nil
}
