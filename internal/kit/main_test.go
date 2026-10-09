package kit

import (
	"os"
	"testing"
)

// TestMain isolates every test from the real ~/.harness (config is read by validate).
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "dh-home-")
	if err != nil {
		panic(err)
	}
	os.Setenv("DH_HOME", dir)
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}
