# Systemd Service Installation Guide

This guide explains how to install and configure the EOS MinIO Prometheus Service Discovery as a systemd service on Linux systems.

## Overview

The service will run the `eos_mb_http_sd` binary as a background service that automatically starts on boot, restarts on failure, and provides proper logging through systemd journal.

## Prerequisites

- Linux system with systemd
- Root access (sudo privileges)
- Built `eos_mb_http_sd` binary
- MinIO server access credentials

## Files Included

- `eos-mb-http-sd.service` - Systemd service file with default values
- `eos-mb-http-sd.service.template` - Template with placeholder values
- `install-service.sh` - Automated installation script
- `uninstall-service.sh` - Service removal script

## Quick Installation

### 1. Download the Binary

Download the pre-built binary from the GitHub releases page:

```bash
# Download the latest release binary for Linux x86_64
wget https://github.com/your-username/eos_mb_http_sd/releases/latest/download/eos_mb_http_sd-linux-amd64

# Rename to the expected binary name
mv eos_mb_http_sd-linux-amd64 eos_mb_http_sd

# Make it executable
chmod +x eos_mb_http_sd
```

**Alternative download methods:**

```bash
# Using curl
curl -L -o eos_mb_http_sd https://github.com/your-username/eos_mb_http_sd/releases/latest/download/eos_mb_http_sd-linux-amd64

# Or download a specific version
wget https://github.com/your-username/eos_mb_http_sd/releases/download/v1.0.0/eos_mb_http_sd-linux-amd64
```

**Available architectures:**
- Linux x86_64: `eos_mb_http_sd-linux-amd64`
- Linux ARM64: `eos_mb_http_sd-linux-arm64`
- Linux ARM: `eos_mb_http_sd-linux-arm`
- macOS x86_64: `eos_mb_http_sd-darwin-amd64`
- macOS ARM64: `eos_mb_http_sd-darwin-arm64`

### 2. Run Installation Script

```bash
# Make script executable
chmod +x install-service.sh

# Run installation (requires sudo)
sudo ./install-service.sh
```

The script will:
- Create a dedicated system user and group
- Install the binary to `/opt/eos-mb-http-sd/`
- Create log directory at `/var/log/eos-mb-http-sd/`
- Install and enable the systemd service
- Set proper permissions

## Manual Installation

If you prefer to install manually or customize the installation:

### 1. Create Service User

```bash
sudo useradd --system --no-create-home --shell /sbin/nologin eos-mb-http-sd
sudo groupadd --system eos-mb-http-sd
sudo usermod -a -G eos-mb-http-sd eos-mb-http-sd
```

### 2. Create Directories

```bash
sudo mkdir -p /opt/eos-mb-http-sd
sudo mkdir -p /var/log/eos-mb-http-sd
sudo chown eos-mb-http-sd:eos-mb-http-sd /opt/eos-mb-http-sd
sudo chown eos-mb-http-sd:eos-mb-http-sd /var/log/eos-mb-http-sd
```

### 3. Install Binary

```bash
sudo cp eos_mb_http_sd /opt/eos-mb-http-sd/
sudo chown eos-mb-http-sd:eos-mb-http-sd /opt/eos-mb-http-sd/eos_mb_http_sd
sudo chmod 755 /opt/eos-mb-http-sd/eos_mb_http_sd
```

### 4. Install Service File

```bash
sudo cp eos-mb-http-sd.service /etc/systemd/system/
sudo chmod 644 /etc/systemd/system/eos-mb-http-sd.service
```

### 5. Enable and Start Service

```bash
sudo systemctl daemon-reload
sudo systemctl enable eos-mb-http-sd
sudo systemctl start eos-mb-http-sd
```

## Configuration

### Environment Variables

The service uses environment variables for configuration. Edit the service file to customize:

```bash
sudo nano /etc/systemd/system/eos-mb-http-sd.service
```

#### MinIO Connection Settings

```ini
Environment="MINIO_ENDPOINT=your-minio-server:9000"
Environment="MINIO_ACCESS_KEY=your-access-key"
Environment="MINIO_SECRET_KEY=your-secret-key"
Environment="MINIO_USE_SSL=false"
```

#### Service Settings

```ini
Environment="LISTEN_ADDR=:8080"
Environment="SCRAPE_INTERVAL=15s"
Environment="METRICS_PATH=/minio/metrics/v3"
```

#### Bucket Filtering

```ini
Environment="BUCKET_PATTERN=*"
Environment="BUCKET_EXCLUDE_PATTERN="
```

#### Logging

```ini
Environment="LOG_LEVEL=info"
```

### Configuration Priority

The service follows this priority order (highest to lowest):
1. Environment variables in service file
2. Configuration file (if specified)
3. Default values

## Service Management

### Start/Stop/Restart

```bash
# Start the service
sudo systemctl start eos-mb-http-sd

# Stop the service
sudo systemctl stop eos-mb-http-sd

# Restart the service
sudo systemctl restart eos-mb-http-sd

# Reload configuration (after editing service file)
sudo systemctl daemon-reload
sudo systemctl restart eos-mb-http-sd
```

### Check Status

```bash
# Check service status
sudo systemctl status eos-mb-http-sd

# Check if service is enabled
sudo systemctl is-enabled eos-mb-http-sd

# Check if service is running
sudo systemctl is-active eos-mb-http-sd
```

### View Logs

```bash
# View real-time logs
sudo journalctl -u eos-mb-http-sd -f

# View recent logs
sudo journalctl -u eos-mb-http-sd -n 100

# View logs since boot
sudo journalctl -u eos-mb-http-sd -b

# View logs for specific time period
sudo journalctl -u eos-mb-http-sd --since "2024-01-01 00:00:00"
```



## Troubleshooting

### Common Issues

#### Service Won't Start

```bash
# Check service status
sudo systemctl status eos-mb-http-sd

# Check logs for errors
sudo journalctl -u eos-mb-http-sd -n 50

# Check if binary exists and is executable
ls -la /opt/eos-mb-http-sd/eos_mb_http_sd

# Check permissions
sudo -u eos-mb-http-sd /opt/eos-mb-http-sd/eos_mb_http_sd --help
```

#### Permission Denied

```bash
# Check user and group ownership
ls -la /opt/eos-mb-http-sd/
ls -la /var/log/eos-mb-http-sd/

# Fix permissions if needed
sudo chown -R eos-mb-http-sd:eos-mb-http-sd /opt/eos-mb-http-sd
sudo chown -R eos-mb-http-sd:eos-mb-http-sd /var/log/eos-mb-http-sd
```

#### Network Issues

```bash
# Check if service is listening
sudo netstat -tlnp | grep :8080
sudo ss -tlnp | grep :8080

# Test connectivity to MinIO
curl -v http://your-minio-server:9000/minio/health/live
```

### Debug Mode

To run the service in debug mode temporarily:

```bash
# Stop the service
sudo systemctl stop eos-mb-http-sd

# Run manually with debug logging
sudo -u eos-mb-http-sd /opt/eos-mb-http-sd/eos_mb_http_sd -log-level=debug
```

## Uninstallation

To completely remove the service:

```bash
# Run uninstall script
sudo ./uninstall-service.sh

# Or manually remove
sudo systemctl stop eos-mb-http-sd
sudo systemctl disable eos-mb-http-sd
sudo rm /etc/systemd/system/eos-mb-http-sd.service
sudo systemctl daemon-reload
sudo rm -rf /opt/eos-mb-http-sd
sudo rm -rf /var/log/eos-mb-http-sd
sudo userdel eos-mb-http-sd
sudo groupdel eos-mb-http-sd
```

## Integration with Monitoring

### Prometheus Configuration

Add this to your Prometheus configuration:

```yaml
scrape_configs:
  - job_name: 'minio-service-discovery'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: /metrics
    scrape_interval: 15s
```

### Health Checks

The service provides health endpoints:

- **Health check**: `GET /health`
- **Metrics**: `GET /metrics`
- **Service discovery**: `GET /sd`

## Support

For issues or questions:

1. Check the service logs: `sudo journalctl -u eos-mb-http-sd`
2. Verify configuration in the service file
3. Test the binary manually
4. Check system resources and permissions
5. Review this documentation

## Files Summary

| File | Purpose | Location |
|------|---------|----------|
| `eos-mb-http-sd.service` | Systemd service file | `/etc/systemd/system/` |
| `eos-mb-http-sd.service.template` | Template with placeholders | Project directory |
| `install-service.sh` | Automated installation | Project directory |
| `uninstall-service.sh` | Service removal | Project directory |
| Binary | Executable | `/opt/eos-mb-http-sd/` |
| Logs | Service logs | Systemd journal |
