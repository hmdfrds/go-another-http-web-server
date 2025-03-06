package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Logger is a thread-safe logger for the server.
type Logger struct {
	logFile           string
	totalRequests     int
	startTime         time.Time
	activeConnections map[string]time.Time // e.g, map[clientIp]connectionStartTime
	mu                sync.Mutex
}

// NewLogger initializes and returns a new Logger.
func NewLogger(logFile string) *Logger {
	return &Logger{
		logFile:           logFile,
		startTime:         time.Now(),
		activeConnections: make(map[string]time.Time),
	}
}

// Log writes a message to the log file with a timestamp.
func (l *Logger) Log(message string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)

	// Open the log file for appending; create it if it doesn't exist.
	f, err := os.OpenFile(l.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}

	defer f.Close()

	if _, err := f.WriteString(logEntry); err != nil {
		fmt.Println("Error writing log entry:", err)
	}
}

// LogRequest logs an HTTP request with client details.
func (l *Logger) LogRequest(clientIP, requestLine string, responseCode int) {
	l.mu.Lock()
	l.totalRequests++
	l.mu.Unlock()

	message := fmt.Sprintf("REQUEST from %s: '%s' responded with %d", clientIP, requestLine, responseCode)

	l.Log(message)
}

// LogError logs an error message.
func (l *Logger) LogError(errorMessage string) {
	l.Log(fmt.Sprintf("ERROR: %s", errorMessage))
}

// LogStats logs periodic server statictics.
func (l *Logger) LogStats() {
	l.mu.Lock()
	uptime := time.Since(l.startTime).Seconds()
	activeCount := len(l.activeConnections)
	totalRequests := l.totalRequests
	l.mu.Unlock()

	statsMessage := fmt.Sprintf("STATS: Total Requests: %d, Active Connections: %d, Uptime %.0f seconds", totalRequests, activeCount, uptime)

	l.Log(statsMessage)
}

// StartPeriodicStats starts a background goroutine to log stats periodically.
func (l *Logger) StartPeriodicStats(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			l.LogStats()
		}
	}()
}

// TotalRequests getter for totalRequests
func (l *Logger) TotalRequests() int {
	return l.totalRequests
}

// StartTime getter for startTime
func (l *Logger) StartTime() time.Time {
	return l.startTime
}

// ActiveConnections getter for activeConnections
func (l *Logger) ActiveConnections() map[string]time.Time {
	return l.activeConnections
}

// LogFile getter for logFile
func (l *Logger) LogFile() string {
	return l.logFile
}
