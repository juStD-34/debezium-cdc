# Kafka Connector Service

A Go-based HTTP service that provides a simplified REST API wrapper around Kafka Connect's native REST API for managing Debezium connectors.

## 🎯 Purpose & Disclaimer

**Important Note**: This service is primarily built as a **learning exercise** and is **not essential** for production use. You can interact directly with Kafka Connect's REST API at `http://kafka-connect:8083` without this wrapper.

### Why This Service Exists

This project serves as a practical example for:
- **Go Project Architecture**: Demonstrating clean code structure, proper error handling, and HTTP service patterns
- **Containerization Learning**: Understanding Docker best practices, multi-stage builds, and container optimization
- **Kubernetes Integration**: Exploring service-to-service communication in a Kubernetes environment
- **API Design**: Creating RESTful endpoints with proper HTTP status codes and JSON responses

## 🏗️ Service Architecture

```
┌─────────────────┐    ┌──────────────────────┐    ┌─────────────────┐
│   Client Apps   │────│ Kafka Connector Svc  │────│  Kafka Connect  │
│                 │    │   (This Service)     │    │   (Port 8083)   │
└─────────────────┘    └──────────────────────┘    └─────────────────┘
```

This service acts as a simple HTTP proxy/wrapper around Kafka Connect's REST API.

## 🚀 Features

- **Simplified API**: Clean REST endpoints for connector management
- **Error Handling**: Proper HTTP status codes and error messages
- **Health Checks**: Built-in health monitoring endpoints
- **Dockerized**: Production-ready container with multi-stage builds
- **Kubernetes Ready**: Includes deployment manifests and service discovery

## 📋 API Endpoints

### Connector Management
```http
POST   /connectors          # Create a new connector
GET    /connectors          # List all connectors
GET    /connectors/{name}   # Get specific connector details
DELETE /connectors/{name}   # Delete a connector
```

### Health & Monitoring
```http
GET    /health              # Service health check
GET    /ready               # Readiness probe
```

## 🛠️ Technology Stack

- **Language**: Go 1.22.4
- **HTTP Client**: Resty (go-resty/resty/v2)
- **HTTP Server**: Native `net/http`
- **Containerization**: Docker with multi-stage builds
- **Orchestration**: Kubernetes
- **Target Service**: Kafka Connect REST API

## 🔧 Local Development

### Prerequisites
- Go 1.22.4
- Docker
- kubectl (for Kubernetes deployment)
- Running Kafka Connect instance

### Environment Setup

1. **Clone and build**:
```bash
git clone <repository-url>
cd kafka-connector-service
go mod tidy
go build -o bin/server ./cmd/server
```

2. **Run locally**:
```bash
# Set environment variables
export KAFKA_CONNECT_URL=http://localhost:8083
export SERVER_PORT=8080

# Run the service
./bin/server
```

3. **Test the endpoints**:
```bash
# Health check
curl http://localhost:8080/health

# List connectors
curl http://localhost:8080/connectors
```

## 🐳 Docker Usage

### Build Image
```bash
docker build -t kafka-connector-service:latest .
```

### Run Container
```bash
docker run -d \
  --name kafka-connector-service \
  -p 8080:8080 \
  -e KAFKA_CONNECT_URL=http://kafka-connect:8083 \
  kafka-connector-service:latest
```

## ☸️ Kubernetes Deployment

The service is designed to work with your existing Kafka Connect setup.

### Deploy the Service
```bash
# Apply the service deployment
kubectl apply -f k8s/

# Check deployment status
kubectl get pods -n debezium
kubectl get svc -n debezium
```

### Service Configuration
```yaml
apiVersion: v1
kind: Service
metadata:
  name: kafka-connector-service
  namespace: debezium
spec:
  selector:
    app: kafka-connector-service
  ports:
    - port: 8080
      targetPort: 8080
```

## 📝 Example Usage

### Creating a MySQL Debezium Connector
```bash
curl -X POST http://kafka-connector-service:8080/connectors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "mysql-connector",
    "config": {
      "connector.class": "io.debezium.connector.mysql.MySqlConnector",
      "database.hostname": "mysql",
      "database.port": "3306",
      "database.user": "mysqluser",
      "database.password": "mysqlpw",
      "database.server.id": "184054",
      "database.server.name": "mysql-server",
      "database.include.list": "inventory",
      "database.history.kafka.bootstrap.servers": "kafka:9092",
      "database.history.kafka.topic": "dbhistory.inventory"
    }
  }'
```

### Listing All Connectors
```bash
curl http://kafka-connector-service:8080/connectors
```

## 🔍 Why Not Use Kafka Connect Directly?

You absolutely can! Here's the direct equivalent:

```bash
# Direct Kafka Connect API call
curl -X POST http://kafka-connect:8083/connectors \
  -H "Content-Type: application/json" \
  -d '{"name": "mysql-connector", "config": {...}}'
```

**This service adds**:
- Simplified request/response handling
- Custom validation and error messages
- Potential for authentication/authorization layers
- Request logging and monitoring
- Custom business logic (if needed)



## 🚦 Production Considerations

While this is a learning project, for production use consider:
- Authentication and authorization
- Rate limiting
- Request/response logging
- Metrics and monitoring (Prometheus)
- Circuit breaker patterns
- Input validation and sanitization
- TLS/SSL termination

## 📚 References

- [Kafka Connect REST API Documentation](https://docs.confluent.io/platform/current/connect/references/restapi.html)
- [Debezium MySQL Connector](https://debezium.io/documentation/reference/connectors/mysql.html)
- [Go HTTP Server Best Practices](https://golang.org/doc/articles/wiki/)
- [Docker Multi-stage Builds](https://docs.docker.com/develop/dev-best-practices/dockerfile_best-practices/)


**Remember**: This service is a learning tool. In production, consider using Kafka Connect's REST API directly or explore enterprise solutions like Confluent Control Center for connector management.