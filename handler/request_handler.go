package handler

import (
	"bufio"
	"fmt"
	"go-another-http-web-server/logger"
	"go-another-http-web-server/utils"
	"io"
	"mime"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

// serverDirectory gtenerates an HTML directory listing.
func (h *RequestHandler) serveDirectory(dirPath, requestLine string) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		h.sendResponse(500, "Internal Server Error", "<html><body><h1>500 Internal Server Error</h1></body></html>")
		h.Log.LogError("Error reading directory: " + err.Error())
		return
	}

	html := "<html><head><title><Directory Listing</title></head><body>"
	html += fmt.Sprintf("<h1>Directory listing for %s</h><ul>", dirPath)
	for _, file := range files {
		name := file.Name()
		if file.IsDir() {
			name += "/"
		}
		fileInfo, err := file.Info()
		if err != nil {
			h.sendResponse(500, "Internal Server Error", "<html><body><h1>500 Internal Server Error</h1></body></html>")
			h.Log.LogError("Error while reading file info" + err.Error())
			return
		}
		modTime := fileInfo.ModTime().Format(time.DateTime)
		html += fmt.Sprintf("<li>%s - Last Modified: %s</li>", name, modTime)
	}
	html += "</ul></body></html>"

	headers := map[string]string{
		"Content-Type":   "text/html",
		"Content-Length": strconv.Itoa(len(html)),
		"Date":           utils.HTTPDateFormat(time.Now()),
		"Server":         "GoHTTP/1.0",
		"Connection":     "close",
	}
	h.sendResponseWithHeaders(200, "OK", html, headers)
	h.Log.LogRequest(h.Conn.RemoteAddr().String(), requestLine, 200)
}

// serveFile reads a file and sends it as an HTTP response.
func (h *RequestHandler) serveFile(filePath, method, requestLine string) {
	f, err := os.Open(filePath)
	if err != nil {
		h.sendResponse(500, "Internal Server Error", "<html><body><h1>500 Internal Server Error</h1></body></html>")
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {

		h.sendResponse(500, "Internal Server Error", "<html><body><h1>500 Internal Server Error</h1></body><html>")
		h.Log.LogError("Error reading file" + err.Error())
		return
	}

	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Prepare headers.
	headers := map[string]string{
		"Content-Type":   mimeType,
		"Content-Length": strconv.Itoa(len(data)),
		"Date":           utils.HTTPDateFormat(time.Now()),
		"Server":         "GoHTTP/1.0",
		"Connection":     "close",
	}

	// For HEAD requests, send headers only.
	if method == "HEAD" {
		h.sendResponseWithHeaders(200, "OK", "", headers)
	} else {
		h.sendResponseWithHeaders(200, "OK", string(data), headers)
	}
	h.Log.LogRequest(h.Conn.RemoteAddr().String(), requestLine, 200)
}

func (h *RequestHandler) Handle() {
	defer h.Conn.Close()
	reader := bufio.NewReader(h.Conn)

	// Read requets line
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		h.Log.LogError("Error reading request line: " + err.Error())
		return
	}
	requestLine = strings.TrimSpace(requestLine)
	parts := strings.Split(requestLine, " ")
	if len(parts) < 3 {
		h.sendResponse(400, "Bad Request", "Invalid request line")
		return
	}
	method, rawPath := parts[0], parts[1]

	// Read headers until an empty line is reached.
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			h.Log.LogError("Error reading headers: " + err.Error())
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		headerParts := strings.SplitN(line, ":", 2)
		if len(headerParts) == 2 {
			headers[strings.TrimSpace(headerParts[0])] = strings.TrimSpace(headerParts[1])
		}
	}

	// Parse URL to ensure proper handling
	parsedUrl, err := url.Parse(rawPath)
	if err != nil {
		h.sendResponse(400, "Bad Request", "Invalid URL")
		return
	}
	path := parsedUrl.Path

	// Resolve the requested file path safely.
	filePath, err := utils.SafePath(h.Config.DocumentRoot, path)
	if err != nil {
		h.sendResponse(403, "Forbidden", "Access denied")
		return
	}

	// Check if the file or directory exists.
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		h.sendResponse(404, "Not Found", "<html><body><h1>404 Not Found</h1></body></html>")
		h.Log.LogError("Error stating file: " + err.Error())
		return
	}

	// If the path is a directory, check for index.html or generate a listing.
	if info.IsDir() {
		indexPath := filepath.Join(filePath, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			filePath = indexPath
		} else {
			h.serveDirectory(filePath, requestLine)
			return
		}
	}

	// Only allow GET and HEAD methods.
	if method != "GET" && method != "HEAD" {
		h.sendResponse(405, "Method Not Allowed", "<html><body><h1>405 Method Not Allowed</h1></body></html>")
		h.Log.LogRequest(h.Conn.RemoteAddr().String(), requestLine, 405)
		return
	}

	// Serve the file
	h.serveFile(filePath, method, requestLine)
}
