package harness

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// wtFx is a hand-built linked-worktree layout, all under t.TempDir() (never the repo or ~).
type wtFx struct{ base, main, common, wt, admin string }

func writeF(t *testing.T, p, s string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
}

// fakeRepo makes <dir>/.git with the HEAD, objects and refs the resolver requires.
func fakeRepo(t *testing.T, dir string) string {
	t.Helper()
	g := filepath.Join(dir, ".git")
	writeF(t, filepath.Join(g, "HEAD"), "ref: refs/heads/main\n")
	for _, d := range []string{"objects", "refs"} {
		if err := os.MkdirAll(filepath.Join(g, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return g
}

func realTemp(t *testing.T) string {
	t.Helper()
	b, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// addWT adds a valid linked worktree <base>/wts/<name> to the repo at main (common dir main/.git).
func addWT(t *testing.T, base, main, name string) wtFx {
	t.Helper()
	f := wtFx{base: base, main: main, common: filepath.Join(main, ".git"), wt: filepath.Join(base, "wts", name)}
	f.admin = filepath.Join(f.common, "worktrees", name)
	writeF(t, filepath.Join(f.admin, "commondir"), "../..\n")
	writeF(t, filepath.Join(f.admin, "gitdir"), filepath.Join(f.wt, ".git")+"\n")
	writeF(t, filepath.Join(f.wt, ".git"), "gitdir: "+f.admin+"\n")
	return f
}

func newWT(t *testing.T) wtFx {
	t.Helper()
	base := realTemp(t)
	main := filepath.Join(base, "main")
	fakeRepo(t, main)
	return addWT(t, base, main, "T-1")
}

func (f wtFx) gitfile(t *testing.T, s string) { t.Helper(); writeF(t, filepath.Join(f.wt, ".git"), s) }

func symlinkOrSkip(t *testing.T, target, link string) {
	t.Helper()
	if err := os.Symlink(target, link); err != nil {
		t.Skip("symlink unavailable: ", err)
	}
}

// R26, R61: a valid layout is recognised and every field is canonical.
func TestR26_R61_ValidLayoutFields(t *testing.T) {
	f := newWT(t)
	got := LinkedWorktree(f.wt)
	if got == nil {
		t.Fatal("valid worktree not recognised")
	}
	want := LinkedWT{Main: f.main, Gitdir: f.admin, Commondir: f.common}
	if *got != want {
		t.Errorf("got %+v want %+v", *got, want)
	}
	if !strings.HasSuffix(got.Gitdir, filepath.Join(".git", "worktrees", "T-1")) {
		t.Errorf("gitdir must be the admin dir, got %s", got.Gitdir)
	}
}

// R26: relative gitdir and relative back-pointer/commondir resolve against their own directory,
// never against the process cwd.
func TestR26_RelativePathsNotResolvedAgainstCwd(t *testing.T) {
	f := newWT(t)
	f.gitfile(t, "gitdir: ../../main/.git/worktrees/T-1\n")
	writeF(t, filepath.Join(f.admin, "gitdir"), "../../../../wts/T-1/.git\n")
	elsewhere := realTemp(t)
	old, _ := os.Getwd()
	if err := os.Chdir(elsewhere); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(old) })
	got := LinkedWorktree(f.wt)
	if got == nil || got.Gitdir != f.admin || got.Main != f.main {
		t.Fatalf("relative layout refused or wrong: %+v", got)
	}
	// a back-pointer that only matches when resolved against cwd must be refused
	writeF(t, filepath.Join(elsewhere, "wts", "T-1", ".git"), "x")
	writeF(t, filepath.Join(f.admin, "gitdir"), "wts/T-1/.git\n")
	if LinkedWorktree(f.wt) != nil {
		t.Error("back-pointer resolved against cwd")
	}
}

// R26 (ADR-008 D3): `..` resolved by the OS (after following existing directories) to the
// SAME valid admin dir is accepted; a text-only `worktrees/../worktrees/<name>` too.
func TestR26_DotDotToSameAdminAccepted(t *testing.T) {
	f := newWT(t)
	f.gitfile(t, "gitdir: ../../main/.git/worktrees/../worktrees/T-1\n")
	if got := LinkedWorktree(f.wt); got == nil || got.Gitdir != f.admin {
		t.Errorf("text-only worktrees/../worktrees: %+v", got)
	}
	if err := os.MkdirAll(filepath.Join(f.wt, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	f.gitfile(t, "gitdir: sub/../../../main/.git/worktrees/T-1\n")
	if got := LinkedWorktree(f.wt); got == nil || got.Gitdir != f.admin {
		t.Errorf("same admin through an existing sub/..: %+v", got)
	}
}

// R26, R27: `..` through a MISSING directory is refused: the OS (and git) cannot resolve it, and
// the resolver has no lexical fallback.
func TestR26_R27_DotDotThroughMissingDirRefused(t *testing.T) {
	f := newWT(t)
	f.gitfile(t, "gitdir: nosuch/../../../main/.git/worktrees/T-1\n")
	refuse(t, f)
}

// deepDir makes <base>/deep/a/b/c/d/e/f and returns it and its ancestors by number of levels up.
func deepDir(t *testing.T, base string) string {
	t.Helper()
	d := filepath.Join(base, "deep", "a", "b", "c", "d", "e", "f")
	if err := os.MkdirAll(d, 0o755); err != nil {
		t.Fatal(err)
	}
	return d
}

// R26, R27: gitdir `L/../../../main/...` where L is a symlink to a deep dir. A textual Clean
// would give <wt>/../../main (the valid admin dir); the OS gives <deep>/a/b/c/main/...
func TestR26_R27_SymlinkDotDotInGitdirRefused(t *testing.T) {
	line := "gitdir: L/../../../main/.git/worktrees/T-1\n"
	t.Run("OS target missing", func(t *testing.T) {
		f := newWT(t)
		symlinkOrSkip(t, deepDir(t, f.base), filepath.Join(f.wt, "L"))
		f.gitfile(t, line)
		refuse(t, f)
	})
	t.Run("OS target is a valid other layout with its own back-pointer", func(t *testing.T) {
		f := newWT(t)
		deep := deepDir(t, f.base)
		symlinkOrSkip(t, deep, filepath.Join(f.wt, "L"))
		c := filepath.Join(f.base, "deep", "a", "b", "c") // 3 levels above deep
		fakeRepo(t, filepath.Join(c, "main"))
		o := addWT(t, c, filepath.Join(c, "main"), "T-1") // back-pointer -> <c>/wts/T-1/.git
		if !isDir(o.admin) {
			t.Fatal("fixture")
		}
		f.gitfile(t, line)
		refuse(t, f)
	})
}

// R26, R27: the same L/.. trick inside admin/commondir and inside the admin/gitdir back-pointer.
func TestR26_R27_SymlinkDotDotInCommondirAndBackPointerRefused(t *testing.T) {
	t.Run("commondir", func(t *testing.T) {
		f := newWT(t)
		symlinkOrSkip(t, deepDir(t, f.base), filepath.Join(f.admin, "L"))
		writeF(t, filepath.Join(f.admin, "commondir"), "L/../../..\n") // text: ../.. = common dir
		refuse(t, f)
	})
	t.Run("back-pointer", func(t *testing.T) {
		f := newWT(t)
		deep := deepDir(t, f.base)
		symlinkOrSkip(t, deep, filepath.Join(f.admin, "L"))
		// decoy at the OS-resolved location so the path resolves and only the comparison can refuse
		writeF(t, filepath.Join(f.base, "deep", "a", "wts", "T-1", ".git"), "x")
		writeF(t, filepath.Join(f.admin, "gitdir"), "L/../../../../../wts/T-1/.git\n")
		refuse(t, f)
	})
}

// R26 controls that must stay ACCEPTED.
func TestR26_SymlinkedAncestorWithoutDotDotAccepted(t *testing.T) {
	f := newWT(t)
	symlinkOrSkip(t, f.main, filepath.Join(f.base, "mainlink"))
	f.gitfile(t, "gitdir: "+filepath.Join(f.base, "mainlink", ".git", "worktrees", "T-1")+"\n")
	got := LinkedWorktree(f.wt)
	if got == nil || got.Main != f.main || got.Gitdir != f.admin || got.Commondir != f.common {
		t.Errorf("symlinked ancestor: %+v", got)
	}
}

// R26, R27: the raw admin path must be <common>/worktrees/<name> literally.
func TestR26_R27_RawAdminShapeRefusals(t *testing.T) {
	cases := []struct {
		name string
		mut  func(t *testing.T, f wtFx)
	}{
		{"middle segment is not worktrees", func(t *testing.T, f wtFx) {
			odd := filepath.Join(f.common, "wt-admin", "T-1")
			writeF(t, filepath.Join(odd, "commondir"), "../..\n")
			writeF(t, filepath.Join(odd, "gitdir"), filepath.Join(f.wt, ".git")+"\n")
			f.gitfile(t, "gitdir: "+odd+"\n")
		}},
		{"empty name (trailing slash)", func(t *testing.T, f wtFx) { f.gitfile(t, "gitdir: "+f.admin+"/\n") }},
		{"name is ..", func(t *testing.T, f wtFx) { f.gitfile(t, "gitdir: "+f.admin+"/..\n") }},
		{"symlink at raw/../worktrees", func(t *testing.T, f wtFx) {
			real := filepath.Join(f.common, "wt-real")
			dir := filepath.Dir(f.admin)
			if err := os.Rename(dir, real); err != nil {
				t.Fatal(err)
			}
			symlinkOrSkip(t, real, dir)
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := newWT(t)
			c.mut(t, f)
			refuse(t, f)
		})
	}
}

// R26 + worktree of a worktree: both point at the same common dir, so both resolve to main.
func TestR26_WorktreeOfWorktree(t *testing.T) {
	f := newWT(t)
	g := addWT(t, f.base, f.main, "T-2") // created "from" T-1: git still records the main common dir
	for _, d := range []string{f.wt, g.wt} {
		if got := LinkedWorktree(d); got == nil || got.Main != f.main || got.Commondir != f.common {
			t.Errorf("%s: %+v", d, got)
		}
	}
	if got := LinkedWorktree(g.wt); got.Gitdir != g.admin {
		t.Errorf("gitdir %s", got.Gitdir)
	}
}

// R21: a plain checkout (.git is a directory) and a plain folder are not linked worktrees.
func TestR21_NormalCheckoutAndPlainDirAreNotWorktrees(t *testing.T) {
	base := realTemp(t)
	main := filepath.Join(base, "main")
	fakeRepo(t, main)
	if LinkedWorktree(main) != nil {
		t.Error(".git directory treated as a linked worktree")
	}
	plain := filepath.Join(base, "plain")
	os.MkdirAll(plain, 0o755)
	if LinkedWorktree(plain) != nil || LinkedWorktree(filepath.Join(base, "missing")) != nil {
		t.Error("non-git dir treated as a linked worktree")
	}
	writeF(t, filepath.Join(main, ".harness", "project.yaml"), "x\n")
	if r := resolveIn(realTemp(t), main, Config{}); r != (Resolved{"repo", filepath.Join(main, ".harness")}) {
		t.Errorf("normal checkout resolution changed: %+v", r)
	}
	if r := resolveIn(realTemp(t), plain, Config{}); r != (Resolved{"none", ""}) {
		t.Errorf("plain dir: %+v", r)
	}
}

// R18, R19: repo mode ignores the .harness copy inside the worktree.
func TestR18_R19_RepoModeResolvesMainHarness(t *testing.T) {
	f := newWT(t)
	writeF(t, filepath.Join(f.main, ".harness", "project.yaml"), "main\n")
	writeF(t, filepath.Join(f.wt, ".harness", "project.yaml"), "copy\n")
	want := Resolved{"repo", filepath.Join(f.main, ".harness")}
	if r := resolveIn(realTemp(t), f.wt, Config{}); r != want {
		t.Errorf("got %+v want %+v", r, want)
	}
	if r := resolveIn(realTemp(t), f.main, Config{}); r != want {
		t.Errorf("main: %+v", r)
	}
}

// R18: global mode returns the same dir for the worktree and for the main checkout.
func TestR18_GlobalModeSameDirAsMain(t *testing.T) {
	f := newWT(t)
	home := realTemp(t)
	writeF(t, filepath.Join(home, "projects", "proj", "project.yaml"), "x\n")
	c := Config{Projects: map[string]string{f.main: "proj"}}
	want := Resolved{"global", filepath.Join(home, "projects", "proj")}
	if r := resolveIn(home, f.wt, c); r != want {
		t.Errorf("worktree: %+v want %+v", r, want)
	}
	if r := resolveIn(home, f.main, c); r != want {
		t.Errorf("main: %+v", r)
	}
	// not linked at all: a worktree folder alone is not found through the registry
	if r := resolveIn(home, f.wt, Config{}); r.Mode != "none" {
		t.Errorf("unregistered worktree: %+v", r)
	}
}

// refuse asserts a layout is rejected and resolution is what it is today (own .harness copy).
func refuse(t *testing.T, f wtFx) {
	t.Helper()
	if got := LinkedWorktree(f.wt); got != nil {
		t.Errorf("accepted: %+v", *got)
	}
	writeF(t, filepath.Join(f.main, ".harness", "project.yaml"), "main\n")
	writeF(t, filepath.Join(f.wt, ".harness", "project.yaml"), "copy\n")
	if r := resolveIn(realTemp(t), f.wt, Config{}); r != (Resolved{"repo", filepath.Join(f.wt, ".harness")}) {
		t.Errorf("resolution changed for a refused layout: %+v", r)
	}
}

// R26, R27: every refusal of ADR-008 section 5.
func TestR26_R27_Refusals(t *testing.T) {
	cases := []struct {
		name string
		mut  func(t *testing.T, f wtFx)
	}{
		{"dotdot to another repository", func(t *testing.T, f wtFx) {
			// other repo has a consistent admin dir of its OWN worktree: only the back-pointer differs
			other := filepath.Join(f.base, "other")
			fakeRepo(t, other)
			o := addWT(t, f.base, other, "O-1")
			rel, _ := filepath.Rel(f.wt, o.admin)
			if !strings.Contains(rel, "..") {
				t.Fatal("fixture: no ..")
			}
			f.gitfile(t, "gitdir: "+rel+"\n")
		}},
		{"absolute path to another repository", func(t *testing.T, f wtFx) {
			other := filepath.Join(f.base, "other")
			fakeRepo(t, other)
			o := addWT(t, f.base, other, "O-1")
			f.gitfile(t, "gitdir: "+o.admin+"\n")
		}},
		{"back-pointer mismatch", func(t *testing.T, f wtFx) {
			writeF(t, filepath.Join(f.admin, "gitdir"), filepath.Join(f.base, "wts", "elsewhere", ".git")+"\n")
		}},
		{"commondir mismatch", func(t *testing.T, f wtFx) {
			other := filepath.Join(f.base, "other")
			fakeRepo(t, other)
			writeF(t, filepath.Join(f.admin, "commondir"), filepath.Join(other, ".git")+"\n")
		}},
		{"commondir missing", func(t *testing.T, f wtFx) { os.Remove(filepath.Join(f.admin, "commondir")) }},
		{"bare common dir", func(t *testing.T, f wtFx) {
			bare := filepath.Join(f.base, "bare.git")
			writeF(t, filepath.Join(bare, "HEAD"), "ref: refs/heads/main\n")
			os.MkdirAll(filepath.Join(bare, "objects"), 0o755)
			os.MkdirAll(filepath.Join(bare, "refs"), 0o755)
			admin := filepath.Join(bare, "worktrees", "T-1")
			writeF(t, filepath.Join(admin, "commondir"), "../..\n")
			writeF(t, filepath.Join(admin, "gitdir"), filepath.Join(f.wt, ".git")+"\n")
			f.gitfile(t, "gitdir: "+admin+"\n")
		}},
		{"submodule style gitfile", func(t *testing.T, f wtFx) {
			os.MkdirAll(filepath.Join(f.common, "modules", "sub"), 0o755)
			f.gitfile(t, "gitdir: "+filepath.Join(f.common, "modules", "sub")+"\n")
		}},
		{"gitfile over 4 KiB", func(t *testing.T, f wtFx) {
			f.gitfile(t, "gitdir: "+f.admin+strings.Repeat("/.", 3000)+"\n")
		}},
		{"gitfile with NUL", func(t *testing.T, f wtFx) { f.gitfile(t, "gitdir: "+f.admin+"\x00\n") }},
		{"two-line gitfile", func(t *testing.T, f wtFx) { f.gitfile(t, "gitdir: "+f.admin+"\nsecond\n") }},
		{"no gitdir prefix", func(t *testing.T, f wtFx) { f.gitfile(t, f.admin+"\n") }},
		{"gitdir without space", func(t *testing.T, f wtFx) { f.gitfile(t, "gitdir:"+f.admin+"\n") }},
		{"empty gitfile", func(t *testing.T, f wtFx) { f.gitfile(t, "") }},
		{"admin dir name not under worktrees", func(t *testing.T, f wtFx) {
			f.gitfile(t, "gitdir: "+filepath.Join(f.common, "T-1")+"\n")
		}},
		{"admin dir missing", func(t *testing.T, f wtFx) { os.RemoveAll(f.admin) }},
		{"common dir without objects", func(t *testing.T, f wtFx) { os.RemoveAll(filepath.Join(f.common, "objects")) }},
		{"symlinked .git", func(t *testing.T, f wtFx) {
			real := filepath.Join(f.base, "realgitfile")
			os.Rename(filepath.Join(f.wt, ".git"), real)
			symlinkOrSkip(t, real, filepath.Join(f.wt, ".git"))
		}},
		{"symlinked admin dir", func(t *testing.T, f wtFx) {
			real := filepath.Join(f.common, "worktrees", "real")
			os.Rename(f.admin, real)
			symlinkOrSkip(t, real, f.admin)
		}},
		{"symlinked common dir", func(t *testing.T, f wtFx) {
			real := filepath.Join(f.main, "realgit")
			os.Rename(f.common, real)
			symlinkOrSkip(t, real, f.common)
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := newWT(t)
			c.mut(t, f)
			refuse(t, f)
		})
	}
}

// R26: `.git` as a directory is not a linked worktree even if its content looks like one.
func TestR26_GitDirectoryIsNotAWorktree(t *testing.T) {
	f := newWT(t)
	os.Remove(filepath.Join(f.wt, ".git"))
	os.MkdirAll(filepath.Join(f.wt, ".git"), 0o755)
	if LinkedWorktree(f.wt) != nil {
		t.Error(".git directory accepted")
	}
}

// R26 hardening: the gitfile read refuses anything that is not a plain file (a FIFO would block).
func TestR26_ReadSmallRejectsNonRegular(t *testing.T) {
	d := realTemp(t)
	if _, ok := readSmall(d); ok {
		t.Error("directory accepted")
	}
	real := filepath.Join(d, "f")
	writeF(t, real, "gitdir: x\n")
	if s, ok := readSmall(real); !ok || s != "gitdir: x" {
		t.Errorf("%q %v", s, ok)
	}
	symlinkOrSkip(t, real, filepath.Join(d, "l"))
	if _, ok := readSmall(filepath.Join(d, "l")); ok {
		t.Error("symlink accepted")
	}
	// not-run: the os.SameFile swap-between-Lstat-and-Open guard needs a race; not deterministic.
}
