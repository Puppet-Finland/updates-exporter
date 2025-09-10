package distros

import (
	"os"
	"testing"
)

func getRelease(r string, t *testing.T) {

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
	osReleaseFile = "/tmp/os-release"
	rhelReleases = []string{
		"/tmp/rocky-release",
	}
	os.Remove(rhelReleases[0])

	os.WriteFile(osReleaseFile, []byte("ubuntu"), 0644)
	got := GetLinuxDistro()
	expected := "ubuntu"
	if got != expected {
		t.Errorf("Expected %s, got %s", expected, got)
	}

	os.WriteFile(osReleaseFile, []byte("fedora"), 0644)
	got = GetLinuxDistro()
	expected = "rhel"
	if got != expected {
		t.Errorf("Expected %s, got %s", expected, got)
	}

	if err := os.WriteFile(rhelReleases[0], []byte("test"), 0644); err != nil {
		t.Errorf("Error creating %s", rhelReleases[0])
	}

	got = GetLinuxDistro()
	if got != expected {
		t.Errorf("Expected %s, got %s", expected, got)
	}
}
