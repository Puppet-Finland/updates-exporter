package distros

import (
	"testing"
)

func TestParseUpdateCount(t *testing.T) {
	output := "3\n"
	expected := 3
	got := ParseUpdateCount(output)
	if got != expected {
		t.Errorf("Expected %d, got %d", expected, got)
	}
}
