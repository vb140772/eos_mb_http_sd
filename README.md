# MinIO Prometheus HTTP Service Discovery

A Go-based service that implements Prometheus HTTP Service Discovery for MinIO (EOS) v3 metrics. This service dynamically discovers MinIO buckets and generates scrape configurations for their metrics endpoints.

## 📚 **Table of Contents**

1. [Overview](#-overview)
2. [Features](#-features)
3. [Architecture](#-architecture)
4. [Configuration](#-configuration)
5. [Bucket Wildcard Patterns](#-bucket-wildcard-patterns)
6. [MinIO v3 Metrics Authentication](#-minio-v3-metrics-authentication)
7. [API Endpoints](#-api-endpoints)
8. [Prometheus Integration](#-prometheus-integration)
9. [Deployment](#-deployment)
10. [Usage Examples](#-usage-examples)
11. [Troubleshooting](#-troubleshooting)
12. [Security Considerations](#-security-considerations)

---

## 🎯 **Overview**

The MinIO Prometheus Service Discovery service is a Go-based HTTP service that implements the Prometheus HTTP Service Discovery protocol. It dynamically discovers MinIO buckets and provides service discovery endpoints for Prometheus to scrape MinIO v3 metrics.

### **What It Does**
- **Dynamic bucket discovery** - Automatically finds all MinIO buckets
- **Prometheus HTTP SD** - Implements the Prometheus service discovery protocol
- **MinIO v3 metrics support** - Provides endpoints for both server and bucket metrics
- **Flexible bucket filtering** - Supports wildcard patterns for bucket inclusion/exclusion


### **Key Benefits**
- ✅ **No manual configuration** - Automatically discovers buckets
- ✅ **Scalable** - Works with any number of buckets
- ✅ **Production ready** - Secure and reliable
- ✅ **Easy integration** - Simple Prometheus configuration
- ✅ **Flexible filtering** - Wildcard-based bucket selection

---

## ✨ **Features**

### **Core Functionality**
- **HTTP Service Discovery** - Prometheus-compatible service discovery endpoints
- **Dynamic Bucket Scanning** - Real-time bucket discovery from MinIO
- **MinIO v3 Metrics Support** - Both server and bucket metrics endpoints
- **Wildcard Bucket Filtering** - Pattern-based bucket inclusion/exclusion


### **Advanced Features**
- **Configurable Scrape Intervals** - Customizable Prometheus scrape timing
- **SSL/TLS Support** - Secure connections to MinIO
- **Health Check Endpoints** - Service monitoring and readiness
- **Structured Logging** - Comprehensive logging with logrus
- **Docker Support** - Containerized deployment ready

---

## 🏗️ **Architecture**

### **Service Components**

```
┌─────────────────────────────────────────────────────────────┐
│                    Prometheus                               │
└─────────────────────┬───────────────────────────────────────┘
                      │ HTTP SD Request
                      ▼
┌─────────────────────────────────────────────────────────────┐
│              MinIO Prometheus SD Service                   │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │  HTTP Router    │  │  MinIO Client   │  │  Config    │ │
│  │  (gorilla/mux)  │  │  (minio-go/v7)  │  │  Manager   │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────┬───────────────────────────────────────┘
                      │ MinIO API Calls
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                    MinIO Server                            │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   Bucket List   │  │  Server Metrics │  │ Bucket     │ │
│  │     API         │  │     v3          │  │ Metrics v3 │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### **Data Flow**
1. **Prometheus** requests service discovery data
2. **Service** queries MinIO for bucket list
3. **Service** applies wildcard filtering to buckets
4. **Service** generates Prometheus targets with labels
5. **Prometheus** receives targets and scrapes metrics

---

## ⚙️ **Configuration**

### **Environment Variables**

The service is configured via environment variables:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `MINIO_ENDPOINT` | MinIO server endpoint (host:port) | `localhost:9000` | No |
| `MINIO_ACCESS_KEY` | MinIO access key | `minioadmin` | No |
| `MINIO_SECRET_KEY` | MinIO secret key | `minioadmin` | No |
| `MINIO_USE_SSL` | Whether to use SSL/TLS | `false` | No |
| `LISTEN_ADDR` | Address to listen on | `:8080` | No |
| `SCRAPE_INTERVAL` | Prometheus scrape interval | `15s` | No |
| `METRICS_PATH` | Custom metrics path | `/minio/metrics/v3` | No |
| `BUCKET_PATTERN` | Wildcard pattern for bucket inclusion | `*` | No |
| `BUCKET_EXCLUDE_PATTERN` | Wildcard pattern for bucket exclusion | (empty) | No |


### **Configuration Examples**

#### **Development Environment**
```bash
export MINIO_ENDPOINT="localhost:9000"
export MINIO_ACCESS_KEY="minioadmin"
export MINIO_SECRET_KEY="minioadmin"
export MINIO_USE_SSL="false"
```

#### **Production Environment**
```bash
export MINIO_ENDPOINT="minio-prod.company.com:9000"
export MINIO_ACCESS_KEY="prod-access-key"
export MINIO_SECRET_KEY="prod-secret-key"
export MINIO_USE_SSL="true"
export BUCKET_PATTERN="prod-*"
export BUCKET_EXCLUDE_PATTERN="*backup*"
```

---

## 🌟 **Bucket Wildcard Patterns**

### **Overview**

The service supports **wildcard-based bucket filtering** using patterns similar to shell globbing. This allows you to:

- **Include specific buckets** based on naming patterns
- **Exclude unwanted buckets** from monitoring
- **Create environment-specific** bucket selections
- **Scale monitoring** without manual configuration

### **Pattern Syntax**

| Pattern | Description | Example |
|---------|-------------|---------|
| `*` | Match any sequence of characters | `*` matches all buckets |
| `?` | Match any single character | `test?` matches `test1`, `test2` |
| `[abc]` | Match any character in the set | `[abc]*` matches `a-data`, `b-backup` |
| `[!abc]` | Match any character NOT in the set | `[!abc]*` excludes `a-`, `b-`, `c-` buckets |

### **Configuration**

#### **Include Pattern (BUCKET_PATTERN)**
```bash
# Include all buckets (default)
export BUCKET_PATTERN="*"

# Include only production buckets
export BUCKET_PATTERN="prod-*"

# Include buckets containing "data"
export BUCKET_PATTERN="*data*"

# Include specific bucket types
export BUCKET_PATTERN="user-*,system-*,log-*"
```

#### **Exclude Pattern (BUCKET_EXCLUDE_PATTERN)**
```bash
# Exclude no buckets (default)
export BUCKET_EXCLUDE_PATTERN=""

# Exclude temporary buckets
export BUCKET_EXCLUDE_PATTERN="temp-*"

# Exclude backup buckets
export BUCKET_EXCLUDE_PATTERN="*backup*"

# Exclude multiple patterns
export BUCKET_EXCLUDE_PATTERN="temp-*,*backup*,archive-*"
```

### **Use Cases**

#### **Development Environment**
```bash
export BUCKET_PATTERN="dev-*"
export BUCKET_EXCLUDE_PATTERN="*test*"
# Result: Only buckets starting with "dev-" and not containing "test"
```

#### **Production Environment**
```bash
export BUCKET_PATTERN="prod-*"
export BUCKET_EXCLUDE_PATTERN="*backup*,*archive*"
# Result: Only production buckets, excluding backups and archives
```

#### **Multi-Environment Setup**
```bash
# Development
export BUCKET_PATTERN="dev-*"
export BUCKET_EXCLUDE_PATTERN=""

# Staging
export BUCKET_PATTERN="staging-*"
export BUCKET_EXCLUDE_PATTERN="*temp*"

# Production
export BUCKET_PATTERN="prod-*"
export BUCKET_EXCLUDE_PATTERN="*backup*,*archive*"
```

### **How It Works**

1. **Bucket Discovery**: Service queries MinIO for all buckets
2. **Pattern Matching**: Applies include/exclude patterns using wildcard matching
3. **Target Generation**: Creates Prometheus targets only for matching buckets
4. **Dynamic Updates**: Re-evaluates patterns on each service discovery request

---



---

## 🌐 **API Endpoints**

### **Service Discovery Endpoint**

#### **`GET /sd?job={job_name}`**

Returns Prometheus service discovery targets for the specified job.

**Parameters:**
- `job` (required): The job name to discover targets for

**Supported Jobs:**
- `minio-server`: MinIO server metrics
- `minio-buckets`: MinIO bucket metrics

**Example Request:**
```bash
curl "http://localhost:8080/sd?job=minio-buckets"
```

**Example Response:**
```json
[
  {
    "targets": ["minio-server:9000"],
    "labels": {
      "__metrics_path__": "/minio/metrics/v3/bucket/api/mybucket",
      "__scheme__": "http",
      "instance": "minio-server:9000",
      "job": "minio-buckets",
      "sd_bucket": "mybucket",
      "sd_bucket_creation": "2024-01-15T10:30:00Z"
    }
  }
]
```

### **Scrape Configs Endpoint**

#### **`GET /scrape_configs`**

Returns all available Prometheus scrape configurations.

**Example Request:**
```bash
curl "http://localhost:8080/scrape_configs"
```

**Example Response:**
```json
[
  {
    "job_name": "minio-server",
    "scrape_interval": "15s",
    "scrape_timeout": "10s",
    "metrics_path": "/minio/metrics/v3",
    "scheme": "http",
    "static_configs": [
      {
        "targets": ["minio-server:9000"],
        "labels": {
          "instance": "minio-server:9000",
          "job": "minio-server"
        }
      }
    ]
  },
  {
    "job_name": "minio-buckets",
    "scrape_interval": "15s",
    "scrape_timeout": "10s",
    "metrics_path": "/minio/metrics/v3/bucket/api",
    "scheme": "http",
    "static_configs": [
      {
        "targets": ["minio-server:9000"],
        "labels": {
          "instance": "minio-server:9000",
          "job": "minio-buckets",
          "bucket_pattern": "*",

        }
      }
    ]
  }
]
```

### **Health Check Endpoint**

#### **`GET /health`**

Returns service health status.

**Example Request:**
```bash
curl "http://localhost:8080/health"
```

**Example Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## 📊 **Prometheus Integration**

### **Overview**

The service provides **Prometheus HTTP Service Discovery** endpoints that allow Prometheus to dynamically discover MinIO metrics targets. This eliminates the need for static Prometheus configuration and enables automatic scaling.

### **Service Discovery Configuration**

#### **MinIO Server Metrics**

```yaml
scrape_configs:
  - job_name: 'minio-server'
    http_sd_configs:
      - url: 'http://localhost:8080/sd?job=minio-server'
        refresh_interval: 30s
    scrape_interval: 15s
    scrape_timeout: 10s
    metrics_path: /minio/metrics/v3
    scheme: https
    relabel_configs:
      - source_labels: [__address__]
        target_label: instance
      - source_labels: [__metrics_path__]
        target_label: __metrics_path__
      - source_labels: [__scheme__]
        target_label: __scheme__
```

#### **MinIO Bucket Metrics**

```yaml
scrape_configs:
  - job_name: 'minio-buckets'
    http_sd_configs:
      - url: 'http://localhost:8080/sd?job=minio-buckets'
        refresh_interval: 30s
    scrape_interval: 15s
    scrape_timeout: 10s
    metrics_path: /minio/metrics/v3/bucket/api
    scheme: https
    relabel_configs:
      - source_labels: [__address__]
        target_label: instance
      - source_labels: [__metrics_path__]
        target_label: __metrics_path__
      - source_labels: [__scheme__]
        target_label: __scheme__
      - source_labels: [bucket]
        target_label: bucket
      - source_labels: [bucket_creation]
        target_label: bucket_creation
```

### **Key Benefits**

- **Dynamic Discovery**: Automatically finds new buckets
- **No Manual Configuration**: Prometheus discovers targets automatically
- **Scalable**: Works with any number of buckets

- **Flexible**: Supports wildcard-based bucket filtering

---

## 🚀 **Deployment**

### **Local Development**

#### **Prerequisites**
- Go 1.21 or later
- MinIO server running
- Valid MinIO access credentials

#### **Command Line Arguments**

The service supports command line arguments that take precedence over environment variables:

```bash
# Show help and available options
minio-prometheus-sd -help

# Override specific settings
minio-prometheus-sd -minio-endpoint=minio:9000 -minio-access-key=mykey

# Use different port and bucket pattern
minio-prometheus-sd -listen-addr=:9090 -bucket-pattern="prod-*"

# Enable SSL and set custom scrape interval
minio-prometheus-sd -minio-use-ssl -scrape-interval=30s
```

**Available Flags:**
- `-minio-endpoint`: MinIO server endpoint
- `-minio-access-key`: MinIO access key
- `-minio-secret-key`: MinIO secret key
- `-minio-use-ssl`: Use SSL for MinIO connection
- `-listen-addr`: Address to listen on
- `-scrape-interval`: Scrape interval
- `-metrics-path`: Metrics path
- `-bucket-pattern`: Wildcard pattern for bucket inclusion
- `-bucket-exclude-pattern`: Wildcard pattern for bucket exclusion
- `-help`: Show help information


#### **Setup Steps**

1. **Clone and setup**:
   ```bash
   git clone <repository>
   cd eos_mb_http_sd
   go mod tidy
   ```

2. **Configure environment**:
   ```bash
   export MINIO_ENDPOINT="localhost:9000"
   export MINIO_ACCESS_KEY="minioadmin"
   export MINIO_SECRET_KEY="minioadmin"

   ```

3. **Run the service**:
   ```bash
   go run main.go
   ```

### **Using Docker**

1. **Build the image**:
   ```bash
   docker build -t minio-prometheus-sd .
   ```

2. **Run the container**:
   ```bash
   # Using environment variables
   docker run -p 8080:8080 \
     -e MINIO_ENDPOINT="your-minio-server:9000" \
     -e MINIO_ACCESS_KEY="your-access-key" \
     -e MINIO_SECRET_KEY="your-secret-key" \
     -e MINIO_USE_SSL="true" \
     minio-prometheus-sd

   # Using command line arguments
   docker run -p 8080:8080 minio-prometheus-sd \
     -minio-endpoint="your-minio-server:9000" \
     -minio-access-key="your-access-key" \
     -minio-secret-key="your-secret-key" \
     -minio-use-ssl
   ```

### **Docker Compose**

#### **Complete Stack**
```yaml
version: '3.8'

services:
  minio:
    image: minio/minio:latest
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    command: server /data --console-address ":9001"
    volumes:
      - minio_data:/data

  minio-prometheus-sd:
    build: .
    ports:
      - "8080:8080"
    # Option 1: Using environment variables
    environment:
      - MINIO_ENDPOINT=minio:9000
      - MINIO_ACCESS_KEY=minioadmin
      - MINIO_SECRET_KEY=minioadmin
      - MINIO_USE_SSL=false
      - LISTEN_ADDR=:8080
      - SCRAPE_INTERVAL=15s
      - BUCKET_PATTERN=*
      - BUCKET_EXCLUDE_PATTERN=
    # Option 2: Using command line arguments (override environment variables)
    command: >
      -minio-endpoint=minio:9000
      -minio-access-key=minioadmin
      -minio-secret-key=minioadmin
      -minio-use-ssl=false
      -listen-addr=:8080
      -scrape-interval=15s
      -bucket-pattern=*
      -bucket-exclude-pattern=
    depends_on:
      minio:
        condition: service_started

  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9091:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
      - '--storage.tsdb.retention.time=200h'
      - '--web.enable-lifecycle'

volumes:
  minio_data:
  prometheus_data:
```

### **Kubernetes Deployment**

#### **Deployment YAML**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: minio-prometheus-sd
  labels:
    app: minio-prometheus-sd
spec:
  replicas: 1
  selector:
    matchLabels:
      app: minio-prometheus-sd
  template:
    metadata:
      labels:
        app: minio-prometheus-sd
    spec:
      containers:
      - name: minio-prometheus-sd
        image: minio-prometheus-sd:latest
        ports:
        - containerPort: 8080
        env:
        - name: MINIO_ENDPOINT
          value: "minio-service:9000"
        - name: MINIO_ACCESS_KEY
          valueFrom:
            secretKeyRef:
              name: minio-secrets
              key: access-key
        - name: MINIO_SECRET_KEY
          valueFrom:
            secretKeyRef:
              name: minio-secrets
              key: secret-key

        - name: BUCKET_PATTERN
          value: "prod-*"
        - name: BUCKET_EXCLUDE_PATTERN
          value: "*backup*"
```

---

## 📋 **Usage Examples**

### **Basic Usage**

#### **Start Service**
```bash
# Using environment variables (default)
minio-prometheus-sd

# Using command line arguments
minio-prometheus-sd -minio-endpoint=minio:9000 -minio-access-key=mykey

# Mix of command line and environment variables
export MINIO_ENDPOINT=minio:9000
minio-prometheus-sd -minio-access-key=mykey -bucket-pattern="prod-*"
```

#### **Test Endpoints**
```bash
# Health check
curl "http://localhost:8080/health"

# Service discovery for buckets
curl "http://localhost:8080/sd?job=minio-buckets"

# All scrape configs
curl "http://localhost:8080/scrape_configs"
```

### **Advanced Configuration**

#### **Production Setup**
```bash
export MINIO_ENDPOINT="minio-prod.company.com:9000"
export MINIO_ACCESS_KEY="prod-access-key"
export MINIO_SECRET_KEY="prod-secret-key"
export MINIO_USE_SSL="true"

export BUCKET_PATTERN="prod-*"
export BUCKET_EXCLUDE_PATTERN="*backup*,*archive*"
export SCRAPE_INTERVAL="30s"
```

#### **Development Setup**
```bash
export MINIO_ENDPOINT="localhost:9000"
export MINIO_ACCESS_KEY="minioadmin"
export MINIO_SECRET_KEY="minioadmin"
export MINIO_USE_SSL="false"

export BUCKET_PATTERN="dev-*"
export BUCKET_EXCLUDE_PATTERN="*test*"
```

### **Testing Scenarios**

#### **Test Bucket Filtering**
```bash
# Create test buckets
mc mb local/test-bucket
mc mb local/prod-data
mc mb local/dev-backup

# Test different patterns
export BUCKET_PATTERN="*"
export BUCKET_EXCLUDE_PATTERN=""
# Result: All buckets

export BUCKET_PATTERN="prod-*"
export BUCKET_EXCLUDE_PATTERN=""
# Result: Only prod-data

export BUCKET_PATTERN="*"
export BUCKET_EXCLUDE_PATTERN="*backup*"
# Result: test-bucket, prod-data (excludes dev-backup)
```

---

## 🚨 **Troubleshooting**

### **Common Issues**

#### **Issue: Service Won't Start**

**Symptoms:**
- Service exits immediately
- Connection refused errors
- Configuration errors

**Solutions:**
1. **Check environment variables**:
   ```bash
   echo $MINIO_ENDPOINT
   ```

2. **Verify MinIO connectivity**:
   ```bash
   curl -I "http://localhost:9000/minio/health/live"
   ```

3. **Check MinIO credentials**:
   ```bash
   mc config host add local http://localhost:9000 minioadmin minioadmin
   mc ls local
   ```

#### **Issue: No Buckets Discovered**

**Symptoms:**
- Empty service discovery response
- No targets in Prometheus
- Bucket-related errors

**Solutions:**
1. **Check bucket patterns**:
   ```bash
   echo $BUCKET_PATTERN
   echo $BUCKET_EXCLUDE_PATTERN
   ```

2. **Verify bucket access**:
   ```bash
   mc ls local
   ```

3. **Test pattern matching**:
   ```bash
   # Test with simple pattern
   export BUCKET_PATTERN="*"
   export BUCKET_EXCLUDE_PATTERN=""
   ```

#### **Issue: Authentication Failures**

**Symptoms:**
- 401 Unauthorized errors

- Prometheus scraping failures

**Solutions:**






### **Debug Mode**

#### **Enable Verbose Logging**
```bash
export LOG_LEVEL=debug
go run main.go
```

#### **Check Service Discovery**
```bash
# Get all configurations
curl "http://localhost:8080/scrape_configs" | jq

# Check specific job
curl "http://localhost:8080/sd?job=minio-buckets" | jq
```

#### **Verify Prometheus Configuration**
```bash
# Check Prometheus config
curl "http://localhost:9090/api/v1/status/config" | jq

# Check targets
curl "http://localhost:9090/api/v1/targets" | jq
```

### **Performance Issues**

#### **Slow Service Discovery**
**Causes:**
- Large number of buckets
- Network latency to MinIO
- Complex wildcard patterns

**Solutions:**
1. **Optimize bucket patterns** - Use specific patterns instead of `*`
2. **Increase refresh interval** - Reduce Prometheus polling frequency
3. **Network optimization** - Ensure low latency to MinIO

#### **High Memory Usage**
**Causes:**
- Large bucket lists
- Inefficient pattern matching
- Memory leaks

**Solutions:**
1. **Limit bucket scope** - Use restrictive patterns
2. **Monitor memory usage** - Set resource limits in containers
3. **Regular restarts** - Implement health checks and restarts

---

## 🔒 **Security Considerations**

### **Authentication & Authorization**

#### **MinIO Access Control**
- **Use dedicated credentials** for monitoring
- **Limit permissions** to read-only access
- **Rotate credentials** regularly
- **Monitor access patterns** for anomalies



### **Network Security**

#### **Transport Security**
- **Use HTTPS/TLS** for production deployments
- **Network isolation** - Restrict access to metrics endpoints
- **Firewall rules** - Limit access to necessary ports
- **VPN access** - Secure remote access

#### **Service Security**
- **Bind to localhost** in development
- **Use reverse proxy** for production
- **Implement rate limiting** to prevent abuse
- **Monitor access logs** for suspicious activity

### **Data Protection**

#### **Metrics Data**
- **Minimize data collection** - Only collect necessary metrics
- **Data retention** - Implement appropriate retention policies
- **Access logging** - Track who accesses metrics
- **Data encryption** - Encrypt sensitive metrics data

#### **Configuration Security**
- **Secrets management** - Use proper secrets management tools
- **Configuration validation** - Validate all configuration inputs
- **Access control** - Limit who can modify configuration
- **Audit logging** - Log configuration changes

### **Best Practices**

1. **Principle of Least Privilege**
   - Use minimal required permissions
   - Limit access to necessary resources
   - Regular permission reviews

2. **Defense in Depth**
   - Multiple security layers
   - Network and application security
   - Monitoring and alerting

3. **Regular Security Updates**
   - Keep dependencies updated
   - Monitor security advisories
   - Regular security assessments

4. **Incident Response**
   - Security incident procedures
   - Monitoring and alerting
   - Response team preparation

---

## 🎉 **Summary**

The MinIO Prometheus Service Discovery service provides a **comprehensive solution** for monitoring MinIO v3 metrics with Prometheus. Key features include:

### **✅ Core Capabilities**
- **Dynamic bucket discovery** - Automatic bucket detection
- **Prometheus HTTP SD** - Standard service discovery protocol
- **MinIO v3 support** - Full metrics endpoint coverage
- **Wildcard filtering** - Flexible bucket selection


### **✅ Production Features**
- **Scalable architecture** - Handles any number of buckets
- **Flexible configuration** - Environment-based setup
- **Security focused** - Proper authentication and authorization
- **Monitoring ready** - Health checks and logging
- **Container support** - Docker and Kubernetes ready

### **✅ Deployment Options**
- **Local development** - Simple Go-based setup
- **Docker containers** - Containerized deployment
- **Docker Compose** - Complete stack orchestration
- **Kubernetes** - Production orchestration platform

### **🚀 Getting Started**

1. **Configure environment**:
   ```bash
   export MINIO_ENDPOINT="your-minio-server:9000"
   ```

3. **Start the service**:
   ```bash
   go run main.go
   ```

4. **Configure Prometheus**:
   ```yaml
   http_sd_configs:
     - url: 'http://localhost:8080/sd?job=minio-buckets'
   ```

The service is now ready to provide **dynamic, scalable, and secure** MinIO v3 metrics discovery for your Prometheus monitoring stack! 🎯

---

## 📚 **Additional Resources**

- **MinIO Documentation**: [https://min.io/docs](https://min.io/docs)
- **Prometheus HTTP SD**: [https://prometheus.io/docs/prometheus/latest/http_sd/](https://prometheus.io/docs/prometheus/latest/http_sd/)
- **Go Documentation**: [https://golang.org/doc/](https://golang.org/doc/)
- **Docker Documentation**: [https://docs.docker.com/](https://docs.docker.com/)

## 🤝 **Support**

For issues, questions, or contributions:
- **GitHub Issues**: Report bugs and feature requests
- **Documentation**: Check this comprehensive guide
- **Community**: Engage with the MinIO and Prometheus communities

## 📄 **License**

This project is licensed under the MIT License - see the LICENSE file for details.
