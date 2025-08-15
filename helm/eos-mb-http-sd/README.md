# EOS MB HTTP SD Helm Chart

A Helm chart for deploying the EOS MB HTTP Service Discovery service along with Prometheus monitoring in Kubernetes.

## Features

- **EOS MB HTTP SD Service**: Service discovery for MinIO clusters
- **Prometheus Integration**: Built-in Prometheus monitoring with service discovery
- **Multi-platform Support**: Supports both AMD64 and ARM64 architectures
- **Configurable**: Easy configuration through values.yaml
- **Production Ready**: Includes health checks, resource limits, and scaling options

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- Access to GitHub Container Registry (GHCR)

## Installation

### 1. Add the Prometheus Helm repository

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
```

### 2. Install the chart

```bash
# Install with default values
helm install eos-mb-http-sd ./helm/eos-mb-http-sd

# Install with custom values
helm install eos-mb-http-sd ./helm/eos-mb-http-sd -f values.yaml

# Install in a specific namespace
helm install eos-mb-http-sd ./helm/eos-mb-http-sd --namespace monitoring --create-namespace
```

### 3. Verify the installation

```bash
kubectl get pods -l "app.kubernetes.io/name=eos-mb-http-sd"
kubectl get services -l "app.kubernetes.io/name=eos-mb-http-sd"
```

## Configuration

### Values.yaml

The following table lists the configurable parameters and their default values:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `eosMbHttpSd.enabled` | Enable the EOS MB HTTP SD service | `true` |
| `eosMbHttpSd.image.repository` | Docker image repository | `ghcr.io/vb140772/eos_mb_http_sd` |
| `eosMbHttpSd.image.tag` | Docker image tag | `latest` |
| `eosMbHttpSd.service.type` | Service type | `ClusterIP` |
| `eosMbHttpSd.service.port` | Service port | `8080` |
| `eosMbHttpSd.resources.limits.cpu` | CPU limit | `500m` |
| `eosMbHttpSd.resources.limits.memory` | Memory limit | `512Mi` |
| `prometheus.enabled` | Enable Prometheus | `true` |
| `bearerToken` | Bearer token for Prometheus authentication | `[your-token]` |

### Custom Configuration

Create a custom `values.yaml` file:

```yaml
eosMbHttpSd:
  config:
    configYaml: |
      minio:
        endpoint: "http://your-minio-service:9000"
        accessKey: "your-access-key"
        secretKey: "your-secret-key"
        useSSL: true
        insecureSkipVerify: false

prometheus:
  server:
    service:
      type: LoadBalancer
    ingress:
      enabled: true
      hosts:
        - host: prometheus.yourdomain.com
          paths:
            - path: /
              pathType: Prefix

bearerToken: "your-bearer-token-here"
```

## Usage

### Accessing the Service

#### EOS MB HTTP SD Service

```bash
# Port forward to access the service
kubectl port-forward svc/eos-mb-http-sd 8080:8080

# Test the service discovery endpoint
curl http://localhost:8080/sd?job=minio-server
curl http://localhost:8080/sd?job=minio-buckets
```

#### Prometheus

```bash
# Port forward to access Prometheus
kubectl port-forward svc/eos-mb-http-sd-prometheus-server 9090:9090

# Open http://localhost:9090 in your browser
```

### Service Discovery Endpoints

The EOS MB HTTP SD service provides the following endpoints:

- `/sd?job=minio-server` - MinIO server metrics targets
- `/sd?job=minio-buckets` - MinIO bucket metrics targets
- `/health` - Health check endpoint
- `/metrics` - Service metrics (if enabled)

### Monitoring

Prometheus is configured to automatically discover and scrape:

1. **MinIO Server Metrics**: Scraped from the service discovery endpoint
2. **MinIO Bucket Metrics**: Scraped from the service discovery endpoint
3. **Self-monitoring**: Prometheus monitoring itself

## Scaling

### Horizontal Pod Autoscaling

Enable HPA by setting:

```yaml
eosMbHttpSd:
  autoscaling:
    enabled: true
    minReplicas: 1
    maxReplicas: 10
    targetCPUUtilizationPercentage: 80
    targetMemoryUtilizationPercentage: 80
```

### Manual Scaling

```bash
# Scale the deployment
kubectl scale deployment eos-mb-http-sd --replicas=3
```

## Troubleshooting

### Check Pod Status

```bash
kubectl get pods -l "app.kubernetes.io/name=eos-mb-http-sd"
kubectl describe pod <pod-name>
```

### Check Logs

```bash
kubectl logs -l "app.kubernetes.io/name=eos-mb-http-sd"
kubectl logs -l "app.kubernetes.io/name=eos-mb-http-sd" --tail=100 -f
```

### Check Service Discovery

```bash
# Test the service discovery endpoint
kubectl exec -it deployment/eos-mb-http-sd -- curl http://localhost:8080/sd?job=minio-server

# Check if the service is responding
kubectl exec -it deployment/eos-mb-http-sd -- curl http://localhost:8080/health
```

### Common Issues

1. **Image Pull Errors**: Ensure you have access to GHCR
2. **Service Discovery Not Working**: Check MinIO endpoint configuration
3. **Prometheus Can't Scrape**: Verify bearer token and network policies

## Upgrading

```bash
# Update the chart
helm upgrade eos-mb-http-sd ./helm/eos-mb-http-sd

# Upgrade with custom values
helm upgrade eos-mb-http-sd ./helm/eos-mb-http-sd -f values.yaml
```

## Uninstalling

```bash
# Uninstall the release
helm uninstall eos-mb-http-sd

# Remove the namespace (if created)
kubectl delete namespace monitoring
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test the chart
5. Submit a pull request

## License

This project is licensed under the MIT License.
