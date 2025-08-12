package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

// Config holds the application configuration
type Config struct {
	MinIOEndpoint        string
	MinIOAccessKey       string
	MinIOSecretKey       string
	MinIOUseSSL          bool
	ListenAddr           string
	ScrapeInterval       time.Duration
	MetricsPath          string
	BucketPattern        string // Wildcard pattern for bucket filtering
	BucketExcludePattern string // Pattern to exclude buckets

	DefaultScrapeConfig ScrapeConfig
}

// ScrapeConfig represents a Prometheus scrape configuration
type ScrapeConfig struct {
	JobName        string          `json:"job_name"`
	StaticConfigs  []StaticConfig  `json:"static_configs"`
	MetricsPath    string          `json:"metrics_path"`
	ScrapeInterval string          `json:"scrape_interval"`
	ScrapeTimeout  string          `json:"scrape_timeout"`
	Scheme         string          `json:"scheme"`
	RelabelConfigs []RelabelConfig `json:"relabel_configs,omitempty"`
}

// StaticConfig represents static targets configuration
type StaticConfig struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

// RelabelConfig represents relabeling configuration
type RelabelConfig struct {
	SourceLabels []string `json:"source_labels"`
	TargetLabel  string   `json:"target_label"`
	Regex        string   `json:"regex"`
	Replacement  string   `json:"replacement"`
}

// ServiceDiscoveryResponse represents the Prometheus service discovery response
type ServiceDiscoveryResponse struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

// MinIOClient wraps the MinIO client
type MinIOClient struct {
	client *minio.Client
	config Config
}

// NewMinIOClient creates a new MinIO client
func NewMinIOClient(config Config) (*MinIOClient, error) {
	client, err := minio.New(config.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinIOAccessKey, config.MinIOSecretKey, ""),
		Secure: config.MinIOUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return &MinIOClient{
		client: client,
		config: config,
	}, nil
}

// ListBuckets retrieves all buckets from MinIO
func (m *MinIOClient) ListBuckets(ctx context.Context) ([]minio.BucketInfo, error) {
	buckets, err := m.client.ListBuckets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list buckets: %w", err)
	}
	return buckets, nil
}

// GenerateScrapeConfigs generates Prometheus scrape configurations for all buckets
func (m *MinIOClient) GenerateScrapeConfigs(ctx context.Context) ([]ScrapeConfig, error) {
	buckets, err := m.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	var configs []ScrapeConfig

	// Generate config for MinIO server metrics (global)
	serverConfig := m.config.DefaultScrapeConfig
	serverConfig.JobName = "minio-server"
	serverConfig.StaticConfigs = []StaticConfig{
		{
			Targets: []string{m.config.MinIOEndpoint},
			Labels: map[string]string{
				"__metrics_path__": "/minio/metrics/v3",
				"__scheme__":       m.getScheme(),
				"instance":         m.config.MinIOEndpoint,
				"job":              "minio-server",
			},
		},
	}
	configs = append(configs, serverConfig)

	// Generate config for all buckets in a single job
	if len(buckets) > 0 {
		bucketConfig := m.config.DefaultScrapeConfig
		bucketConfig.JobName = "minio-buckets"
		bucketConfig.StaticConfigs = []StaticConfig{
			{
				Targets: []string{m.config.MinIOEndpoint},
				Labels: map[string]string{
					"__metrics_path__": "/minio/metrics/v3/bucket/api",
					"__scheme__":       m.getScheme(),
					"instance":         m.config.MinIOEndpoint,
					"job":              "minio-buckets",
					"bucket_pattern":   "*",
				},
			},
		}
		configs = append(configs, bucketConfig)
	}

	return configs, nil
}

// getScheme returns the scheme based on SSL configuration
func (m *MinIOClient) getScheme() string {
	if m.config.MinIOUseSSL {
		return "https"
	}
	return "http"
}

// matchPattern checks if a string matches a wildcard pattern
// Supports * for any sequence of characters and ? for single character
func matchPattern(str, pattern string) bool {
	if pattern == "*" || pattern == "" {
		return true
	}

	// Simple glob pattern matching
	// Convert glob pattern to regex
	regexPattern := "^" + strings.ReplaceAll(pattern, "*", ".*") + "$"
	regexPattern = strings.ReplaceAll(regexPattern, "?", ".")

	matched, err := regexp.MatchString(regexPattern, str)
	if err != nil {
		return false
	}
	return matched
}

// filterBuckets filters buckets based on include/exclude patterns
func (m *MinIOClient) filterBuckets(buckets []minio.BucketInfo) []minio.BucketInfo {
	if m.config.BucketPattern == "*" && m.config.BucketExcludePattern == "" {
		return buckets // No filtering needed
	}

	var filtered []minio.BucketInfo
	for _, bucket := range buckets {
		// Check include pattern
		if !matchPattern(bucket.Name, m.config.BucketPattern) {
			continue
		}

		// Check exclude pattern
		if m.config.BucketExcludePattern != "" && matchPattern(bucket.Name, m.config.BucketExcludePattern) {
			continue
		}

		filtered = append(filtered, bucket)
	}

	return filtered
}

// handleServiceDiscovery handles the /sd endpoint for Prometheus service discovery
func (m *MinIOClient) handleServiceDiscovery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get job name from query parameter
	jobName := r.URL.Query().Get("job")
	if jobName == "" {
		logrus.Warnf("Service discovery request missing job parameter from %s", r.RemoteAddr)
		http.Error(w, "job parameter is required", http.StatusBadRequest)
		return
	}

	logrus.Infof("Service discovery request for job '%s' from %s", jobName, r.RemoteAddr)

	// Generate all scrape configs
	configs, err := m.GenerateScrapeConfigs(ctx)
	if err != nil {
		logrus.Errorf("Failed to generate scrape configs: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Find the requested job
	var targetConfig *ScrapeConfig
	for _, config := range configs {
		if config.JobName == jobName {
			targetConfig = &config
			break
		}
	}

	if targetConfig == nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	// Convert to service discovery format
	var response []ServiceDiscoveryResponse

	// Special handling for minio-buckets job - return all buckets
	if jobName == "minio-buckets" {
		buckets, err := m.ListBuckets(ctx)
		if err != nil {
			logrus.Errorf("Failed to list buckets: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Apply wildcard filtering
		filteredBuckets := m.filterBuckets(buckets)
		logrus.Infof("Found %d buckets, applying pattern '%s' and exclude '%s'", len(buckets), m.config.BucketPattern, m.config.BucketExcludePattern)
		logrus.Infof("After filtering, %d buckets remain", len(filteredBuckets))

		for _, bucket := range filteredBuckets {
			response = append(response, ServiceDiscoveryResponse{
				Targets: []string{m.config.MinIOEndpoint},
				Labels: map[string]string{
					"__metrics_path__":   fmt.Sprintf("/minio/metrics/v3/bucket/api/%s", bucket.Name),
					"__scheme__":         m.getScheme(),
					"instance":           m.config.MinIOEndpoint,
					"job":                "minio-buckets",
					"sd_bucket":          bucket.Name,
					"sd_bucket_creation": bucket.CreationDate.Format(time.RFC3339),
				},
			})
		}
	} else {
		// For other jobs (like minio-server), use the standard approach
		for _, staticConfig := range targetConfig.StaticConfigs {
			for _, target := range staticConfig.Targets {
				// Copy the labels
				labels := make(map[string]string)
				for k, v := range staticConfig.Labels {
					labels[k] = v
				}

				response = append(response, ServiceDiscoveryResponse{
					Targets: []string{target},
					Labels:  labels,
				})
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleScrapeConfigs handles the /scrape_configs endpoint to get all configurations
func (m *MinIOClient) handleScrapeConfigs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logrus.Debugf("Scrape configs request from %s", r.RemoteAddr)

	configs, err := m.GenerateScrapeConfigs(ctx)
	if err != nil {
		logrus.Errorf("Failed to generate scrape configs: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logrus.Debugf("Generated %d scrape configurations", len(configs))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(configs)
}

// handleHealth handles the /health endpoint
func (m *MinIOClient) handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logrus.Debugf("Health check request from %s", r.RemoteAddr)

	// Test MinIO connection
	_, err := m.client.ListBuckets(ctx)
	if err != nil {
		logrus.Warnf("Health check failed - MinIO connection error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "unhealthy",
			"error":     err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	logrus.Debugf("Health check passed - MinIO connection successful")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// loadConfig loads configuration from command line arguments first, then environment variables
func loadConfig() Config {
	// Define command line flags
	var (
		help                 = flag.Bool("help", false, "Show help information")
		minioEndpoint        = flag.String("minio-endpoint", "", "MinIO server endpoint (e.g., localhost:9000)")
		minioAccessKey       = flag.String("minio-access-key", "", "MinIO access key")
		minioSecretKey       = flag.String("minio-secret-key", "", "MinIO secret key")
		minioUseSSL          = flag.Bool("minio-use-ssl", false, "Use SSL for MinIO connection")
		listenAddr           = flag.String("listen-addr", "", "Address to listen on (e.g., :8080)")
		scrapeInterval       = flag.String("scrape-interval", "", "Scrape interval (e.g., 15s)")
		metricsPath          = flag.String("metrics-path", "", "Metrics path (e.g., /minio/metrics/v3)")
		bucketPattern        = flag.String("bucket-pattern", "", "Wildcard pattern for bucket inclusion")
		bucketExcludePattern = flag.String("bucket-exclude-pattern", "", "Wildcard pattern for bucket exclusion")
	)

	// Parse command line flags
	flag.Parse()

	// Show help if requested
	if *help {
		fmt.Println("MinIO Prometheus Service Discovery")
		fmt.Println("Usage:")
		fmt.Println("  minio-prometheus-sd [flags]")
		fmt.Println("")
		fmt.Println("Flags:")
		flag.PrintDefaults()
		fmt.Println("")
		fmt.Println("Environment Variables (used if flags not provided):")
		fmt.Println("  MINIO_ENDPOINT, MINIO_ACCESS_KEY, MINIO_SECRET_KEY, MINIO_USE_SSL")
		fmt.Println("  LISTEN_ADDR, SCRAPE_INTERVAL, METRICS_PATH, BUCKET_PATTERN, BUCKET_EXCLUDE_PATTERN")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  minio-prometheus-sd -minio-endpoint=minio:9000 -minio-access-key=mykey")
		fmt.Println("  minio-prometheus-sd -listen-addr=:9090 -bucket-pattern=prod-*")
		os.Exit(0)
	}

	// Helper function to get value with priority: command line > environment variable > default
	getValue := func(cmdValue *string, envKey, defaultValue string) string {
		if *cmdValue != "" {
			return *cmdValue
		}
		return getEnv(envKey, defaultValue)
	}

	// Helper function to get boolean value with priority: command line > environment variable > default
	getBoolValue := func(cmdValue *bool, envKey string, defaultValue bool) bool {
		if *cmdValue != defaultValue {
			return *cmdValue
		}
		return getEnvAsBool(envKey, defaultValue)
	}

	// Helper function to get duration value with priority: command line > environment variable > default
	getDurationValue := func(cmdValue *string, envKey string, defaultValue time.Duration) time.Duration {
		if *cmdValue != "" {
			if duration, err := time.ParseDuration(*cmdValue); err == nil {
				return duration
			}
		}
		return getEnvAsDuration(envKey, defaultValue)
	}

	config := Config{
		MinIOEndpoint:        getValue(minioEndpoint, "MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey:       getValue(minioAccessKey, "MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey:       getValue(minioSecretKey, "MINIO_SECRET_KEY", "minioadmin"),
		MinIOUseSSL:          getBoolValue(minioUseSSL, "MINIO_USE_SSL", false),
		ListenAddr:           getValue(listenAddr, "LISTEN_ADDR", ":8080"),
		ScrapeInterval:       getDurationValue(scrapeInterval, "SCRAPE_INTERVAL", 15*time.Second),
		MetricsPath:          getValue(metricsPath, "METRICS_PATH", "/minio/metrics/v3"),
		BucketPattern:        getValue(bucketPattern, "BUCKET_PATTERN", "*"),
		BucketExcludePattern: getValue(bucketExcludePattern, "BUCKET_EXCLUDE_PATTERN", ""),
		DefaultScrapeConfig: ScrapeConfig{
			MetricsPath:    "/minio/metrics/v3",
			ScrapeInterval: "15s",
			ScrapeTimeout:  "10s",
			Scheme:         "http",
		},
	}

	// Update scheme based on SSL configuration
	if config.MinIOUseSSL {
		config.DefaultScrapeConfig.Scheme = "https"
	}

	return config
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsBool gets an environment variable as a boolean or returns a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvAsDuration gets an environment variable as a duration or returns a default value
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// maskSensitive masks sensitive configuration values for logging
func maskSensitive(value string) string {
	if value == "" {
		return "(empty)"
	}
	if len(value) <= 4 {
		return "***"
	}
	return value[:2] + "***" + value[len(value)-2:]
}

func main() {
	// Configure logging
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Load configuration
	config := loadConfig()

	// Log configuration
	logrus.Infof("Configuration loaded:")
	logrus.Infof("  MinIO Endpoint: %s", config.MinIOEndpoint)
	logrus.Infof("  MinIO Access Key: %s", maskSensitive(config.MinIOAccessKey))
	logrus.Infof("  MinIO Secret Key: %s", maskSensitive(config.MinIOSecretKey))
	logrus.Infof("  MinIO Use SSL: %t", config.MinIOUseSSL)
	logrus.Infof("  Listen Address: %s", config.ListenAddr)
	logrus.Infof("  Scrape Interval: %v", config.ScrapeInterval)
	logrus.Infof("  Metrics Path: %s", config.MetricsPath)
	logrus.Infof("  Bucket Pattern: %s", config.BucketPattern)
	logrus.Infof("  Bucket Exclude Pattern: %s", config.BucketExcludePattern)

	logrus.Infof("Starting MinIO Prometheus Service Discovery service...")

	// Create MinIO client
	logrus.Infof("Creating MinIO client for endpoint: %s", config.MinIOEndpoint)
	minioClient, err := NewMinIOClient(config)
	if err != nil {
		logrus.Fatalf("Failed to create MinIO client: %v", err)
	}
	logrus.Infof("MinIO client created successfully")

	// Create router
	logrus.Infof("Setting up HTTP router and middleware")
	router := mux.NewRouter()

	// Add middleware for logging
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logrus.Infof("%s %s %v", r.Method, r.URL.Path, time.Since(start))
		})
	})

	// Register routes
	logrus.Infof("Registering HTTP routes:")
	logrus.Infof("  GET /sd - Service discovery endpoint")
	logrus.Infof("  GET /scrape_configs - Scrape configurations endpoint")
	logrus.Infof("  GET /health - Health check endpoint")
	logrus.Infof("  GET / - Documentation endpoint")
	router.HandleFunc("/sd", minioClient.handleServiceDiscovery).Methods("GET")
	router.HandleFunc("/scrape_configs", minioClient.handleScrapeConfigs).Methods("GET")
	router.HandleFunc("/health", minioClient.handleHealth).Methods("GET")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>MinIO Prometheus Service Discovery</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .endpoint { background: #f5f5f5; padding: 10px; margin: 10px 0; border-radius: 5px; }
        .method { color: #0066cc; font-weight: bold; }
        .url { font-family: monospace; }
        .required { color: #cc0000; font-weight: bold; }
    </style>
</head>
<body>
    <h1>MinIO Prometheus Service Discovery</h1>
    <p>This service provides Prometheus HTTP service discovery for MinIO v3 metrics.</p>
    
    <h2>Available Endpoints:</h2>
    
    <div class="endpoint">
        <span class="method">GET</span> <span class="url">/sd?job=minio-server</span>
        <p>Get service discovery targets for MinIO server metrics</p>
    </div>
    
    <div class="endpoint">
        <span class="method">GET</span> <span class="url">/sd?job=minio-buckets</span>
        <p>Get service discovery targets for all MinIO bucket metrics (dynamically filtered)</p>
    </div>
    
    <div class="endpoint">
        <span class="method">GET</span> <span class="url">/scrape_configs</span>
        <p>Get all available scrape configurations</p>
    </div>
    
    <div class="endpoint">
        <span class="method">GET</span> <span class="url">/health</span>
        <p>Health check endpoint (returns JSON status)</p>
    </div>
    
    <h2>Configuration:</h2>
    <p>Set the following environment variables to configure the service:</p>
    <ul>
        <li><strong>MINIO_ENDPOINT</strong>: MinIO server endpoint (default: localhost:9000)</li>
        <li><strong>MINIO_ACCESS_KEY</strong>: MinIO access key (default: minioadmin)</li>
        <li><strong>MINIO_SECRET_KEY</strong>: MinIO secret key (default: minioadmin)</li>
        <li><strong>MINIO_USE_SSL</strong>: Use SSL for MinIO connection (default: false)</li>
        <li><strong>LISTEN_ADDR</strong>: Address to listen on (default: :8080)</li>
        <li><strong>SCRAPE_INTERVAL</strong>: Scrape interval (default: 15s)</li>
        <li><strong>BUCKET_PATTERN</strong>: Wildcard pattern for bucket inclusion (default: *)</li>
        <li><strong>BUCKET_EXCLUDE_PATTERN</strong>: Wildcard pattern for bucket exclusion (default: empty)</li>
    </ul>
    
    <h2>Features:</h2>
    <ul>
        <li>Dynamic bucket discovery with wildcard filtering</li>
        <li>MinIO v3 metrics support (server and bucket metrics)</li>
        <li>Prometheus HTTP Service Discovery compatible</li>
        <li>Configurable bucket inclusion/exclusion patterns</li>
    </ul>
    
    <h2>Getting Started:</h2>
    <ol>
        <li>Start the service: <code>go run main.go</code></li>
        <li>Configure Prometheus to use the service discovery endpoints</li>
    </ol>
</body>
</html>
`)
	}).Methods("GET")

	// Start server
	logrus.Infof("Starting HTTP server on %s", config.ListenAddr)
	logrus.Infof("Service is ready to accept requests")
	if err := http.ListenAndServe(config.ListenAddr, router); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
