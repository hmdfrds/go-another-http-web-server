package admin

import (
	"bufio"
	"fmt"
	"go-another-http-web-server/logger"
	"net/http"
	"os"
	"time"
)

// AdminInterface represents the admin web interface.
type AdminInterface struct {
	Host      string
	AdminPort int
	Username  string
	Password  string
	Log       *logger.Logger
}

// NewAdminInterface creates a new AdminInterface with default credentials.
func NewAdminInterface(host string, port int, log *logger.Logger) *AdminInterface {
	return &AdminInterface{
		Host:      host,
		AdminPort: port,
		Username:  "admin",
		Password:  "adminpass",
		Log:       log,
	}
}

// Start launches the admin interface server in a separate goroutine.
func (a *AdminInterface) Start() {
	addr := fmt.Sprintf("%s:%d", a.Host, a.AdminPort)
	server := &http.Server{
		Addr:    addr,
		Handler: a,
	}
	go func() {
		fmt.Printf("Admin interface is runnign on %s\n", addr)
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Admin server error: %v\n", err)
		}
	}()
}

// ServeHTTP handles incoming HTTP requests to the admin interface.
func (a *AdminInterface) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check for valid Basic Authentication.
	user, pass, ok := r.BasicAuth()
	if !ok || user != a.Username || pass != a.Password {
		w.Header().Set("WWW-Authenticate", `Basic realm="Admin Interface"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Generate the amin page HTML.

	html := a.generateAdminPage()

	// Set headers and write response.
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(html)))
	fmt.Fprint(w, html)
}

// generateAdminPage builds the HTML content of the admin page.
func (a *AdminInterface) generateAdminPage() string {
	// Retrieve statistics from the logger.
	totalRequests := a.Log.TotalRequests()
	uptime := int(time.Since(a.Log.StartTime()).Seconds())
	activeConns := a.Log.ActiveConnections()

	// Read the last 10 lines from the log file.
	logLines := readLastLines(a.Log.LogFile(), 10)

	// Build HTML content with embedded CSS and meta refresh.
	html := "<html><head><title>Admin Interface</title>"
	html += "<meta http-equiv='refresh' content='30'>"
	html += `<style>
	body {font-family: Arial, sans-serif; margin: 20px;}
	table {border-collapse: collapse; width:80px;}
	th {background-color: #f2f2f2;}
	</style>`

	html += "</head><body>"
	html += "<h1>Admin Interface</h1>"
	html += fmt.Sprintf("<p><strong>Total Requests:</strong> %d</p>", totalRequests)
	html += fmt.Sprintf("<p><strong>Server Uptime:</strong> %d seconds</p>", uptime)
	html += "<h2>Active Connections</h2>"
	if len(activeConns) > 0 {
		html += "<table><tr>Client IP</th><th>Connection Time</th></tr>"
		for ip, connTime := range activeConns {
			html += fmt.Sprintf("<tr><td>%s</td></tr><tr><td>%s</td></tr>", ip, connTime.Format(time.RFC1123))
		}
		html += "</table>"
	} else {
		html += "<p>No active connections.</p>"
	}

	html += "<h2>Last 10 Log Entries</h2><pre>"
	for _, line := range logLines {
		html += line + "\n"
	}
	html += "</pre></body></html>"
	return html
}

// readLastLines reads the last n lines from a given file.
func readLastLines(filepath string, n int) []string {
	var lines []string
	file, err := os.Open(filepath)
	if err != nil {
		return []string{"Error reading log file."}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	// Return only the last n lines
	if len(lines) <= n {
		return lines
	}
	return lines[len(lines)-n:]
}
