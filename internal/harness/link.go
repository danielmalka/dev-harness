package harness

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// canon is the registry key: Abs+Clean, then EvalSymlinks when the path exists, so a symlinked
// path and its physical path share one entry.
// ponytail: paths are not case-folded (names are, via EqualFold); a case-insensitive filesystem may need two links, fold here if that bites.
func canon(p string) (string, error) {
	a, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	a = filepath.Clean(a)
	if r, err := filepath.EvalSymlinks(a); err == nil {
		return r, nil
	}
	return a, nil
}

// Link registers repo in the registry (name defaults to the folder name) and, when the repo
// has no internal .harness/, creates <home>/projects/<name>/ (reusing it if present). With
// off it removes the entry and keeps the folder. Nothing is ever created inside the repo.
// The whole load-modify-save runs under <home>/config.lock.
func Link(home, repo, name string, off bool) (Resolved, error) {
	plain, err := filepath.Abs(repo)
	if err != nil {
		return Resolved{}, err
	}
	plain = filepath.Clean(plain)
	repo, err = canon(repo)
	if err != nil {
		return Resolved{}, err
	}
	unlock, err := lock(home)
	if err != nil {
		return Resolved{}, err
	}
	defer unlock()
	c, err := LoadConfig(home)
	if err != nil {
		return Resolved{}, err
	}
	if off { // works on a stale path: it only drops the entry
		delete(c.Projects, repo)
		delete(c.Projects, plain) // an entry saved under the symlinked spelling
		if err := SaveConfig(home, c); err != nil {
			return Resolved{}, err
		}
		return resolveIn(home, repo, c), nil
	}
	if !isDir(repo) {
		return Resolved{}, fmt.Errorf("%s is not a directory", repo)
	}
	if name == "" {
		name = filepath.Base(repo)
	}
	if !validName(name) {
		return Resolved{}, fmt.Errorf("invalid project name %q: use a single folder name (pass an explicit name)", name)
	}
	for p, n := range c.Projects {
		if !filepath.IsAbs(p) || p != filepath.Clean(p) { // hand-edited key: never stat it
			continue
		}
		if !strings.EqualFold(n, name) || p == repo {
			continue
		}
		if cp, err := canon(p); err == nil && cp == repo {
			delete(c.Projects, p) // same repository under another spelling: replace
			continue
		}
		if _, serr := os.Stat(p); !errors.Is(serr, os.ErrNotExist) {
			return Resolved{}, fmt.Errorf("name %q is already used by %s: pass an explicit name (dh link <name>)", name, p)
		}
		delete(c.Projects, p) // stale entry: its repository is gone, relink
	}
	if !isDir(filepath.Join(repo, ".harness")) {
		d := filepath.Join(home, "projects", name)
		if err := os.MkdirAll(d, dirPerm); err != nil {
			return Resolved{}, err
		}
		if !isDir(d) {
			return Resolved{}, errors.New(d + " exists and is not a directory")
		}
		if err := tighten(d, dirPerm); err != nil {
			return Resolved{}, err
		}
	}
	c.Projects[repo] = name
	if err := SaveConfig(home, c); err != nil {
		return Resolved{}, err
	}
	return resolveIn(home, repo, c), nil
}
