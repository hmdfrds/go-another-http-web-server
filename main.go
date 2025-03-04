package main

import (
	"encoding/json"
	"fmt"
	"go-another-http-web-server/handler"
	"go-another-http-web-server/logger"
	"net"
	"os"
)

func loadConfig(configFileName string) (*handler.Config, error) {
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
	var config handler.Config
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

// startServer initializes the TCP listener and starts the accept loop.
func startServer(config *handler.Config, log *logger.Logger) {
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.LogError(fmt.Sprintf("Failed to bind server on %s: %v", address, err))
		os.Exit(1)
	}
	fmt.Printf("HTTP Server listening on %s\n", address)
	log.Log("Server started on " + address)

	// Main accept loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.LogError(fmt.Sprintf("Error accepting connection: %v", err))
			continue
		}
		// For each connection, spawn a new goroutine to handle the request.
		go handleConnection(conn, config, log)
	}
}

// handleConnection is a placeholder for processing client connections.
func handleConnection(conn net.Conn, config *handler.Config, log *logger.Logger) {
	defer conn.Close()

	// For testing purposes, simply send back a basic HTTP response.
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"Hello, World!"

	conn.Write([]byte(response))

	// Log the simple interaction.
	log.LogRequest(conn.RemoteAddr().String(), "GET /", 200)
}

func main() {

	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		os.Exit(0)
	}

	fmt.Printf("Configuration loaded successfully:\n%+v\n", config)

	log := logger.NewLogger(config.LogFile)
	log.StartPeriodicStats(60 * 1e9)

	startServer(config, log)
}
