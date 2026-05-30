package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
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
	Port     int  `mapstructure:"port"`
	Interval int  `mapstructure:"interval"`
	Version  bool `mapstructure:"version"`
	Debug    bool `mapstructure:"debug"`
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

	if d.GetRebootRequired() {
		rebootRequired.Set(1)
	} else {
		rebootRequired.Set(0)
	}

}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if cfg.Debug {
		cfgBytes, err := json.Marshal(cfg)
		if err != nil {
			log.Fatalf("Failed to marshal config for debugging: %v", err)
		}

		log.Printf("Starting with: %s", string(cfgBytes))
	}

	prometheus.MustRegister(securityUpdates)
	prometheus.MustRegister(totalUpdates)
	prometheus.MustRegister(rebootRequired)

	distro := getDistro()
	if distro == nil {
		log.Fatal("Error: Distro not detected")
	}

	go func() {
		for {
			updateMetrics(distro)
			time.Sleep(time.Duration(cfg.Interval) * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	fmt.Printf("Starting updates_exporter on :%d, updating every %d seconds\n", cfg.Port, cfg.Interval)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil))
}

func loadConfig() (*Config, error) {
	pflag.IntP("port", "p", DEFAULT_PORT, "HTTP port")
	pflag.IntP("interval", "i", DEFAULT_INTERVAL, "Metrics refresh interval (seconds)")
	pflag.BoolP("version", "v", false, "Print version")
	pflag.BoolP("debug", "d", false, "Enables Debug messages")

	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, fmt.Errorf("unable to bind flags: %w", err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // Look for the file in the current working directory
	viper.AddConfigPath("/etc/updates-exporter/")

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
