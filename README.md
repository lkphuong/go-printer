# Go Printer

A cross-platform printing service built with Go and Gin framework. This application provides REST API endpoints for managing printers, print configurations, and print jobs.

## Features

- List available printers
- Get and set print configurations
- Submit print jobs
- Clear print cache
- Cross-platform support (macOS, Linux, Windows)
- RESTful API with CORS support

## Installation

### Prerequisites

- Go 1.23 or later
- CUPS (on macOS/Linux) or Windows printing services

### Clone the repository

```bash
git clone https://github.com/lkphuong/go-printer.git
cd go-printer
```

### Install dependencies

```bash
go mod download
```

## Usage

### Running the server

```bash
make run
```

Or directly:

```bash
go run ./cmd/server/main.go
```

The server will start on port 9099.

### Building for different platforms

```bash
# Build for Windows
make build-windows

# Build for macOS (ARM64)
make build-mac
```

## API Endpoints

All endpoints are prefixed with `/api/v1`.

### Printers

- `GET /printers` - Get list of available printers
- `GET /printers/:printer/config` - Get print configuration for a specific printer
- `POST /printers/config` - Configure a printer
- `POST /printers/jobs` - Submit a print job
- `DELETE /printers/cache` - Clear print cache

### Example API Usage

```bash
# Get printers
curl http://localhost:9099/api/v1/printers

# Get config for a printer
curl http://localhost:9099/api/v1/printers/MyPrinter/config

# Submit a print job (example payload)
curl -X POST http://localhost:9099/api/v1/printers/jobs \
  -H "Content-Type: application/json" \
  -d '{"printer": "MyPrinter", "file": "path/to/file.pdf"}'
```

## Configuration

The application uses `config/config.json` for storing printer configurations. The file is automatically created if it doesn't exist.

## Project Structure

```
go-printer/
├── cmd/server/          # Main application entry point
├── internal/
│   ├── app/            # Application setup and initialization
│   ├── handlers/       # HTTP request handlers
│   ├── services/       # Business logic services
│   ├── routers/        # Route definitions
│   ├── dto/            # Data transfer objects
│   └── utils/          # Utility functions
├── config/             # Configuration files
├── uploads/            # Uploaded files directory
├── docs/               # Documentation
├── build/              # Build artifacts
└── README.md
```

## Dependencies

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [gin-contrib/cors](https://github.com/gin-contrib/cors) - CORS middleware

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
