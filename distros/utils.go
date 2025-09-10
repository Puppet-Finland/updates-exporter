package distros

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var rhelReleases = []string{
	"/etc/rocky-release",
	"/etc/almalinux-release",
	"/etc/redhat-release",
}

var osReleaseFile = "/etc/os-release"

func ParseUpdateCount(out string) int {
	count, _ := strconv.Atoi(strings.TrimSpace(out))
	return count
}

func GetLinuxDistro() string {
	if runtime.GOOS != "linux" {
		return "unknown"
	}

	for _, releaseFile := range rhelReleases {
		if _, err := os.Stat(releaseFile); err == nil {
			return "rhel"
		}
	}

	out, err := exec.Command("sh", "-c", fmt.Sprintf("cat %s", osReleaseFile)).Output()
	if err != nil {
		return "error"
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
