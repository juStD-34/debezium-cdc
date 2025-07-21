# Change Data Capture (CDC) Service on Kubernetes

## 📖 Overview

This project demonstrates a complete **Change Data Capture (CDC)** implementation using **Debezium** on Kubernetes. CDC is a modern approach to real-time data synchronization that captures database changes as they occur, providing instant data updates without the overhead of traditional polling mechanisms.

> **Inspiration**: This project is inspired by the article [Why Cron Jobs are Dead and CDC is the Killer](https://medium.com/@kanishksinghpujari/why-cron-jobs-are-dead-and-cdc-is-the-killer-a7aad011c98f), which explores how CDC revolutionizes data integration patterns.

## 🎯 Purpose

This implementation serves as a learning platform for:
- **Change Data Capture (CDC)** fundamentals and best practices
- **Docker** containerization and multi-service orchestration
- **Kubernetes** deployment patterns and service mesh configuration
- **Debezium** connector setup and MySQL binlog streaming
- Real-time data streaming architecture

## 🏗️ Architecture

The CDC service consists of four core components deployed as separate pods in Kubernetes:

```
┌─────────────────────────────────────────────────────────────────┐
│                        Kubernetes Cluster                       │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                    debezium namespace                       │ │
│  │                                                             │ │
│  │  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐  │ │
│  │  │   ZooKeeper  │    │    Kafka     │    │Kafka Connect │  │ │
│  │  │              │◄───┤              │◄───┤   (Debezium) │  │ │
│  │  │   Port: 2181 │    │  Port: 9092  │    │  Port: 8083  │  │ │
│  │  └──────────────┘    └──────────────┘    └──────────────┘  │ │
│  │         │                     │                     │      │ │
│  │         │                     │                     │      │ │
│  │         │                     │                     │      │ │
│  │  ┌──────▼─────────────────────▼─────────────────────▼────┐ │ │
│  │  │                   MySQL                               │ │ │
│  │  │                                                       │ │ │
│  │  │  • Binary Logging (ROW format)                       │ │ │
│  │  │  • GTID Mode: ON                                      │ │ │
│  │  │  • Sample Data: customers, orders                     │ │ │
│  │  │  • Port: 3306                                         │ │ │
│  │  └───────────────────────────────────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

### Component Breakdown

| Component | Purpose | Key Features |
|-----------|---------|--------------|
| **MySQL** | Source Database | Binary logging, GTID mode, sample inventory data |
| **ZooKeeper** | Coordination Service | Kafka cluster coordination and metadata management |
| **Kafka** | Message Broker | Event streaming platform for CDC events |
| **Kafka Connect** | Connector Platform | Debezium MySQL connector for CDC capture |

## 🛠️ Technology Stack

- **Container Runtime**: Docker
- **Orchestration**: Kubernetes (via Minikube)
- **CDC Platform**: Debezium 2.5
- **Message Broker**: Apache Kafka 7.4.0
- **Database**: MySQL 8.0
- **Coordination**: Apache ZooKeeper 7.4.0
- **CLI Tools**: kubectl, docker

## 🚀 Quick Start

### Prerequisites
Ensure you have the following tools installed (i'm using ubuntu22.04):

```bash
# Verify installations
minikube version (v1.36.0)
docker --version (28.2.2)
kubectl version --client ( Client Version: v1.33.2 // Kustomize Version: v5.6.0)
```

### Deployment Steps 
```bash
cd /k8s
```
1. **Start Minikube**: modify base on your configuration
   ```bash
   minikube start --memory=8192 --cpus=4 --disk-size=20g --driver=docker
   ```
2. **Create Namespace** 
   ```bash
   kubectl apply -f namespace.yaml
   ```

3. **Deploy Persistent Volumes**
   ```bash
   kubectl apply -f persistent-volumes.yaml
   ```

4. **Deploy Services in Order**
   ```bash
   # Deploy ZooKeeper first
   kubectl apply -f zookeeper.yaml
   
   # Wait for ZooKeeper to be ready, then deploy Kafka
   kubectl apply -f kafka.yaml
   
   # Deploy MySQL with CDC configuration
   kubectl apply -f mysql.yaml
   
   # Finally, deploy Kafka Connect with Debezium
   kubectl apply -f kafka-connect.yaml
   ```

5. **Verify Deployment**
   ```bash
   kubectl get pods -n debezium
   kubectl get services -n debezium
   ```

## 📊 Debezium Configuration

This setup uses **Debezium MySQL Connector** to capture changes from the MySQL database. The key CDC configurations include:

### MySQL CDC Settings
- **Binary Logging**: `--log-bin=mysql-bin` - Enables change tracking
- **Binlog Format**: `--binlog-format=ROW` - Captures full row changes
- **GTID Mode**: `--gtid-mode=ON` - Global transaction identifier for consistent replication
- **Row Image**: `--binlog-row-image=FULL` - Complete before/after row data

### Database Privileges
The MySQL user has specific CDC permissions:
```sql
GRANT SELECT, RELOAD, SHOW DATABASES, REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'mysqluser'@'%';
```

### Sample Data Schema
- **customers** table: User profile information
- **orders** table: Transaction records with foreign key relationships

## 🔍 Monitoring and Testing

### Check Service Health
```bash
# Pod status
kubectl get pods -n debezium -w

# Service endpoints
kubectl get svc -n debezium

# Kafka Connect status
kubectl port-forward svc/kafka-connect 8083:8083 -n debezium
curl http://localhost:8083/connector-plugins
```

### Database Operations
```bash
# Connect to MySQL
kubectl exec -it deployment/mysql -n debezium -- mysql -u mysqluser -pmysqlpw inventory

# Insert test data to trigger CDC
INSERT INTO customers (first_name, last_name, email) VALUES ('Alice', 'Cooper', 'alice@example.com');
```

### Log trace for kafka event
```bash
# Get detailed info about a topic
kubectl exec -it deployment/kafka -n debezium -- kafka-topics \
  --bootstrap-server localhost:9092 \
  --describe \
  --topic dbserver1.inventory.customers
```

## 🎓 Learning Outcomes

By working with this CDC implementation, you'll gain hands-on experience with:

- **Real-time Data Streaming**: Understanding how CDC eliminates polling delays
- **Event-Driven Architecture**: Building systems that react to data changes instantly
- **Kubernetes Orchestration**: Managing complex multi-service deployments
- **Debezium Ecosystem**: Configuring and operating CDC connectors
- **Container Networking**: Service discovery and inter-pod communication
- **Persistent Storage**: Managing stateful services in Kubernetes

## 🔧 Troubleshooting

### Common Issues

**Pod Startup Issues**
```bash
kubectl describe pod <pod-name> -n debezium
kubectl logs <pod-name> -n debezium
```

**Storage Permissions**
The MySQL deployment includes an init container to fix volume permissions:
```yaml
initContainers:
  - name: fix-permissions
    image: busybox
    command: ['sh', '-c', 'rm -rf /var/lib/mysql/* && chown -R 999:999 /var/lib/mysql']
```

**Network Connectivity**
Verify service-to-service communication:
```bash
kubectl exec -it deployment/kafka-connect -n debezium -- nc -zv kafka 9092
```

## 📚 Next Steps

- Configure Debezium MySQL connector via REST API
- Set up Kafka consumer applications to process CDC events
- Implement data transformation using Kafka Streams
- Add monitoring with Prometheus and Grafana
- Explore multi-source CDC with additional databases

## 🤝 Contributing

This is a learning project! Feel free to:
- Experiment with different Debezium connectors
- Add monitoring and alerting capabilities
- Implement data validation and transformation
- Document your learnings and improvements

---

**Happy Learning!** 🎉 Dive into the world of Change Data Capture and modern data streaming architectures.