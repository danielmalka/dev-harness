package dashboard

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/danielmalka/dev-harness/internal/snapshot"
)

func big(t *testing.T, path string) {
	t.Helper()
	mk(t, path, strings.Repeat("x", snapshot.MaxFileBytes+1))
}

func TestOversizedFilesAreSkipped(t *testing.T) {
	p := t.TempDir()
	mk(t, filepath.Join(p, "docs/prd/PRD-001-a.md"), "| Status | aprovado |\n")
	big(t, filepath.Join(p, "docs/prd/PRD-002-b.md"))
	big(t, filepath.Join(p, ".harness/tasks/T-1/TASK.md"))
	pr := ProjectProgress(p, filepath.Join(p, ".harness"))
	if len(pr.Open) != 1 || pr.Open[0].PRD != "PRD-001" || len(pr.Warnings) != 2 {
		t.Fatalf("%+v", pr)
	}
	sd := t.TempDir()
	big(t, filepath.Join(sd, "poses.json"))
	if s := LoadSprites(sd); len(s.Warnings) != 1 || s.PoseFor(StateIdle).Pose != "frente" {
		t.Fatalf("poses: %+v", s.Warnings)
	}
	big(t, filepath.Join(sd, "frente.png"))
	if b, _ := LoadSprites(sd).Frame("frente", 0); len(b) > snapshot.MaxFileBytes {
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
	pr := ProjectProgress(p, filepath.Join(p, ".harness"))
	if len(pr.Delivered) != 0 || len(pr.Open) != 1 {
		t.Fatalf("%+v", pr)
	}
}

func TestCountCaps(t *testing.T) {
	pp := t.TempDir()
	for i := 0; i < MaxTicketsPerProj+1; i++ {
		mk(t, filepath.Join(pp, ".harness/tasks", fmt.Sprintf("T-%d", i), "TASK.md"), "| Status | pronta |\n| PRD (RF-<n>) | PRD-001 |\n")
	}
	mk(t, filepath.Join(pp, "docs/prd/PRD-001-a.md"), "| Status | aprovado |\n")
	pr := ProjectProgress(pp, filepath.Join(pp, ".harness"))
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
	home, snaps := t.TempDir(), t.TempDir()
	a, b := t.TempDir(), t.TempDir()
	mk(t, filepath.Join(a, ".harness", "x"), "")
	mk(t, filepath.Join(b, ".harness", "x"), "")
	register(t, home, a)
	h := Handler(Config{Home: home, SnapshotDir: snaps, Stale: DefaultStale, DoneDecay: time.Minute, CacheTTL: time.Hour})
	first := do(h, "GET", "127.0.0.1", "/api/state").Body.String()
	register(t, home, a, b)
	if second := do(h, "GET", "127.0.0.1", "/api/state").Body.String(); second != first {
		t.Fatal("state rebuilt inside the TTL")
	}
}

func TestSymlinkToDevZeroIsSkippedPromptly(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks need privileges on Windows")
	}
	p := t.TempDir()
	dir := filepath.Join(p, ".harness/tasks/T-1")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("/dev/zero", filepath.Join(dir, "TASK.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("/dev/zero", filepath.Join(p, "snap.json")); err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	pr := ProjectProgress(p, filepath.Join(p, ".harness"))
	_, err := snapshot.Load(filepath.Join(p, "snap.json"))
	if time.Since(start) > 2*time.Second || len(pr.Warnings) != 1 || err == nil {
		t.Fatalf("slow or not skipped: %v %+v %v", time.Since(start), pr.Warnings, err)
	}
}

func TestSnapshotLoadSizeBoundary(t *testing.T) {
	d := t.TempDir()
	mk(t, filepath.Join(d, "big.json"), strings.Repeat(" ", snapshot.MaxFileBytes+1))
	if _, err := snapshot.Load(filepath.Join(d, "big.json")); !errors.Is(err, snapshot.ErrTooLarge) {
		t.Fatalf("got %v", err)
	}
	mk(t, filepath.Join(d, "ok.json"), `{"session_id":"x"}`+strings.Repeat(" ", snapshot.MaxFileBytes-18))
	if s, err := snapshot.Load(filepath.Join(d, "ok.json")); err != nil || s.SessionID != "x" {
		t.Fatalf("at the limit must load: %v", err)
	}
}

func TestPRDFileCap(t *testing.T) {
	p := t.TempDir()
	for i := 0; i < MaxTicketsPerProj+2; i++ {
		mk(t, filepath.Join(p, "docs/prd", fmt.Sprintf("PRD-%04d-a.md", i)), "| Status | aprovado |\n")
	}
	pr := ProjectProgress(p, filepath.Join(p, ".harness"))
	if len(pr.Open) != MaxTicketsPerProj || len(pr.Warnings) != 1 {
		t.Fatalf("open %d warns %v", len(pr.Open), pr.Warnings)
	}
}

func TestReviewFilesDoNotEatThePRDCap(t *testing.T) {
	p := t.TempDir()
	for i := 0; i < MaxTicketsPerProj; i++ {
		mk(t, filepath.Join(p, "docs/prd", fmt.Sprintf("PRD-%04d-a.review.md", i)), "| Status | aprovado |\n")
	}
	mk(t, filepath.Join(p, "docs/prd", "PRD-9999-real.md"), "| Status | aprovado |\n")
	pr := ProjectProgress(p, filepath.Join(p, ".harness"))
	if len(pr.Open) != 1 || pr.Open[0].PRD != "PRD-9999" || len(pr.Warnings) != 0 {
		t.Fatalf("%+v", pr)
	}
}

func TestFrameCountIsCappedAndReadsNothing(t *testing.T) {
	d := t.TempDir()
	for i := 0; i < MaxFrames+5; i++ {
		mk(t, filepath.Join(d, fmt.Sprintf("celular_%02d.png", i)), "x")
	}
	if n := LoadSprites(d).FrameCount("celular"); n != MaxFrames {
		t.Fatalf("frames %d", n)
	}
}
