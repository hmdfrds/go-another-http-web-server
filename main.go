package main

import (
	"encoding/json"
	"fmt"
	"os"

	"go-another-http-web-server.git/logger"
)

type Config struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	AdminPort    int    `json:"admin_port"`
	DocumentRoot string `json:"document_root"`
	MaxThreads   int    `json:"max_threads"`
	LogFile      string `json:"log_file"`
}

func loadConfig(configFileName string) (*Config, error) {
	// Check if the file exists
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file %s not found", configFileName)
	}

	// Read the file contents
	data, err := os.ReadFile(configFileName)
	if err != nil {
		return nil, fmt.Errorf("error parsing config JSON: %v", err)
	}

	// Validate required fields
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config JSON: %v", err)
	}

	// Validate required fields
	if config.Host == "" {
		return nil, fmt.Errorf("host is required in config")
	}
	if config.Port == 0 {
		return nil, fmt.Errorf("port is required and must be non-zero")
	}
	if config.AdminPort == 0 {
		return nil, fmt.Errorf("admin_port is required and must be non-zero")
	}
	if config.DocumentRoot == "" {
		return nil, fmt.Errorf("document_root is required")
	}
	if config.MaxThreads == 0 {
		return nil, fmt.Errorf("max_threads is required and must be non-zero")
	}
	if config.LogFile == "" {
		return nil, fmt.Errorf("log_file is required")
	}

	return &config, nil
}

func main() {

	log := logger.NewLogger("server.log")

	log.LogRequest("127.0.0.1", "GET /index.html HTTP/1.1", 200)
	log.LogError("Test error message")

	os.Exit(0)
}
