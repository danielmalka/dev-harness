// Package harness owns the machine-global dev-harness folder: <home>, its config.yaml,
// the project registry and the single resolver behind `dh harness-path`, `dh link` and `dh projects`.
package harness

import (
	"errors"
	"os"
	"path/filepath"
)

const (
	dirPerm  = 0o700
	filePerm = 0o600
)

// Home is $DH_HOME when set, else $HOME/.harness, $USERPROFILE/.harness, <UserHomeDir>/.harness. It never creates anything and is
// empty when neither can be determined.
func Home() string {
	if v := os.Getenv("DH_HOME"); v != "" {
		if abs, err := filepath.Abs(v); err == nil {
			return abs
		}
		return filepath.Clean(v)
	}
	// Same order as the TypeScript mod on every OS: HOME, then USERPROFILE, then the Go fallback.
	for _, v := range []string{"HOME", "USERPROFILE"} {
		if h := os.Getenv(v); h != "" {
			return filepath.Join(h, ".harness")
		}
	}
	if h, err := os.UserHomeDir(); err == nil && h != "" {
		return filepath.Join(h, ".harness")
	}
	return ""
}

// SessionsDir and DashboardDir are empty when Home is unknown, never a relative path.
func SessionsDir() string  { return sub("sessions") }
func DashboardDir() string { return sub("dashboard") }

func sub(name string) string {
	if h := Home(); h != "" {
		return filepath.Join(h, name)
	}
	return ""
}

// EnsureHome creates <home>, projects/, sessions/, dashboard/ and an empty config.yaml when
// missing. It is idempotent and never touches a repository. It returns the home path.
func EnsureHome() (string, error) {
	home := Home()
	if home == "" {
		return "", errors.New("cannot determine the harness home: set DH_HOME, HOME or USERPROFILE")
	}
	if r, err := filepath.EvalSymlinks(home); err == nil && r != home { // <home> itself may be a symlink
		if err := tighten(r, dirPerm); err != nil {
			return "", err
		}
	}
	for _, d := range []string{home, "projects", "sessions", "dashboard"} {
		if d != home {
			d = filepath.Join(home, d)
		}
		if err := os.MkdirAll(d, dirPerm); err != nil {
			return "", err
		}
		if err := tighten(d, dirPerm); err != nil { // e.g. sessions/ made by the mod with a wider mode
			return "", err
		}
	}
	// direct children of projects/ only; the records inside are not touched
	if ents, err := os.ReadDir(filepath.Join(home, "projects")); err == nil {
		for _, e := range ents {
			if e.IsDir() {
				if err := tighten(filepath.Join(home, "projects", e.Name()), dirPerm); err != nil {
					return "", err
				}
			}
		}
	}
	f, err := os.OpenFile(filepath.Join(home, "config.yaml"), os.O_WRONLY|os.O_CREATE|os.O_EXCL, filePerm)
	if err == nil {
		_, werr := f.WriteString("projects:\n")
		if cerr := f.Close(); werr == nil {
			werr = cerr
		}
		if werr != nil {
			return "", werr
		}
	} else if !errors.Is(err, os.ErrExist) {
		return "", err
	} else if err := tighten(filepath.Join(home, "config.yaml"), filePerm); err != nil {
		return "", err
	}
	return home, nil
}

// tighten chmods p to perm when it is group/other accessible. Symlinks are left alone. A failed
// chmod (read-only or foreign-owned path) is skipped, never an error: creation errors are what fail.
func tighten(p string, perm os.FileMode) error {
	if fi, err := os.Lstat(p); err == nil && fi.Mode()&os.ModeSymlink == 0 && fi.Mode().Perm()&0o077 != 0 {
		_ = os.Chmod(p, perm)
	}
	return nil
}
