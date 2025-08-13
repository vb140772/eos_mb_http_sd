#!/bin/bash

# EOS MinIO Prometheus Service Discovery - Service Installation Script
# This script installs the eos_mb_http_sd binary as a systemd service

set -e

# Configuration
SERVICE_NAME="eos-mb-http-sd"
BINARY_NAME="eos_mb_http_sd"
INSTALL_DIR="/opt/eos-mb-http-sd"
SERVICE_USER="eos-mb-http-sd"
SERVICE_GROUP="eos-mb-http-sd"
SERVICE_FILE="eos-mb-http-sd.service"
LOG_DIR="/var/log/eos-mb-http-sd"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   print_error "This script must be run as root (use sudo)"
   exit 1
fi

print_status "Starting installation of $SERVICE_NAME service..."

# Check if binary exists in current directory
if [[ ! -f "./$BINARY_NAME" ]]; then
    print_error "Binary '$BINARY_NAME' not found in current directory"
    print_error "Please build the binary first or run this script from the directory containing the binary"
    exit 1
fi

# Create service user and group
print_status "Creating service user and group..."
if ! id "$SERVICE_USER" &>/dev/null; then
    useradd --system --no-create-home --shell /sbin/nologin "$SERVICE_USER"
    print_status "Created user: $SERVICE_USER"
else
    print_status "User $SERVICE_USER already exists"
fi

if ! getent group "$SERVICE_GROUP" &>/dev/null; then
    groupadd --system "$SERVICE_GROUP"
    print_status "Created group: $SERVICE_GROUP"
else
    print_status "Group $SERVICE_GROUP already exists"
fi

# Add user to group
usermod -a -G "$SERVICE_GROUP" "$SERVICE_USER"

# Create installation directory
print_status "Creating installation directory..."
mkdir -p "$INSTALL_DIR"
chown "$SERVICE_USER:$SERVICE_GROUP" "$INSTALL_DIR"
chmod 755 "$INSTALL_DIR"

# Create log directory
print_status "Creating log directory..."
mkdir -p "$LOG_DIR"
chown "$SERVICE_USER:$SERVICE_GROUP" "$LOG_DIR"
chmod 755 "$LOG_DIR"

# Copy binary
print_status "Installing binary..."
cp "./$BINARY_NAME" "$INSTALL_DIR/"
chown "$SERVICE_USER:$SERVICE_GROUP" "$INSTALL_DIR/$BINARY_NAME"
chmod 755 "$INSTALL_DIR/$BINARY_NAME"

# Copy service file
print_status "Installing systemd service..."
cp "$SERVICE_FILE" "/etc/systemd/system/"
chmod 644 "/etc/systemd/system/$SERVICE_FILE"

# Reload systemd
print_status "Reloading systemd daemon..."
systemctl daemon-reload

# Enable service
print_status "Enabling service..."
systemctl enable "$SERVICE_NAME"

print_status "Installation completed successfully!"
echo
print_status "Service details:"
echo "  Service name: $SERVICE_NAME"
echo "  Binary location: $INSTALL_DIR/$BINARY_NAME"
echo "  Service file: /etc/systemd/system/$SERVICE_FILE"
echo "  Log directory: $LOG_DIR"
echo "  Service user: $SERVICE_USER"
echo
print_status "To start the service, run:"
echo "  sudo systemctl start $SERVICE_NAME"
echo
print_status "To check service status, run:"
echo "  sudo systemctl status $SERVICE_NAME"
echo
print_status "To view logs, run:"
echo "  sudo journalctl -u $SERVICE_NAME -f"
echo
print_warning "IMPORTANT: Please review and update the environment variables in the service file:"
echo "  sudo nano /etc/systemd/system/$SERVICE_FILE"
echo
print_warning "Key variables to configure:"
echo "  - MINIO_ENDPOINT: Your MinIO server endpoint"
echo "  - MINIO_ACCESS_KEY: Your MinIO access key"
echo "  - MINIO_SECRET_KEY: Your MinIO secret key"
echo "  - MINIO_USE_SSL: Set to 'true' if using HTTPS"
echo "  - LISTEN_ADDR: Service listening address (default: :8080)"
echo
print_status "After updating the service file, reload and restart:"
echo "  sudo systemctl daemon-reload"
echo "  sudo systemctl restart $SERVICE_NAME"
