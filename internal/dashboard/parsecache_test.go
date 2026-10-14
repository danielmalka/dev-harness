package dashboard

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

const cacheTicket = "| Status | pronta |\n| PRD (RF-<n>) | PRD-001 |\n"

func cacheProject(t *testing.T) (repo, ticket string) {
	t.Helper()
	repo = t.TempDir()
	ticket = filepath.Join(repo, ".harness/tasks/T-1/TASK.md")
	mk(t, ticket, cacheTicket)
	mk(t, filepath.Join(repo, "docs/prd/PRD-001-a.md"), "| Status | aprovado |\n")
	return repo, ticket
}

func TestParseCacheChangedReparsedUnchangedServed(t *testing.T) {
	repo, ticket := cacheProject(t)
	h := filepath.Join(repo, ".harness")
	pc := newProgressCache()
	if pr := projectProgress("p", repo, h, pc); pr.Open[0].Done != 0 || pr.Open[0].Total != 1 {
		t.Fatalf("%+v", pr.Open)
	}
	// same size and mtime, different content: only a cache hit can still say "open"
	info, _ := os.Lstat(ticket)
	mk(t, ticket, "| Status | done   |\n| PRD (RF-<n>) | PRD-001 |\n")
	if err := os.Chtimes(ticket, time.Now(), info.ModTime()); err != nil {
		t.Fatal(err)
	}
	if after, _ := os.Lstat(ticket); after.Size() != info.Size() {
		t.Fatalf("test needs equal sizes: %d vs %d", after.Size(), info.Size())
	}
	if pr := projectProgress("p", repo, h, pc); pr.Open[0].Done != 0 {
		t.Fatalf("unchanged file was re-read: %+v", pr.Open)
	}
	// a different size is a changed file: re-parsed
	mk(t, ticket, "| Status | concluída |\n| PRD (RF-<n>) | PRD-001 |\n")
	if pr := projectProgress("p", repo, h, pc); pr.Open[0].Done != 1 {
		t.Fatalf("changed file was not re-read: %+v", pr.Open)
	}
}

func TestParseCacheSymlinkedTicketIsRefused(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks need privileges on Windows")
	}
	repo, ticket := cacheProject(t)
	h := filepath.Join(repo, ".harness")
	pc := newProgressCache()
	projectProgress("p", repo, h, pc)
	target := filepath.Join(t.TempDir(), "TASK.md")
	mk(t, target, cacheTicket)
	if err := os.Remove(ticket); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, ticket); err != nil {
		t.Fatal(err)
	}
	cached, plain := projectProgress("p", repo, h, pc), projectProgress("p", repo, h, nil)
	if len(plain.Warnings) != 1 || !reflect.DeepEqual(cached, plain) {
		t.Fatalf("cached %+v, uncached %+v", cached, plain)
	}
}

func TestParseCacheSweepDropsUnseen(t *testing.T) {
	repo, ticket := cacheProject(t)
	pc := newProgressCache()
	projectProgress("p", repo, filepath.Join(repo, ".harness"), pc)
	pc.sweep()
	if len(pc.m) != 2 {
		t.Fatalf("entries after first build: %d", len(pc.m))
	}
	if err := os.Remove(ticket); err != nil {
		t.Fatal(err)
	}
	projectProgress("p", repo, filepath.Join(repo, ".harness"), pc)
	pc.sweep()
	if _, ok := pc.m[ticket]; ok || len(pc.m) != 1 {
		t.Fatalf("stale entry kept: %v", pc.m)
	}
}

func TestBuildStateSameWithAndWithoutParseCache(t *testing.T) {
	repo, _ := cacheProject(t)
	home := t.TempDir()
	register(t, home, repo)
	cfg := Config{Home: home, SnapshotDir: filepath.Join(home, "sessions"), Now: func() time.Time { return time.Unix(1_800_000_000, 0) }}
	old := progress
	t.Cleanup(func() { progress = old })
	progress = nil
	want := buildState(cfg, &spriteHolder{})
	progress = newProgressCache()
	for i := 0; i < 2; i++ { // cold, then served from cache
		if got := buildState(cfg, &spriteHolder{}); !reflect.DeepEqual(got, want) {
			t.Fatalf("round %d: cached %+v, uncached %+v", i, got, want)
		}
	}
	if len(want.Projects) != 1 || len(want.Projects[0].Open) != 1 {
		t.Fatalf("fixture not exercised: %+v", want.Projects)
	}
}
