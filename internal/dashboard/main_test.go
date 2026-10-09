package dashboard

import (
	"os"
	"testing"
)

// TestMain keeps every test off the real ~/.harness.
func TestMain(m *testing.M) {
	d, err := os.MkdirTemp("", "dh-dashboard-test-")
	if err != nil {
		panic(err)
	}
	os.Setenv("DH_HOME", d)
	code := m.Run()
	os.RemoveAll(d)
	os.Exit(code)
}
