#!/bin/bash

# EOS MinIO Prometheus Service Discovery - Service Uninstallation Script
# This script removes the eos_mb_http_sd systemd service and cleans up

set -e

# Configuration
SERVICE_NAME="eos-mb-http-sd"
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

print_status "Starting uninstallation of $SERVICE_NAME service..."

# Stop and disable service if running
if systemctl is-active --quiet "$SERVICE_NAME"; then
    print_status "Stopping service..."
    systemctl stop "$SERVICE_NAME"
else
    print_status "Service is not running"
fi

if systemctl is-enabled --quiet "$SERVICE_NAME"; then
    print_status "Disabling service..."
    systemctl disable "$SERVICE_NAME"
else
    print_status "Service is not enabled"
fi

# Remove service file
if [[ -f "/etc/systemd/system/$SERVICE_FILE" ]]; then
    print_status "Removing service file..."
    rm -f "/etc/systemd/system/$SERVICE_FILE"
else
    print_status "Service file not found"
fi

# Reload systemd
print_status "Reloading systemd daemon..."
systemctl daemon-reload

# Remove binary and installation directory
if [[ -d "$INSTALL_DIR" ]]; then
    print_status "Removing installation directory..."
    rm -rf "$INSTALL_DIR"
else
    print_status "Installation directory not found"
fi

# Remove log directory
if [[ -d "$LOG_DIR" ]]; then
    print_status "Removing log directory..."
    rm -rf "$LOG_DIR"
else
    print_status "Log directory not found"
fi

# Remove service user and group
if id "$SERVICE_USER" &>/dev/null; then
    print_status "Removing service user..."
    userdel "$SERVICE_USER"
else
    print_status "Service user not found"
fi

if getent group "$SERVICE_GROUP" &>/dev/null; then
    print_status "Removing service group..."
    groupdel "$SERVICE_GROUP"
else
    print_status "Service group not found"
fi

print_status "Uninstallation completed successfully!"
echo
print_status "The following items have been removed:"
echo "  - Systemd service: $SERVICE_NAME"
echo "  - Service file: /etc/systemd/system/$SERVICE_FILE"
echo "  - Binary and installation directory: $INSTALL_DIR"
echo "  - Log directory: $LOG_DIR"
echo "  - Service user: $SERVICE_USER"
echo "  - Service group: $SERVICE_GROUP"
echo
print_status "Note: Any custom configuration files or data in these directories have been permanently deleted."
print_status "If you need to preserve any data, please backup before running this script."
