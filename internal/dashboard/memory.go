package dashboard

import (
	"net/http"
	"strings"

	"github.com/danielmalka/dev-harness/internal/harness"
)

// memoryFiles is the closed list of files /api/memory can serve; nothing from the query reaches a path.
var memoryFiles = map[string]string{"memory": "MEMORY.md", "risks": "RISKS.md"}

// memoryH serves MEMORY.md or RISKS.md of a registered project. Every refusal is the bare 404 (no echo).
func memoryH(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		name, file := q.Get("project"), q.Get("file")
		base, ok := memoryFiles[file]
		if !ok || name == "" || cfg.Home == "" || strings.ContainsAny(name, `/\`) || strings.Contains(name, "..") {
			http.NotFound(w, r)
			return
		}
		projects, _ := harness.Projects(cfg.Home)
		for _, p := range projects {
			if p.Name != name {
				continue
			}
			if symlinked(p.Harness) {
				break
			}
			b, err := readCappedIn(p.Harness, base) // first match in registry order
			if err != nil {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Header().Set("Cache-Control", "no-store")
			_, _ = w.Write(b)
			return
		}
		http.NotFound(w, r)
	}
}
