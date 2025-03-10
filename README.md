# Go Another HTTP Web Server

Another HTTP web server written in Go.

## Features

- **Static File Serving:** Serves files from `www/`.
- **Directory Listing:** Automatically lists directory contents if no index.html exists.
- **HTTP Methods:** Supports GET and HEAD.
- **Admin Interface:** Accessible at `/` on the admin port with HTTP Basic Auth (default: `admin/adminpass`).
- **Thread-safe Logging:** Logs requests, errors, and server stats.

## Project Structure

```text
go-another-http-web-server/ 
├── admin/ 
│ └── admin_interface.go 
├── handler/ 
│ └── request_handler.go 
├── logger/ 
│ └── logger.go 
├── utils/ 
│ └── utils.go 
├── config.json 
├── main.go 
├── www/ 
│ └── index.html 
└── README.md
```

## Setup & Run

1. **Clone the repository:**

    ```bash
    git clone https://github.com/hmdfrds/go-another-http-web-server
    cd go-another-http-web-server
    ```

2. **Build and run:**

    ```bash
    go build -o web-server.exe
    ./web-server
    ```

3. **Access the server:**
    - Main site: <http://localhost:8080/>
    - Admin interface: <http://localhost:8081/> (Credentials: admin / adminpass)

## Configuration

Adjust settings in config.json as needed.

## License

MIT [LICENSE](./LICENSE).
