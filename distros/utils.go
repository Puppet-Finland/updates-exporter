package distros

import (
    "runtime"
    "strings"
    "os/exec"
)

// DetectLinuxDistro returns "ubuntu", "rhel", or "unknown"
func GetLinuxDistro() string {
    if runtime.GOOS != "linux" {
        return "unknown"
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
