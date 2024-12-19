package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	EthereumRPCURL string
	HTTPPort       int
	LogLevel       string
	WebhookURL     string
}

// NewConfig creates a default configuration
func NewConfig() *Config {
	return &Config{
		EthereumRPCURL: "https://ethereum-rpc.publicnode.com",
		HTTPPort:       8060,
		LogLevel:       "info",
		WebhookURL:     "https://example-webhook-url.com/notify", // Default webhook URL
	}
}

// LoadEnvironmentVariables loads configuration values from environment variables
func (c *Config) LoadEnvironmentVariables() {

	if rpcURL := os.Getenv("ETHEREUM_RPC_URL"); rpcURL != "" {
		c.EthereumRPCURL = rpcURL
	}

	if portStr := os.Getenv("HTTP_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			c.HTTPPort = port
		}
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		c.LogLevel = logLevel
	}

	if webhookURL := os.Getenv("WEBHOOK_URL"); webhookURL != "" {
		c.WebhookURL = webhookURL
	}
}
