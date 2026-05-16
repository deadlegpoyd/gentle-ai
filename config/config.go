// Package config handles loading and validation of application configuration.
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration values for the gentle-ai application.
type Config struct {
	// Server settings
	Port int
	Host string

	// AI provider settings
	OpenAIAPIKey    string
	OpenAIModel     string
	MaxTokens       int
	Temperature     float64

	// Application settings
	Debug           bool
	LogLevel        string
	SystemPrompt    string
}

// Load reads configuration from environment variables and returns a Config.
// Required variables: OPENAI_API_KEY
// Optional variables: PORT, HOST, OPENAI_MODEL, MAX_TOKENS, TEMPERATURE, DEBUG, LOG_LEVEL, SYSTEM_PROMPT
func Load() (*Config, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	cfg := &Config{
		OpenAIAPIKey: apiKey,
		Host:         getEnv("HOST", "0.0.0.0"),
		OpenAIModel:  getEnv("OPENAI_MODEL", "gpt-4o-mini"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		SystemPrompt: getEnv("SYSTEM_PROMPT", "You are a helpful, harmless, and honest AI assistant."),
	}

	var err error

	cfg.Port, err = getEnvInt("PORT", 8080)
	if err != nil {
		return nil, fmt.Errorf("invalid PORT value: %w", err)
	}

	cfg.MaxTokens, err = getEnvInt("MAX_TOKENS", 2048)
	if err != nil {
		return nil, fmt.Errorf("invalid MAX_TOKENS value: %w", err)
	}

	cfg.Temperature, err = getEnvFloat("TEMPERATURE", 0.7)
	if err != nil {
		return nil, fmt.Errorf("invalid TEMPERATURE value: %w", err)
	}

	cfg.Debug, err = getEnvBool("DEBUG", false)
	if err != nil {
		return nil, fmt.Errorf("invalid DEBUG value: %w", err)
	}

	return cfg, nil
}

// Addr returns the full address string for the HTTP server.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) (int, error) {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal, nil
	}
	return strconv.Atoi(val)
}

func getEnvFloat(key string, defaultVal float64) (float64, error) {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal, nil
	}
	return strconv.ParseFloat(val, 64)
}

func getEnvBool(key string, defaultVal bool) (bool, error) {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal, nil
	}
	return strconv.ParseBool(val)
}
