package rhel

import (
	"errors"
	"log/slog"
	"os/exec"
	"time"

	utils "github.com/Puppet-Finland/updates-exporter/distros"
)

type Rhel struct{}

var cache_dir = "/var/cache/libdnf5"

func (Rhel) GetLatestCacheChange() time.Time {
	if !utils.DirExists(cache_dir) {
		cache_dir = "/var/cache/dnf"
	}
	change, err := utils.GetLatestChangeInDir(cache_dir, "")
	if err != nil {
		slog.Error("Failed to get latest cache change", slog.Any("error", err))
		return time.Time{}
	}
	return change
}

func (Rhel) GetSecurityUpdates() int {
	cmd := exec.Command("sh", "-c", "dnf updateinfo list --assumeyes --cacheonly --sec-severity=Critical --sec-severity=Important --all | wc -l")
	out, err := cmd.Output()
	if err != nil {
		slog.Error("Failed getting security updates", slog.Any("error", err))
		return -1
	}
	return utils.ParseUpdateCount(string(out))
}

func (Rhel) GetTotalUpdates() int {
	cmd := exec.Command("sh", "-c", "dnf updateinfo list --assumeyes --cacheonly --all | wc -l")
	out, err := cmd.Output()
	if err != nil {
		slog.Error("Failed getting total updates", slog.Any("error", err))
		return -1
	}
	return utils.ParseUpdateCount(string(out))
}

func (Rhel) GetRebootRequired() bool {
	cmd := exec.Command("needs-restarting", "-r")
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode := exitErr.ExitCode()

			if exitCode == 1 {
				slog.Debug("System needs restarting", slog.Int("exit_code", exitCode))
				return true
			}
		}
		slog.Error("Failed executing 'needs-restarting'", slog.Any("error", err))
		return false
	}
	slog.Debug("System restart not needed")
	return false
}
