package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/codemeapixel/cadence/internal/detector"
)

type Config struct {
	Thresholds   detector.Thresholds
	ExcludeFiles []string
	Webhook      WebhookConfig
	AI           AIConfig
}

// WebhookConfig holds webhook server configuration
type WebhookConfig struct {
	Enabled      bool
	Host         string
	Port         int
	Secret       string
	MaxWorkers   int
	ReadTimeout  int
	WriteTimeout int
}

// AIConfig holds AI analysis configuration
type AIConfig struct {
	Enabled  bool
	Provider string
	APIKey   string
	Model    string
}

func Load(configFile string) (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("thresholds.suspicious_additions", 500)
	v.SetDefault("thresholds.suspicious_deletions", 1000)
	v.SetDefault("thresholds.max_additions_per_min", 100)
	v.SetDefault("thresholds.max_deletions_per_min", 500)
	v.SetDefault("thresholds.min_time_delta_seconds", 60)
	v.SetDefault("thresholds.max_files_per_commit", 50)
	v.SetDefault("thresholds.max_addition_ratio", 0.95)
	v.SetDefault("thresholds.min_deletion_ratio", 0.95)
	v.SetDefault("thresholds.min_commit_size_ratio", 100)
	v.SetDefault("thresholds.enable_precision_analysis", true)

	if configFile != "" {
		v.SetConfigFile(configFile)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	v.SetEnvPrefix("CADENCE")
	v.AutomaticEnv()

	config := &Config{}

	config.Thresholds.SuspiciousAdditions = v.GetInt64("thresholds.suspicious_additions")
	config.Thresholds.SuspiciousDeletions = v.GetInt64("thresholds.suspicious_deletions")
	config.Thresholds.MaxAdditionsPerMin = v.GetFloat64("thresholds.max_additions_per_min")
	config.Thresholds.MaxDeletionsPerMin = v.GetFloat64("thresholds.max_deletions_per_min")
	config.Thresholds.MinTimeDeltaSeconds = v.GetInt64("thresholds.min_time_delta_seconds")
	config.Thresholds.MaxFilesPerCommit = v.GetInt("thresholds.max_files_per_commit")
	config.Thresholds.MaxAdditionRatio = v.GetFloat64("thresholds.max_addition_ratio")
	config.Thresholds.MinDeletionRatio = v.GetFloat64("thresholds.min_deletion_ratio")
	config.Thresholds.MinCommitSizeRatio = v.GetInt64("thresholds.min_commit_size_ratio")
	config.Thresholds.EnablePrecisionAnalysis = v.GetBool("thresholds.enable_precision_analysis")

	config.ExcludeFiles = v.GetStringSlice("exclude_files")

	// Load webhook configuration
	config.Webhook.Enabled = v.GetBool("webhook.enabled")
	config.Webhook.Host = v.GetString("webhook.host")
	if config.Webhook.Host == "" {
		config.Webhook.Host = "0.0.0.0"
	}
	config.Webhook.Port = v.GetInt("webhook.port")
	if config.Webhook.Port == 0 {
		config.Webhook.Port = 3000
	}
	config.Webhook.Secret = v.GetString("webhook.secret")
	config.Webhook.MaxWorkers = v.GetInt("webhook.max_workers")
	if config.Webhook.MaxWorkers == 0 {
		config.Webhook.MaxWorkers = 4
	}
	config.Webhook.ReadTimeout = v.GetInt("webhook.read_timeout")
	if config.Webhook.ReadTimeout == 0 {
		config.Webhook.ReadTimeout = 30
	}
	config.Webhook.WriteTimeout = v.GetInt("webhook.write_timeout")
	if config.Webhook.WriteTimeout == 0 {
		config.Webhook.WriteTimeout = 30
	}

	// Load AI configuration
	config.AI.Enabled = v.GetBool("ai.enabled")
	config.AI.Provider = v.GetString("ai.provider")
	config.AI.APIKey = v.GetString("ai.api_key")
	config.AI.Model = v.GetString("ai.model")
	if config.AI.Model == "" {
		config.AI.Model = "gpt-4o-mini"
	}

	return config, nil
}

func GenerateSampleConfig(path string) error {
	sample := `# Cadence Configuration - AI-Generated Code Detection
# Analyzes git repositories to detect potential AI-generated code patterns

thresholds:
  # SIZE-BASED DETECTION
  suspicious_additions: 500
  suspicious_deletions: 1000
  
  # VELOCITY-BASED DETECTION
  max_additions_per_min: 100
  max_deletions_per_min: 500
  
  # TIMING-BASED DETECTION
  min_time_delta_seconds: 60
  
  # FILE DISPERSION DETECTION
  max_files_per_commit: 50
  
  # RATIO-BASED DETECTION
  max_addition_ratio: 0.95
  min_deletion_ratio: 0.95
  min_commit_size_ratio: 100
  
  # PRECISION ANALYSIS
  enable_precision_analysis: true

# File patterns to exclude from analysis
exclude_files:
  - package-lock.json
  - yarn.lock
  - "*.min.js"
  - "*.min.css"
  - "node_modules/**"

# WEBHOOK SERVER CONFIGURATION
webhook:
  # Enable/disable webhook server
  enabled: false
  
  # Server host and port
  host: "0.0.0.0"
  port: 3000
  
  # Webhook secret for signature verification (set this!)
  secret: "your-webhook-secret-key-here"
  
  # Number of concurrent workers for processing webhook events
  max_workers: 4
  
  # Request timeouts in seconds
  read_timeout: 30
  write_timeout: 30

# AI ANALYSIS CONFIGURATION (Optional - requires API key)
ai:
  # Enable/disable AI-powered code analysis
  enabled: false
  
  # AI provider ("openai" only for now)
  provider: "openai"
  
  # OpenAI API key (or set via CADENCE_AI_KEY environment variable)
  api_key: ""
  
  # OpenAI model (gpt-4o-mini recommended for efficiency)
  model: "gpt-4o-mini"
`

	return os.WriteFile(path, []byte(sample), 0o600)
}
