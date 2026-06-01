package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Puppet-Finland/updates-exporter/distros"
	"github.com/Puppet-Finland/updates-exporter/distros/rhel"
	"github.com/Puppet-Finland/updates-exporter/distros/ubuntu"
)

type Config struct {
	Port     int    `mapstructure:"port"`
	Interval int    `mapstructure:"interval"`
	Version  bool   `mapstructure:"version"`
	LogLevel string `mapstructure:"log-level"`
	Config   string `mapstructure:"config"`
}

const (
	DEFAULT_INTERVAL = 3600
	DEFAULT_PORT     = 9101
)

var (
	securityUpdates = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pending_security_updates",
		Help: "Number of pending security updates",
	})
	totalUpdates = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pending_updates",
		Help: "Total number of pending updates",
	})
	rebootRequired = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "reboot_required",
		Help: "1 if a reboot is required, 0 otherwise",
	})
	latestCacheUpdate = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "latest_cache_update",
		Help: "Unix timestamp of latest change to update cache",
	})
	Version = "dev"
)

func getDistro() distros.Distro {
	switch distros.GetLinuxDistro() {
	case "ubuntu":
		return ubuntu.Ubuntu{}
	case "rhel":
		return rhel.Rhel{}
	default:
		return nil
	}
}

func updateMetrics(d distros.Distro) {
	if d == nil {
		return
	}
	securityUpdates.Set(float64(d.GetSecurityUpdates()))
	totalUpdates.Set(float64(d.GetTotalUpdates()))
	latestCacheUpdate.Set(float64(d.GetLatestCacheChange().Unix()))

	if d.GetRebootRequired() {
		rebootRequired.Set(1)
	} else {
		rebootRequired.Set(0)
	}

}

func main() {
	logLevel := &slog.LevelVar{}
	logLevel.Set(slog.LevelInfo)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	// slog.Info("Application initialized", "version", "v1.4.2", "environment", "production")
	cfg, err := loadConfig()
	if err != nil {
		slog.Error("Unable to load config", "error", err)
		os.Exit(1)
	}

	// 3. Parse and dynamically set the log level from the configuration string
	var parsedLevel slog.Level
	// UnmarshalText handles string inputs like "DEBUG", "info", "Warn", "ERROR" case-insensitively
	if err := parsedLevel.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		slog.Error("Invalid log level provided, falling back to INFO", "invalid_level", cfg.LogLevel, "error", err)
		parsedLevel = slog.LevelInfo
	}

	logLevel.Set(parsedLevel)

	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		slog.Error("Unable to marshal config", "error", err)
	}

	slog.Debug("Application starting with config", "config", string(cfgBytes))

	if cfg.Version {
		fmt.Println(Version)
		os.Exit(0)
	}

	prometheus.MustRegister(securityUpdates)
	prometheus.MustRegister(totalUpdates)
	prometheus.MustRegister(rebootRequired)

	distro := getDistro()
	if distro == nil {
		slog.Error("Distro not detected")
		os.Exit(1)
	}

	go func() {
		slog.Debug("Initialized go routine")

		for {
			updateMetrics(distro)
			slog.Debug("Sleeping", "interval", cfg.Interval)
			time.Sleep(time.Duration(cfg.Interval) * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	addr := fmt.Sprintf(":%d", cfg.Port)
	slog.Info("Starting HTTP server", "addr", addr, "interval", cfg.Interval)

	if err := http.ListenAndServe(addr, nil); err != nil {
		slog.Error("HTTP server collapsed", "error", err)
		os.Exit(1)
	}

	slog.Info("HTTP server stopped gracefully")
}

func loadConfig() (*Config, error) {
	pflag.IntP("port", "p", DEFAULT_PORT, "HTTP port")
	pflag.IntP("interval", "i", DEFAULT_INTERVAL, "Metrics refresh interval (seconds)")
	pflag.BoolP("version", "v", false, "Print version")
	pflag.StringP("log-level", "l", "info", "Log verbosity level (debug, info, warn, error)")
	pflag.StringP("config", "c", "", "Config file")

	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, fmt.Errorf("unable to bind flags: %w", err)
	}

	// Environment variable matching setup
	viper.SetEnvPrefix("UPDATES_EXPORTER")
	viper.AutomaticEnv()
	// Critical: Converts dashes in flags ("log-level") to underscores for ENVs ("LOG_LEVEL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	configFlag := viper.GetString("config")

	if configFlag != "" {
		slog.Info("Loading configuration from explicit config flag", "path", configFlag)
		viper.SetConfigFile(configFlag)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".") // Look for the file in the current working directory
		viper.AddConfigPath("/etc/updates_exporter/")
	}

	if err := viper.ReadInConfig(); err != nil {
		// It is acceptable if the config file is missing; we fall back to flags/defaults.
		// But if it exists and has a syntax error, we want to fail fast.
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}
