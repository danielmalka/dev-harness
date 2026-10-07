package dashboard

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestChildArgsNeverDetach(t *testing.T) {
	got := strings.Join(childArgs(7000, time.Minute, 2*time.Minute), " ")
	if got != "dashboard --port 7000 --stale 1m0s --done-decay 2m0s" || strings.Contains(got, "detach") {
		t.Fatalf("argv: %q", got)
	}
}

// The child (this test binary, bad flag) exits at once; the probe answers only after it, as when a concurrent
// --detach wins the port. Expect success, not an error.
func TestDetachRaceChildExitedButDashboardUp(t *testing.T) {
	calls := 0
	probe := func(int) bool { calls++; return calls > 1 }
	var out, errb bytes.Buffer
	if rc := startDetachedWith(1, "http://x/", []string{"-test.run=^$"}, probe, &out, &errb); rc != 0 || !strings.Contains(out.String(), "http://x/") {
		t.Fatalf("rc=%d out=%q err=%q", rc, out.String(), errb.String())
	}
}
