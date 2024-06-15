package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	envKey              = "ENV"
	databaseUsernameKey = "DATABASE_USERNAME"
	databasePasswordKey = "DATABASE_PASSWORD"
	databaseHostKey     = "DATABASE_HOST"
	databasePortKey     = "DATABASE_PORT"
	databaseNameKey     = "DATABASE_NAME"

	redisHostKey     = "REDIS_HOST"
	redisPortKey     = "REDIS_PORT"
	redisPasswordKey = "REDIS_PASSWORD"

	kafkaAddrsKey = "KAFKA_ADDRS"
)

type Config struct {
	Env            string
	DatabaseConfig DatabaseConfig
	RedisConfig    RedisConfig
	KafkaConfig    KafkaConfig
}

type DatabaseConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type KafkaConfig struct {
	Addrs []string
}

func ParseConfig() (Config, error) {
	env := os.Getenv(envKey)
	if env == "" {
		msg := fmt.Sprintf("You did not provide env: %s", envKey)
		panic(msg)
	}

	dbUsername := os.Getenv(databaseUsernameKey)
	if dbUsername == "" {
		return Config{}, fmt.Errorf("you did not provide env: %s", databaseUsernameKey)
	}

	dbPassword := os.Getenv(databasePasswordKey)
	if dbPassword == "" {
		return Config{}, fmt.Errorf("you did not provide env: %s", databasePasswordKey)
	}

	dbHost := os.Getenv(databaseHostKey)
	if dbHost == "" {
		return Config{}, fmt.Errorf("you did not provide env: %s", databaseHostKey)
	}

	dbPort := os.Getenv(databasePortKey)
	if dbPort == "" {
		return Config{}, fmt.Errorf("you did not provide env: %s", databasePortKey)
	}

	dbName := os.Getenv(databaseNameKey)
	if dbName == "" {
		return Config{}, fmt.Errorf("you did not provide env: %s", databaseNameKey)
	}

	redisHost := os.Getenv(redisHostKey)
	if redisHost == "" {
		return Config{}, fmt.Errorf("you did not provide env: %s", redisHostKey)
	}

	redisPort := os.Getenv(redisPortKey)
	if redisPort == "" {
		return Config{}, fmt.Errorf("you did not provide env: %s", redisPortKey)
	}

	redisPassword := os.Getenv(redisPasswordKey)
	if redisPassword == "" {
		return Config{}, fmt.Errorf("you did not provide env: %s", redisPasswordKey)
	}

	kafkaAddrsRaw := os.Getenv(kafkaAddrsKey)
	if kafkaAddrsRaw == "" {
		return Config{}, fmt.Errorf("you did not provide env: %s", kafkaAddrsKey)
	}
	kafkaAddrs := strings.Split(kafkaAddrsRaw, ",")

	return Config{
		Env: env,
		DatabaseConfig: DatabaseConfig{
			Username: dbUsername,
			Password: dbPassword,
			Host:     dbHost,
			Port:     dbPort,
			Name:     dbName,
		},
		RedisConfig: RedisConfig{
			Host:     redisHost,
			Port:     redisPort,
			Password: redisPassword,
		},
		KafkaConfig: KafkaConfig{
			Addrs: kafkaAddrs,
		},
	}, nil
}
