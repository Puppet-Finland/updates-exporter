package ubuntu

import (
	"os"
	"testing"
)

var ubuntu = Ubuntu{}

func TestUbuntuRebootRequired(t *testing.T) {
	// Create a fake reboot-required file
	tmpfile := "/tmp/reboot-required"
	err := os.WriteFile(tmpfile, []byte("reboot"), 0644)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	oldFile := rebootFile
	rebootFile = tmpfile
	defer func() {
		rebootFile = oldFile
	}()

	if !ubuntu.GetRebootRequired() {
		t.Errorf("Expected reboot required = true")
	}
}
