package internal

import "testing"

// TestBuilds is a placeholder ensuring the module compiles and `go test` runs.
func TestBuilds(t *testing.T) {
	if 1+1 != 2 {
		t.Fatal("arithmetic broken")
	}
}
