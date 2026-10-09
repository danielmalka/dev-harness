package dashboard

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/danielmalka/dev-harness/internal/harness"
	"github.com/danielmalka/dev-harness/internal/snapshot"
)

const pdocsCap = 8 << 20

// pdocsCSP is only for /pdocs/*: sandbox without allow-same-origin gives the page an opaque origin,
// so a script from the CDN cannot read /api/state or /api/memory (PRD-017 R10).
const pdocsCSP = "default-src 'self' 'unsafe-inline'; img-src 'self' data:; script-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net/npm/mermaid@11.17.2/dist/mermaid.min.js; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src https://fonts.gstatic.com; connect-src 'self'; sandbox allow-scripts allow-popups"

var pdocsTypes = map[string]string{".html": "text/html; charset=utf-8", ".svg": "image/svg+xml", ".png": "image/png"}

// pdocsH serves GET /pdocs/<name>/<file>. Every refusal is the bare 404 with the route CSP (no echo, no warnings).
func pdocsH(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", pdocsCSP)
		name, file, ok := strings.Cut(strings.TrimPrefix(r.URL.Path, "/pdocs/"), "/")
		if file == "" {
			file = "index.html"
		}
		ct := pdocsTypes[filepath.Ext(file)]
		if !ok || name == "" || cfg.Home == "" || ct == "" || strings.ContainsAny(file, `/\:`) || strings.Contains(file, "..") {
			http.NotFound(w, r)
			return
		}
		projects, _ := harness.Projects(cfg.Home)
		for _, p := range projects {
			if p.Name != name {
				continue
			}
			b, err := snapshot.ReadCappedIn(harness.PdocsDir(cfg.Home, p.Name), file, pdocsCap)
			if err != nil {
				break
			}
			w.Header().Set("Content-Type", ct)
			w.Header().Set("Cache-Control", "no-cache")
			_, _ = w.Write(b)
			return
		}
		http.NotFound(w, r)
	}
}

// docsFlag reports whether <pdocs>/index.html is a regular file within the cap. Lstat only: never reads it.
// The warning (no absolute path) is set for a symlink, a non-regular file or an oversized one.
func docsFlag(home, name string) (bool, string) {
	fi, err := os.Lstat(filepath.Join(harness.PdocsDir(home, name), "index.html"))
	switch {
	case err != nil:
		return false, ""
	case fi.Mode()&os.ModeSymlink != 0:
		return false, fmt.Sprintf("%s: pdocs/index.html is a symlink: no documentation link.", name)
	case !fi.Mode().IsRegular():
		return false, fmt.Sprintf("%s: pdocs/index.html is not a regular file: no documentation link.", name)
	case fi.Size() > pdocsCap:
		return false, fmt.Sprintf("%s: pdocs/index.html is over 8 MiB: no documentation link.", name)
	}
	return true, ""
}
