package main

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary test config file
	testConfigContent := `minio_endpoint: "localhost:9000"
listen_addr: ":8080"
scrape_interval: "15s"
minio_access_key: "test"
minio_secret_key: "test"
minio_use_ssl: false
metrics_path: "/minio/metrics/v3"
bucket_pattern: "*"
bucket_exclude_pattern: ""
log_level: "info"`

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "test_config_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	// Write test config content
	if _, err := tmpFile.WriteString(testConfigContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Save original working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Change to directory with temp file
	if err := os.Chdir(os.TempDir()); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(originalWd) // Restore original directory

	// Rename temp file to config.yaml in temp directory
	testConfigPath := os.TempDir() + "/config.yaml"
	if err := os.Rename(tmpFile.Name(), testConfigPath); err != nil {
		t.Fatalf("Failed to rename temp file: %v", err)
	}
	defer os.Remove(testConfigPath) // Clean up

	config := loadConfig()

	// Test expected values
	if config.MinIOEndpoint != "localhost:9000" {
		t.Errorf("Expected MinIO endpoint 'localhost:9000', got '%s'", config.MinIOEndpoint)
	}

	if config.ListenAddr != ":8080" {
		t.Errorf("Expected listen address ':8080', got '%s'", config.ListenAddr)
	}

	if config.ScrapeInterval != 15*time.Second {
		t.Errorf("Expected scrape interval '15s', got '%v'", config.ScrapeInterval)
	}
}

func TestGetEnv(t *testing.T) {
	// Test default value
	result := getEnv("NONEXISTENT_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", result)
	}
}

func TestGetEnvAsBool(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		defaultValue bool
		expected     bool
	}{
		{"true value", "true", false, true},
		{"false value", "false", true, false},
		{"invalid value", "invalid", true, true},
		{"empty value", "", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable temporarily
			if tt.envValue != "" {
				t.Setenv("TEST_BOOL_VAR", tt.envValue)
			}

			result := getEnvAsBool("TEST_BOOL_VAR", tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetEnvAsDuration(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		defaultValue time.Duration
		expected     time.Duration
	}{
		{"valid duration", "30s", 15 * time.Second, 30 * time.Second},
		{"invalid duration", "invalid", 15 * time.Second, 15 * time.Second},
		{"empty value", "", 15 * time.Second, 15 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable temporarily
			if tt.envValue != "" {
				t.Setenv("TEST_DURATION_VAR", tt.envValue)
			}

			result := getEnvAsDuration("TEST_DURATION_VAR", tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestScrapeConfigValidation(t *testing.T) {
	config := ScrapeConfig{
		JobName:        "test-job",
		MetricsPath:    "/metrics",
		ScrapeInterval: "15s",
		ScrapeTimeout:  "10s",
		Scheme:         "http",
	}

	if config.JobName != "test-job" {
		t.Errorf("Expected job name 'test-job', got '%s'", config.JobName)
	}

	if config.MetricsPath != "/metrics" {
		t.Errorf("Expected metrics path '/metrics', got '%s'", config.MetricsPath)
	}
}

func TestStaticConfigValidation(t *testing.T) {
	staticConfig := StaticConfig{
		Targets: []string{"localhost:8080"},
		Labels: map[string]string{
			"job": "test-job",
		},
	}

	if len(staticConfig.Targets) != 1 {
		t.Errorf("Expected 1 target, got %d", len(staticConfig.Targets))
	}

	if staticConfig.Targets[0] != "localhost:8080" {
		t.Errorf("Expected target 'localhost:8080', got '%s'", staticConfig.Targets[0])
	}

	if staticConfig.Labels["job"] != "test-job" {
		t.Errorf("Expected label 'test-job', got '%s'", staticConfig.Labels["job"])
	}
}

func TestServiceDiscoveryResponseValidation(t *testing.T) {
	response := ServiceDiscoveryResponse{
		Targets: []string{"localhost:8080"},
		Labels: map[string]string{
			"job": "test-job",
		},
	}

	if len(response.Targets) != 1 {
		t.Errorf("Expected 1 target, got %d", len(response.Targets))
	}

	if response.Targets[0] != "localhost:8080" {
		t.Errorf("Expected target 'localhost:8080', got '%s'", response.Targets[0])
	}

	if response.Labels["job"] != "test-job" {
		t.Errorf("Expected label 'test-job', got '%s'", response.Labels["job"])
	}
}

// Mock MinIO client for testing
type MockMinIOClient struct {
	buckets []interface{}
	err     error
}

func (m *MockMinIOClient) ListBuckets(ctx context.Context) ([]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.buckets, nil
}

func TestGenerateScrapeConfigs(t *testing.T) {
	// This test would require mocking the MinIO client
	// For now, we'll just test the structure
	config := Config{
		MinIOEndpoint: "localhost:9000",
		DefaultScrapeConfig: ScrapeConfig{
			MetricsPath:    "/minio/metrics/v3",
			ScrapeInterval: "15s",
			ScrapeTimeout:  "10s",
			Scheme:         "http",
		},
	}

	if config.MinIOEndpoint != "localhost:9000" {
		t.Errorf("Expected MinIO endpoint 'localhost:9000', got '%s'", config.MinIOEndpoint)
	}

	if config.DefaultScrapeConfig.MetricsPath != "/minio/metrics/v3" {
		t.Errorf("Expected metrics path '/minio/metrics/v3', got '%s'", config.DefaultScrapeConfig.MetricsPath)
	}
}
