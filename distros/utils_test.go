package distros

import (
	"os"
	"testing"
)

func getRelease(r string, t *testing.T) {
	osReleaseFile = "/tmp/os-release"

	os.WriteFile(osReleaseFile, []byte(r), 0644)
	got := GetLinuxDistro()
	expected := r
	if got != expected {
		t.Errorf("Expected %s, got %s", expected, got)
	}
}

func TestParseUpdateCount(t *testing.T) {
	output := "3\n"
	expected := 3
	got := ParseUpdateCount(output)
	if got != expected {
		t.Errorf("Expected %d, got %d", expected, got)
	}
}

func TestGetDistros(t *testing.T) {
	getRelease("alma", t)
	getRelease("rhel", t)
	getRelease("ubuntu", t)
}

func TestGetDistrosRocky(t *testing.T) {
	rockyReleaseFile = "/tmp/rocky-release"

	os.WriteFile(rockyReleaseFile, []byte("rocky"), 0644)
	got := GetLinuxDistro()
	expected := "rocky"
	if got != expected {
		t.Errorf("Expected %s, got %s", expected, got)
	}
}
