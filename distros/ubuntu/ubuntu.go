package ubuntu

import (
	"log"
	"os"
	"os/exec"

	utils "github.com/Puppet-Finland/updates-exporter/distros"
)

type Ubuntu struct{}

var rebootFile = "/var/run/reboot-required"

func (Ubuntu) GetSecurityUpdates() int {
	cmd := exec.Command("sh", "-c", `apt-get -s dist-upgrade | grep "^Inst" | grep security | wc -l`)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error running apt-get: %v", err)
		return -1
	}
	return utils.ParseUpdateCount(string(output))
}

func (Ubuntu) GetTotalUpdates() int {
	cmd := exec.Command("sh", "-c", `apt-get -s dist-upgrade | grep "^Inst" | wc -l`)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error running apt-get: %v", err)
		return -1
	}
	return utils.ParseUpdateCount(string(output))
}

func (Ubuntu) GetRebootRequired() bool {
	if _, err := os.Stat(rebootFile); err == nil {
		return true
	}
	return false
}
