# DataFlow - Real-time Data Orchestration Tool

A powerful real-time data orchestration platform with visual pipeline builder for complex data workflows. Features drag-and-drop interface, automatic data validation, and seamless integration with multiple data sources and destinations.

## Features

- **Visual Pipeline Designer**: Drag-and-drop interface for building data pipelines
- **Real-time Processing**: Monitor and process data streams in real-time
- **Multiple Connectors**: Support for databases, APIs, files, and more
- **Data Validation**: Automatic validation and error handling
- **Pipeline Versioning**: Version control and rollback capabilities
- **RESTful API**: Programmatic access to all functionality
- **Web Dashboard**: Comprehensive pipeline management interface

## Quick Start

### Prerequisites

- Go 1.19+
- PostgreSQL 12+
- Node.js 16+ (for frontend development)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/your-org/dataflow.git
cd dataflow
```

2. Install Go dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
export PORT=8080
export DATABASE_URL=postgres://user:password@localhost/dataflow
```

4. Run the application:
```bash
go run main.go
```

5. Open your browser and navigate to `http://localhost:8080`

## API Endpoints

### Pipelines

- `GET /api/v1/pipelines` - List all pipelines
- `POST /api/v1/pipelines` - Create a new pipeline
- `GET /api/v1/pipelines/:id` - Get pipeline details
- `POST /api/v1/pipelines/:id/execute` - Execute a pipeline

### Health Check

- `GET /api/v1/health` - Service health status

## Pipeline Structure

A pipeline consists of:

- **Nodes**: Processing steps (source, transform, sink)
- **Connections**: Data flow between nodes
- **Configuration**: Node-specific settings
- **Metadata**: Pipeline name, status, timestamps

### Example Pipeline JSON

```json
{
  "name": "Data Processing Pipeline",
  "nodes": [
    {
      "id": "source_1",
      "type": "source",
      "name": "Database Source",
      "config": {
        "connection_string": "postgres://...",
        "query": "SELECT * FROM users"
      },
      "position": { "x": 100, "y": 100 }
    },
    {
      "id": "transform_1",
      "type": "transform",
      "name": "Data Cleaner",
      "config": {
        "operations": ["trim", "lowercase"]
      },
      "position": { "x": 300, "y": 100 }
    }
  ],
  "connections": [
    {
      "id": "conn_1",
      "source": "source_1",
      "target": "transform_1"
    }
  ]
}
```

## Development

### Project Structure

```
├── main.go                 # Application entry point
├── internal/
│   ├── api/               # HTTP handlers and routes
│   │   └── handlers.go
│   └── pipeline/          # Core pipeline engine
│       └── engine.go
├── web/
│   ├── static/           # Frontend assets
│   │   └── dashboard.js
│   └── templates/        # HTML templates
└── README.md
```

### Building

```bash
# Build for current platform
go build -o dataflow main.go

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o dataflow-linux main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o dataflow.exe main.go
```

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Support

For questions and support, please open an issue on GitHub or contact the development team.