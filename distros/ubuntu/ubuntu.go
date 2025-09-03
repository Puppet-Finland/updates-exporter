package ubuntu

import (
    "log"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

type Ubuntu struct{}

var rebootFile = "/var/run/reboot-required"

func parseUpdateCount(out string) int {
    count, _ := strconv.Atoi(strings.TrimSpace(out))
    return count
}

func (Ubuntu) GetSecurityUpdates() int {
    cmd := exec.Command("sh", "-c", `apt-get -s dist-upgrade | grep "^Inst" | grep security | wc -l`)
    output, err := cmd.Output()
    if err != nil {
        log.Printf("Error running apt-get: %v", err)
        return -1
    }
    return parseUpdateCount(string(output))
}

func (Ubuntu) GetTotalUpdates() int {
    cmd := exec.Command("sh", "-c", `apt-get -s dist-upgrade | grep "^Inst" | wc -l`)
    output, err := cmd.Output()
    if err != nil {
        log.Printf("Error running apt-get: %v", err)
        return -1
    }
    return parseUpdateCount(string(output))
}

func (Ubuntu) GetRebootRequired() bool {
    if _, err := os.Stat(rebootFile); err == nil {
        return true
    }
    return false
}

