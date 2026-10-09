package harness

import (
	"os"
	"path/filepath"
)

// PdocsDir is <home>/projects/<name>/pdocs, or "" for an empty home or invalid name. Creates nothing.
func PdocsDir(home, name string) string {
	if home == "" || !validName(name) {
		return ""
	}
	return filepath.Join(home, "projects", name, "pdocs")
}

// IsHarness: dir is a directory that is not a pdocs-only folder (PRD-017 R2). A pdocs-only folder
// has pdocs/ and none of the six harness records; an empty folder still counts as harness.
func IsHarness(dir string) bool {
	if !isDir(dir) {
		return false
	}
	for _, n := range []string{"project.yaml", "MEMORY.md", "EPOCHAL.md", "RISKS.md", "tasks", "prd"} {
		if _, err := os.Stat(filepath.Join(dir, n)); err == nil {
			return true
		}
	}
	return !isDir(filepath.Join(dir, "pdocs"))
}
