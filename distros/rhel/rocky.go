package rhel

import (
	"log"
	"os/exec"

	utils "github.com/Puppet-Finland/updates-exporter/distros"
)

type Rocky struct{}

func (Rocky) GetSecurityUpdates() int {
	cmd := exec.Command("sh", "-c", "dnf updateinfo list --sec-severity=Critical --sec-severity=Important --all | wc -l")
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Error running dnf: %v", err)
		return -1
	}
	return utils.ParseUpdateCount(string(out))
}

func (Rocky) GetTotalUpdates() int {
	cmd := exec.Command("sh", "-c", "dnf updateinfo list --all | wc -l")
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Error running dnf: %v", err)
		return -1
	}
	return utils.ParseUpdateCount(string(out))
}

func (Rocky) GetRebootRequired() bool {
	cmd := exec.Command("needs-restarting", "-r")
	err := cmd.Run()
	return err != nil
}
