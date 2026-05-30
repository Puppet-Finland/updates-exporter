package rhel

import (
	"log/slog"
	"os/exec"

	utils "github.com/Puppet-Finland/updates-exporter/distros"
)

type Rhel struct{}

func (Rhel) GetSecurityUpdates() int {
	cmd := exec.Command("sh", "-c", "dnf updateinfo list --assumeyes --cacheonly --sec-severity=Critical --sec-severity=Important --all | wc -l")
	out, err := cmd.Output()
	if err != nil {
		slog.Error("Failed getting security updates", "error", err)
		return -1
	}
	return utils.ParseUpdateCount(string(out))
}

func (Rhel) GetTotalUpdates() int {
	cmd := exec.Command("sh", "-c", "dnf updateinfo list --assumeyes --cacheonly --all | wc -l")
	out, err := cmd.Output()
	if err != nil {
		slog.Error("Failed getting total updates", "error", err)
		return -1
	}
	return utils.ParseUpdateCount(string(out))
}

func (Rhel) GetRebootRequired() bool {
	cmd := exec.Command("needs-restarting", "-r")
	if err := cmd.Run(); err != nil {
		slog.Error("Failed executing 'needs-restarting'", "error", err)
		return false
	}
	return true
}
