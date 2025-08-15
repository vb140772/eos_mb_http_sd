#!/bin/bash

# EOS MB HTTP SD Helm Chart Installation Script

set -e

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

# Check if Helm is installed
if ! command -v helm &> /dev/null; then
    print_error "Helm is not installed. Please install Helm 3.0+ first."
    exit 1
fi

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    print_error "kubectl is not installed. Please install kubectl first."
    exit 1
fi

# Check if we have access to the cluster
if ! kubectl cluster-info &> /dev/null; then
    print_error "Cannot connect to Kubernetes cluster. Please check your kubeconfig."
    exit 1
fi

print_status "Starting EOS MB HTTP SD Helm chart installation..."

# Add Prometheus Helm repository
print_status "Adding Prometheus Helm repository..."
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# Get chart directory
CHART_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Default values
RELEASE_NAME="eos-mb-http-sd"
NAMESPACE="monitoring"
CREATE_NAMESPACE="--create-namespace"
VALUES_FILE=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--namespace)
            NAMESPACE="$2"
            shift 2
            ;;
        -r|--release-name)
            RELEASE_NAME="$2"
            shift 2
            ;;
        -f|--values)
            VALUES_FILE="$2"
            shift 2
            ;;
        --no-create-namespace)
            CREATE_NAMESPACE=""
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo "Options:"
            echo "  -n, --namespace NAMESPACE     Kubernetes namespace (default: monitoring)"
            echo "  -r, --release-name NAME       Helm release name (default: eos-mb-http-sd)"
            echo "  -f, --values FILE             Custom values file"
            echo "  --no-create-namespace         Don't create namespace if it doesn't exist"
            echo "  -h, --help                    Show this help message"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Build helm install command
HELM_CMD="helm install $RELEASE_NAME $CHART_DIR --namespace $NAMESPACE $CREATE_NAMESPACE"

if [ -n "$VALUES_FILE" ]; then
    if [ ! -f "$VALUES_FILE" ]; then
        print_error "Values file not found: $VALUES_FILE"
        exit 1
    fi
    HELM_CMD="$HELM_CMD -f $VALUES_FILE"
fi

print_status "Installing chart with command: $HELM_CMD"

# Install the chart
if eval $HELM_CMD; then
    print_status "Chart installed successfully!"
    
    echo
    print_status "Next steps:"
    echo "1. Check the status: kubectl get pods -n $NAMESPACE -l 'app.kubernetes.io/name=eos-mb-http-sd'"
    echo "2. View logs: kubectl logs -n $NAMESPACE -l 'app.kubernetes.io/name=eos-mb-http-sd'"
    echo "3. Access the service: kubectl port-forward -n $NAMESPACE svc/$RELEASE_NAME 8080:8080"
    echo "4. Access Prometheus: kubectl port-forward -n $NAMESPACE svc/$RELEASE_NAME-prometheus-server 9090:9090"
    echo
    print_status "Installation complete! 🎉"
else
    print_error "Chart installation failed!"
    exit 1
fi
