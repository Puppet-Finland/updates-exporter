package ubuntu

import (
    "log"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

type Ubuntu struct{}

func (Ubuntu) GetSecurityUpdates() int {
    cmd := exec.Command("sh", "-c", `apt-get -s dist-upgrade | grep "^Inst" | grep security | wc -l`)
    output, err := cmd.Output()
    if err != nil {
        log.Printf("Error running apt-get: %v", err)
        return 0
    }
    count, _ := strconv.Atoi(strings.TrimSpace(string(output)))
    return count
}

func (Ubuntu) GetRebootRequired() bool {
    if _, err := os.Stat("/var/run/reboot-required"); err == nil {
        return true
    }
    return false
}

