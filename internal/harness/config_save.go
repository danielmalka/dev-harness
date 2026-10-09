package harness

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var plainRe = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)

func quoteVal(s string) string {
	if plainRe.MatchString(s) {
		return s
	}
	return strconv.Quote(s)
}

func nl(s string) string {
	if s != "" && !strings.HasSuffix(s, "\n") {
		return s + "\n"
	}
	return s
}

// SaveConfig writes <home>/config.yaml atomically (temp + rename, 0600). It re-parses the
// existing file first: an invalid file is an error and is never overwritten. Lines the reader
// does not model, and modeled keys whose value did not change, are kept byte for byte.
func SaveConfig(home string, c Config) error {
	data, err := readConfig(home)
	if err != nil {
		return err
	}
	old, segs, err := parseConfig(data)
	if err != nil {
		return err
	}
	var out strings.Builder
	have := map[string]bool{}
	for _, s := range segs {
		have[s.key] = true
		raw := strings.Join(s.lines, "")
		switch s.key {
		case "language":
			if c.Language != old.Language {
				raw = scalarLine("language", c.Language)
			}
		case "mode":
			if c.Mode != old.Mode {
				raw = scalarLine("mode", c.Mode)
			}
		case "dashboard":
			if c.Sprites != old.Sprites {
				raw = replaceSprites(s, c.Sprites)
			}
		case "reviewers":
			if c.ReviewersRaw != old.ReviewersRaw {
				raw = nl(c.ReviewersRaw)
			}
		case "projects":
			if c.Projects != nil { // nil = not modeled by the caller: keep the block as is
				raw = rewriteProjects(s, c.Projects)
			}
		}
		out.WriteString(nl(raw))
	}
	tail := ""
	if !have["language"] {
		tail += scalarLine("language", c.Language)
	}
	if !have["mode"] {
		tail += scalarLine("mode", c.Mode)
	}
	if !have["dashboard"] && c.Sprites != "" {
		tail += "dashboard:\n  sprites: " + quoteVal(c.Sprites) + "\n"
	}
	if !have["reviewers"] {
		tail += nl(c.ReviewersRaw)
	}
	if !have["projects"] {
		tail += rewriteProjects(&seg{key: "projects", lines: []string{"projects:\n"}}, c.Projects)
	}
	return writeAtomic(home, out.String()+tail)
}

func scalarLine(key, v string) string {
	if v == "" {
		return ""
	}
	return key + ": " + quoteVal(v) + "\n"
}

func replaceSprites(s *seg, v string) string {
	var b strings.Builder
	b.WriteString(nl(s.lines[0]))
	done := false
	for _, l := range s.lines[1:] {
		if strings.HasPrefix(strings.TrimSpace(l), "sprites:") {
			if v != "" {
				b.WriteString("  sprites: " + quoteVal(v) + "\n")
			}
			done = true
			continue
		}
		b.WriteString(nl(l))
	}
	if !done && v != "" {
		b.WriteString("  sprites: " + quoteVal(v) + "\n")
	}
	return b.String()
}

// rewriteProjects keeps the header's comments, unchanged entries and foreign lines in place,
// rewrites renamed entries, drops removed ones and appends new ones after the last entry.
func rewriteProjects(s *seg, want map[string]string) string {
	out := []string{"projects:\n"}
	if h := s.lines[0]; strings.TrimSpace(eol(h)) != "projects:" && isBlank(strings.TrimPrefix(strings.TrimSpace(eol(h)), "projects:")) {
		out[0] = nl(h) // header with only a comment: keep verbatim
	}
	done := map[string]bool{}
	insertAt := 1
	for _, l := range s.lines[1:] {
		if isBlank(l) {
			out = append(out, nl(l))
			continue
		}
		p, name, err := entry(l, 0)
		if err != nil {
			out = append(out, nl(l))
			continue
		}
		w, ok := want[p]
		if !ok {
			continue
		}
		done[p] = true
		if w == name {
			out = append(out, nl(l))
		} else {
			out = append(out, "  "+strconv.Quote(p)+": "+quoteVal(w)+"\n")
		}
		insertAt = len(out)
	}
	var fresh []string
	for p := range want {
		if !done[p] {
			fresh = append(fresh, p)
		}
	}
	sort.Strings(fresh)
	add := make([]string, 0, len(fresh))
	for _, p := range fresh {
		add = append(add, "  "+strconv.Quote(p)+": "+quoteVal(want[p])+"\n")
	}
	out = append(out[:insertAt], append(add, out[insertAt:]...)...)
	return strings.Join(out, "")
}

func writeAtomic(home, content string) error {
	if err := os.MkdirAll(home, dirPerm); err != nil {
		return err
	}
	f, err := os.CreateTemp(home, ".config-*.tmp")
	if err != nil {
		return err
	}
	tmp := f.Name()
	_, err = f.WriteString(content)
	if err == nil {
		err = f.Chmod(filePerm)
	}
	if err == nil {
		err = f.Sync()
	}
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Rename(tmp, filepath.Join(home, "config.yaml"))
	}
	if err != nil {
		os.Remove(tmp)
	}
	return err
}
