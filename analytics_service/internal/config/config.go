package config

import (
	"fmt"
	"os"
)

const (
	envKey = "ENV"

	clickhouseUsernameKey = "CLICKHOUSE_USERNAME"
	clickhousePasswordKey = "CLICKHOUSE_PASSWORD"
	clickhouseHostKey     = "CLICKHOUSE_HOST"
	clickhousePortKey     = "CLICKHOUSE_PORT"
	clickhouseDatabaseKey = "CLICKHOUSE_DATABASE"
)

type Config struct {
	Env              string
	ClickhouseConfig ClickhouseConfig
}

type ClickhouseConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func ParseConfig() (Config, error) {
	env := os.Getenv(envKey)
	if env == "" {
		msg := fmt.Sprintf("You did not provide env: %s", envKey)
		panic(msg)
	}

	username := os.Getenv(clickhouseUsernameKey)
	if username == "" {
		return Config{}, fmt.Errorf("you did not provide env: %s", clickhouseUsernameKey)
	}
	password := os.Getenv(clickhousePasswordKey)
	if password == "" {
		return Config{}, fmt.Errorf("you did not provice env: %s", clickhousePasswordKey)
	}
	host := os.Getenv(clickhouseHostKey)
	if host == "" {
		return Config{}, fmt.Errorf("you did not provide env: %s", clickhouseHostKey)
	}
	port := os.Getenv(clickhousePortKey)
	if port == "" {
		return Config{}, fmt.Errorf("you did not provide env: %s", clickhousePortKey)
	}

	database := os.Getenv(clickhouseDatabaseKey)
	if database == "" {
		return Config{}, fmt.Errorf("you did not provice env: %s", clickhouseDatabaseKey)
	}

	return Config{
		Env: env,
		ClickhouseConfig: ClickhouseConfig{
			Username: username,
			Password: password,
			Host:     host,
			Port:     port,
			Database: database,
		},
	}, nil
}
