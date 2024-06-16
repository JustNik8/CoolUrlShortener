package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	envKey = "ENV"

	urlServiceHostKey = "URL_SERVICE_HOST"
	urlServicePortKey = "URL_SERVICE_PORT"

	analyticsServiceHostKey = "ANALYTICS_SERVICE_HOST"
	analyticsServicePortKey = "ANALYTICS_SERVICE_PORT"

	serverDomainKey = "SERVER_DOMAIN"

	rateLimitTokenPerSecondKey = "RATE_LIMIT_TOKEN_PER_SECOND"
	rateLimitBurstSizeKey      = "RATE_LIMIT_BURST_SIZE"
)

type Config struct {
	Env                    string
	ServerDomain           string
	UrlServiceConfig       UrlServiceConfig
	AnalyticsServiceConfig AnalyticsServiceConfig
	RateLimitConfig        RateLimitConfig
}

type AnalyticsServiceConfig struct {
	Host string
	Port string
}

type UrlServiceConfig struct {
	Host string
	Port string
}

type RateLimitConfig struct {
	TokensPerSecond float64
	BurstSize       int
}

func ParseConfig() (Config, error) {
	env := os.Getenv(envKey)
	if env == "" {
		return Config{}, fmt.Errorf("you did not provide env: %s", envKey)
	}

	urlServiceHost := os.Getenv(urlServiceHostKey)
	if urlServiceHost == "" {
		return Config{}, fmt.Errorf("you did not provide env: %s", urlServiceHostKey)
	}

	urlServicePort := os.Getenv(urlServicePortKey)
	if urlServicePort == "" {
		return Config{}, fmt.Errorf("you did not provide env: %s", urlServicePortKey)
	}

	analyticsServiceHost := os.Getenv(analyticsServiceHostKey)
	if analyticsServiceHost == "" {
		return Config{}, fmt.Errorf("you did not provide env: %s", analyticsServiceHostKey)
	}

	analyticsServicePort := os.Getenv(analyticsServicePortKey)
	if analyticsServicePort == "" {
		return Config{}, fmt.Errorf("you did not provide env: %s", analyticsServicePortKey)
	}

	serverDomain := os.Getenv(serverDomainKey)
	if serverDomain == "" {
		return Config{}, fmt.Errorf("you did not provide env: %s", serverDomainKey)
	}

	rateLimitTokenPerSecondRaw := os.Getenv(rateLimitTokenPerSecondKey)
	if rateLimitTokenPerSecondRaw == "" {
		return Config{}, fmt.Errorf("you did not provide env: %s", rateLimitTokenPerSecondKey)
	}
	rateLimitTokenPerSecond, err := strconv.ParseFloat(rateLimitTokenPerSecondRaw, 64)
	if err != nil {
		return Config{}, err
	}

	rateLimitBurstSizeRaw := os.Getenv(rateLimitBurstSizeKey)
	if rateLimitBurstSizeRaw == "" {
		return Config{}, fmt.Errorf("you did not provide env: %s", rateLimitBurstSizeKey)
	}
	rateLimitBurstSize, err := strconv.Atoi(rateLimitBurstSizeRaw)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Env:          env,
		ServerDomain: serverDomain,
		UrlServiceConfig: UrlServiceConfig{
			Host: urlServiceHost,
			Port: urlServicePort,
		},
		AnalyticsServiceConfig: AnalyticsServiceConfig{
			Host: analyticsServiceHost,
			Port: analyticsServicePort,
		},
		RateLimitConfig: RateLimitConfig{
			TokensPerSecond: rateLimitTokenPerSecond,
			BurstSize:       rateLimitBurstSize,
		},
	}, nil
}
