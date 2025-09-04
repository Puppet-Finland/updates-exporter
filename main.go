package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Puppet-Finland/updates-exporter/distros"
	"github.com/Puppet-Finland/updates-exporter/distros/rhel"
	"github.com/Puppet-Finland/updates-exporter/distros/ubuntu"
)

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
)

func getDistro() distros.Distro {
	switch distros.GetLinuxDistro() {
	case "ubuntu":
		return ubuntu.Ubuntu{}
	case "rocky":
		return rhel.Rocky{}
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
	port := flag.Int("port", DEFAULT_PORT, "HTTP port")
	interval := flag.Int("interval", DEFAULT_INTERVAL, "Metrics refresh interval (seconds)")

	flag.Parse()

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
			time.Sleep(time.Duration(*interval) * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	fmt.Printf("Starting updates_exporter on :%d, updating every %d seconds\n", *port, *interval)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
