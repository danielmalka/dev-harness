package harness

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// LinkedWT describes a valid linked git worktree. All paths are canonical.
type LinkedWT struct {
	Main      string `json:"main"`      // main checkout (parent of the common dir)
	Gitdir    string `json:"gitdir"`    // admin dir <common>/worktrees/<name>, not the gitfile
	Commondir string `json:"commondir"` // the common dir (<main>/.git)
}

const maxGitfile = 4096

// readSmall reads a regular, non-symlink file of at most 4 KiB with no NUL, and returns its
// content without one trailing newline. Any other shape returns ok=false.
func readSmall(p string) (string, bool) {
	fi, err := os.Lstat(p)
	if err != nil || !fi.Mode().IsRegular() || fi.Size() > maxGitfile {
		return "", false
	}
	// O_NONBLOCK: opening a FIFO swapped in after Lstat would otherwise block until a writer appears.
	f, err := os.OpenFile(p, os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return "", false
	}
	defer f.Close()
	// Reject a swap between Lstat and Open (portable; no O_NOFOLLOW for windows builds).
	if fi2, err := f.Stat(); err != nil || !fi2.Mode().IsRegular() || !os.SameFile(fi, fi2) {
		return "", false
	}
	b, err := io.ReadAll(io.LimitReader(f, maxGitfile+1))
	if err != nil || len(b) > maxGitfile || bytes.IndexByte(b, 0) >= 0 {
		return "", false
	}
	s := strings.TrimSuffix(strings.TrimSuffix(string(b), "\n"), "\r")
	if s == "" || strings.ContainsAny(s, "\r\n") {
		return "", false
	}
	return s, true
}

func notSymlink(p string) bool {
	fi, err := os.Lstat(p)
	return err == nil && fi.Mode()&os.ModeSymlink == 0
}

// rawJoin joins base and p WITHOUT Clean: "x/.." must reach the OS as written, because the
// kernel (and git) resolve it as the parent of x's symlink target, not as text.
func rawJoin(base, p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return base + string(filepath.Separator) + p
}

// resolveRaw resolves a raw target the way the OS does. Any error (missing path, loop) is
// not-a-worktree: there is no fallback to a cleaned path.
func resolveRaw(base, p string) (string, bool) {
	r, err := filepath.EvalSymlinks(rawJoin(base, p))
	return r, err == nil
}

// LinkedWorktree returns the worktree description when dir holds a valid linked-worktree
// gitfile (ADR-008 D3), else nil. It reads files only and never runs git. The gitfile is
// attacker-controlled input: every check below must hold or the result is nil.
func LinkedWorktree(dir string) *LinkedWT {
	d, err := canon(dir)
	if err != nil {
		return nil
	}
	gf := filepath.Join(d, ".git")
	line, ok := readSmall(gf)
	if !ok || !strings.HasPrefix(line, "gitdir: ") {
		return nil
	}
	raw := rawJoin(d, strings.TrimPrefix(line, "gitdir: "))
	// Split the raw text (no Clean) into <rawCommon>/worktrees/<name>; refuse symlinks on each.
	i := strings.LastIndexAny(raw, `/\`)
	if i < 0 {
		return nil
	}
	name := raw[i+1:]
	j := strings.LastIndexAny(raw[:i], `/\`)
	if j < 0 || raw[j+1:i] != "worktrees" || !validName(name) {
		return nil
	}
	rawCommon := raw[:j]
	if rawCommon == "" {
		rawCommon = string(filepath.Separator)
	}
	if !notSymlink(raw) || !notSymlink(raw[:i]) || !notSymlink(rawCommon) {
		return nil
	}
	admin, ok := resolveRaw(d, strings.TrimPrefix(line, "gitdir: "))
	if !ok || !isDir(admin) {
		return nil
	}
	common := filepath.Dir(filepath.Dir(admin))
	if filepath.Base(admin) != name || filepath.Base(filepath.Dir(admin)) != "worktrees" ||
		filepath.Base(common) != ".git" {
		return nil
	}
	if cd, ok := readSmall(filepath.Join(admin, "commondir")); !ok {
		return nil
	} else if c, ok := resolveRaw(admin, cd); !ok || c != common {
		return nil
	}
	if bp, ok := readSmall(filepath.Join(admin, "gitdir")); !ok {
		return nil
	} else if c, ok := resolveRaw(admin, bp); !ok || c != gf {
		return nil
	}
	if fi, err := os.Stat(filepath.Join(common, "HEAD")); err != nil || !fi.Mode().IsRegular() ||
		!isDir(filepath.Join(common, "objects")) || !isDir(filepath.Join(common, "refs")) {
		return nil
	}
	return &LinkedWT{Main: filepath.Dir(common), Gitdir: admin, Commondir: common}
}
