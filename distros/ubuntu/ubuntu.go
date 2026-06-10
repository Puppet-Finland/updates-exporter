package ubuntu

import (
	"log/slog"
	"os"
	"os/exec"
	"time"

	utils "github.com/Puppet-Finland/updates-exporter/distros"
)

type Ubuntu struct{}

const (
	REBOOT_FILE = "/var/run/reboot-required"
	CACHE_DIR   = " /var/lib/apt/lists/"
)

func (Ubuntu) GetLatestCacheChange() time.Time {
	change, err := utils.GetLatestChangeInDir(CACHE_DIR, ".lz4")
	if err != nil {
		slog.Error("Failed to get latest cache change", slog.Any("error", err))
		return time.Time{}
	}
	return change
}

func (Ubuntu) GetSecurityUpdates() int {
	cmd := exec.Command("sh", "-c", `apt-get -s dist-upgrade | grep "^Inst" | grep security | wc -l`)
	output, err := cmd.Output()
	if err != nil {
		slog.Error("Failed running apt-get", slog.Any("error", err), slog.String("cmd", cmd.String()))
		return -1
	}
	return utils.ParseUpdateCount(string(output))
}

func (Ubuntu) GetTotalUpdates() int {
	cmd := exec.Command("sh", "-c", `apt-get -s dist-upgrade | grep "^Inst" | wc -l`)
	output, err := cmd.Output()
	if err != nil {
		slog.Error("Failed running apt-get", slog.Any("error", err), slog.String("cmd", cmd.String()))
		return -1
	}
	return utils.ParseUpdateCount(string(output))
}

func (Ubuntu) GetRebootRequired() bool {
	if _, err := os.Stat(REBOOT_FILE); err == nil {
		return true
	}
	return false
}
