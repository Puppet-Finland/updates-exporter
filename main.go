package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
    DEFAULT_INTERVAL = 86400
    DEFAULT_PORT     = 9101
)
var (
    securityUpdates = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "ubuntu_pending_security_updates",
            Help: "Number of pending security updates on Ubuntu",
        },
    )
    rebootRequired = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "ubuntu_reboot_required",
            Help: "1 if a reboot is required, 0 otherwise",
        },
    )
)

func init() {
    prometheus.MustRegister(securityUpdates)
    prometheus.MustRegister(rebootRequired)
}

func getPendingSecurityUpdates() int {
    cmd := exec.Command("sh", "-c", `apt-get -s dist-upgrade | grep "^Inst" | grep security | wc -l`)
    output, err := cmd.Output()
    if err != nil {
        log.Printf("Error running apt-get: %v", err)
        return -1
    }
    count, err := strconv.Atoi(strings.TrimSpace(string(output)))
    if err != nil {
        log.Printf("Error parsing security update count: %v", err)
        return -1
    }
    return count
}

func checkRebootRequired() bool {
    if _, err := os.Stat("/var/run/reboot-required"); err == nil {
        return true
    }
    return false
}

func updateMetrics() {
    securityUpdates.Set(float64(getPendingSecurityUpdates()))
    if checkRebootRequired() {
        rebootRequired.Set(1)
    } else {
        rebootRequired.Set(0)
    }
}

func main() {
    // Command-line flags
    port := flag.Int("port", DEFAULT_PORT, "HTTP port to serve metrics")
    interval := flag.Int("interval", DEFAULT_INTERVAL, "Metrics update interval in seconds")
    flag.Parse()

    // Update metrics periodically
    go func() {
        for {
            updateMetrics()
            time.Sleep(time.Duration(*interval) * time.Second)
        }
    }()

    http.Handle("/metrics", promhttp.Handler())
    fmt.Printf("Starting updates_exporter on :%d, updating every %d seconds\n", *port, *interval)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
