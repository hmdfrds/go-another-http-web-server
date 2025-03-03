package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Logger is a thread-safe logger for the server.
type Logger struct {
	LogFile           string
	totalRequests     int
	startTime         time.Time
	activeConnections map[string]time.Time // e.g, map[clientIp]connectionStartTime
	mu                sync.Mutex
}

// NewLogger initializes and returns a new Logger.
func NewLogger(logFile string) *Logger {
	return &Logger{
		LogFile:           logFile,
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
	f, err := os.OpenFile(l.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}

	defer f.Close()

	if _, err := f.WriteString(logEntry); err != nil {
		fmt.Println("Error writing log entry:", err)
	}
}
