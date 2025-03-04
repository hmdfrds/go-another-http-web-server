package handler

import (
	"bufio"
	"fmt"
	"go-another-http-web-server/logger"
	"net"
)

type Config struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	AdminPort    int    `json:"admin_port"`
	DocumentRoot string `json:"document_root"`
	MaxThreads   int    `json:"max_threads"`
	LogFile      string `json:"log_file"`
}

// RequestHandler processes a single HTTP request.
type RequestHandler struct {
	Conn   net.Conn
	Config *Config
	Log    *logger.Logger
}

// NewRequestHandler creates a new RequestHandler.
func NewRequestHandler(conn net.Conn, config *Config, log *logger.Logger) *RequestHandler {
	return &RequestHandler{
		Conn:   conn,
		Config: config,
		Log:    log,
	}
}

// sendResponse is a helper to send a basic HTTP response.
func (h *RequestHandler) sendResponse(statusCode int, statusText, body string) {
	h.sendResponseWithHeaders(statusCode, statusText, body, nil)
}

// sendResponseWithHeaders formats and sends a HTTP response with given headers.
func (h *RequestHandler) sendResponseWithHeaders(statusCode int, statusText, body string, headers map[string]string) {
	writer := bufio.NewWriter(h.Conn)
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusText)
	writer.WriteString(statusLine)
	if headers != nil {
		for k, v := range headers {
			writer.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
		}
	}
	writer.WriteString("\r\n")
	writer.WriteString(body)
	writer.Flush()
}
