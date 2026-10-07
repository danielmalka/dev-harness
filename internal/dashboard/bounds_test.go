package dashboard

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/danielmalka/dev-harness/internal/snapshot"
)

func big(t *testing.T, path string) {
	t.Helper()
	mk(t, path, strings.Repeat("x", MaxFileBytes+1))
}

func TestOversizedFilesAreSkipped(t *testing.T) {
	p := t.TempDir()
	mk(t, filepath.Join(p, "docs/prd/PRD-001-a.md"), "| Status | aprovado |\n")
	big(t, filepath.Join(p, "docs/prd/PRD-002-b.md"))
	big(t, filepath.Join(p, ".harness/tasks/T-1/TASK.md"))
	pr := ProjectProgress(p)
	if len(pr.Open) != 1 || pr.Open[0].PRD != "PRD-001" || len(pr.Warnings) != 2 {
		t.Fatalf("%+v", pr)
	}
	sd := t.TempDir()
	big(t, filepath.Join(sd, "poses.json"))
	if s := LoadSprites(sd); len(s.Warnings) != 1 || s.PoseFor(StateIdle).Pose != "frente" {
		t.Fatalf("poses: %+v", s.Warnings)
	}
	big(t, filepath.Join(sd, "frente.png"))
	if b, _ := LoadSprites(sd).Frame("frente", 0); len(b) > MaxFileBytes {
		t.Fatal("oversized sprite served")
	}
	snaps := t.TempDir()
	big(t, filepath.Join(snaps, "a.json"))
	if _, err := snapshot.Load(filepath.Join(snaps, "a.json")); err == nil {
		t.Fatal("oversized snapshot read")
	}
}

func TestReviewFilesAreNotPRDs(t *testing.T) {
	p := t.TempDir()
	mk(t, filepath.Join(p, "docs/prd/PRD-001-a.md"), "| Status | aprovado |\n")
	mk(t, filepath.Join(p, "docs/prd/PRD-001-a.review.md"), "| Status | entregue em 2026-01-01 |\n")
	pr := ProjectProgress(p)
	if len(pr.Delivered) != 0 || len(pr.Open) != 1 {
		t.Fatalf("%+v", pr)
	}
}

func TestCountCaps(t *testing.T) {
	root := t.TempDir()
	for i := 0; i < MaxProjects+1; i++ {
		mk(t, filepath.Join(root, fmt.Sprintf("p%03d", i), ".harness", "x"), "")
	}
	if p, w := Projects(root); len(p) != MaxProjects || w == "" {
		t.Fatalf("projects %d warn %q", len(p), w)
	}
	pp := t.TempDir()
	for i := 0; i < MaxTicketsPerProj+1; i++ {
		mk(t, filepath.Join(pp, ".harness/tasks", fmt.Sprintf("T-%d", i), "TASK.md"), "| Status | pronta |\n| PRD (RF-<n>) | PRD-001 |\n")
	}
	mk(t, filepath.Join(pp, "docs/prd/PRD-001-a.md"), "| Status | aprovado |\n")
	pr := ProjectProgress(pp)
	if pr.Open[0].Total != MaxTicketsPerProj || len(pr.Warnings) != 1 {
		t.Fatalf("%+v", pr)
	}
	sd := t.TempDir()
	for i := 0; i < MaxSnapshots+2; i++ {
		f := filepath.Join(sd, fmt.Sprintf("s%03d.json", i))
		mk(t, f, fmt.Sprintf(`{"session_id":"s%d"}`, i))
		old := time.Now().Add(-time.Duration(MaxSnapshots+2-i) * time.Minute) // higher i = newer
		_ = os.Chtimes(f, old, old)
	}
	got, skipped := snapshot.LoadDirMax(sd, MaxSnapshots)
	if len(got) != MaxSnapshots || skipped != 2 {
		t.Fatalf("got %d skipped %d", len(got), skipped)
	}
	for _, s := range got {
		if s.SessionID == "s0" || s.SessionID == "s1" {
			t.Fatalf("oldest kept: %s", s.SessionID)
		}
	}
}

func TestStateCache(t *testing.T) {
	root, snaps := t.TempDir(), t.TempDir()
	mk(t, filepath.Join(root, "a", ".harness", "x"), "")
	h := Handler(Config{Roots: root, SnapshotDir: snaps, Stale: DefaultStale, DoneDecay: time.Minute, CacheTTL: time.Hour})
	first := do(h, "GET", "127.0.0.1", "/api/state").Body.String()
	mk(t, filepath.Join(root, "b", ".harness", "x"), "")
	if second := do(h, "GET", "127.0.0.1", "/api/state").Body.String(); second != first {
		t.Fatal("state rebuilt inside the TTL")
	}
}
