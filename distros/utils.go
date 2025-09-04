package distros

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var rockyReleaseFile = "/etc/rocky-release"

func ParseUpdateCount(out string) int {
	count, _ := strconv.Atoi(strings.TrimSpace(out))
	return count
}

func GetLinuxDistro() string {
	if runtime.GOOS != "linux" {
		return "unknown"
	}

	if _, err := os.Stat(rockyReleaseFile); err == nil {
		return "rocky"
	}

	out, err := exec.Command("sh", "-c", "cat /etc/os-release").Output()
	if err != nil {
		return "unknown"
	}

	s := strings.ToLower(string(out))
	switch {
	case strings.Contains(s, "ubuntu"):
		return "ubuntu"
	case strings.Contains(s, "rhel"), strings.Contains(s, "centos"), strings.Contains(s, "fedora"):
		return "rhel"
	default:
		return "unknown"
	}
}
