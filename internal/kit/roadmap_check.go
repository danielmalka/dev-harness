package kit

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// checkRoadmaps keeps docs/roadmap.html and docs/en/roadmap.html in step with
// CHANGELOG.md and with each other (PRD-009, R1-R8). It reads the HTML with
// regular expressions anchored on the exact markup those two files use; if the
// markup changes, the check fails loudly (chip or range "not found") instead
// of passing silently.
//
// Read from each roadmap: the date chip, the status chip ("<v> · vivo" /
// "<v> · live"), the path chip, the "CHANGELOG.md — releases 0.1.0 a|to X"
// line, and the <h3> cards between <h2 id="entregue"> and the next <h2>. A
// card whose title starts with "N.N.N · " is a version card; any other card is
// a detail card and is ignored. The <h3> title may contain <code>, so only the
// prefix is matched. An <h3> in Entregue whose title (tags stripped, joined
// across lines) starts like a version but not in the strict "<h3>N.N.N · "
// form (attribute, wrapped title, leading space, NBSP, em dash) is an error,
// never silently ignored. Language words are enforced: pt "data"/"a", en
// "date"/"to".
//
// Skipped only when both roadmaps are absent (generated package, consumer
// project). Every failure is "<file>:<line>: expected ..., found ...".
var (
	changelogHeading = regexp.MustCompile(`^## (\d+\.\d+\.\d+) — (\d{4}-\d{2}-\d{2})\s*$`)
	changelogAnyRel  = regexp.MustCompile(`^##\s*\[?v?\d`)
	roadmapDateChip  = regexp.MustCompile(`<span class="chip">(data|date) <b>([^<]*)</b>`)
	roadmapStatus    = regexp.MustCompile(`<span class="chip">status <b>([^<]*)</b>`)
	roadmapPathChip  = regexp.MustCompile(`<span class="chip">path <b>([^<]*)</b>`)
	roadmapRange     = regexp.MustCompile(`<code>CHANGELOG\.md</code> — releases (\d+\.\d+\.\d+) (a|to) (\d+\.\d+\.\d+)`)
	roadmapCard      = regexp.MustCompile(`<h3>(\d+\.\d+\.\d+) · `)
	roadmapH3Open    = regexp.MustCompile(`<h3\b`)
	roadmapTag       = regexp.MustCompile(`<[^>]*>`)
	roadmapVersionAt = regexp.MustCompile(`^\d+\.\d+\.\d+`)
	roadmapH2        = regexp.MustCompile(`<h2[ >]`)
	roadmapEntregue  = regexp.MustCompile(`<h2 id="entregue"`)
)

type roadmapValue struct {
	text string
	word string // language word (data/date, a/to) where the value has one
	line int    // 0 when not found
}

type roadmapCardRef struct {
	version string
	line    int
}

type roadmapBadCard struct {
	found string
	line  int
}

type roadmapDoc struct {
	rel                  string
	date, status, path   roadmapValue
	rangeStart, rangeEnd roadmapValue
	entregueLine         int
	cards                []roadmapCardRef
	badCards             []roadmapBadCard
}

func parseRoadmap(rel, text string) roadmapDoc {
	doc := roadmapDoc{rel: rel}
	inEntregue := false
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		n := i + 1
		if doc.date.line == 0 {
			if m := roadmapDateChip.FindStringSubmatch(line); m != nil {
				doc.date = roadmapValue{text: m[2], word: m[1], line: n}
			}
		}
		if doc.status.line == 0 {
			if m := roadmapStatus.FindStringSubmatch(line); m != nil {
				doc.status = roadmapValue{text: m[1], line: n}
			}
		}
		if doc.path.line == 0 {
			if m := roadmapPathChip.FindStringSubmatch(line); m != nil {
				doc.path = roadmapValue{text: m[1], line: n}
			}
		}
		if doc.rangeEnd.line == 0 {
			if m := roadmapRange.FindStringSubmatch(line); m != nil {
				doc.rangeStart = roadmapValue{text: m[1], line: n}
				doc.rangeEnd = roadmapValue{text: m[3], word: m[2], line: n}
			}
		}
		if roadmapEntregue.MatchString(line) {
			doc.entregueLine, inEntregue = n, true
			continue
		}
		if inEntregue && roadmapH2.MatchString(line) {
			inEntregue = false
		}
		if !inEntregue || !roadmapH3Open.MatchString(line) {
			continue
		}
		if m := roadmapCard.FindStringSubmatch(line); m != nil {
			doc.cards = append(doc.cards, roadmapCardRef{m[1], n})
			continue
		}
		// ponytail: wrapped titles are read for at most 8 lines, raise if a title ever wraps further.
		// Not the strict form: read the whole title (up to </h3>, at most 8
		// lines), strip tags, and flag it if it still starts like a version.
		title := line[roadmapH3Open.FindStringIndex(line)[0]:]
		for j := i; !strings.Contains(lines[j], "</h3>") && j+1 < len(lines) && j < i+8; j++ {
			title += " " + lines[j+1]
		}
		if k := strings.Index(title, ">"); k >= 0 {
			title = title[k+1:]
		}
		if k := strings.Index(title, "</h3>"); k >= 0 {
			title = title[:k]
		}
		title = roadmapTag.ReplaceAllString(title, "")
		probe := strings.NewReplacer("&nbsp;", " ", "\u00a0", " ").Replace(title)
		if roadmapVersionAt.MatchString(strings.TrimSpace(probe)) {
			r := []rune(title)
			if len(r) > 24 {
				r = r[:24]
			}
			doc.badCards = append(doc.badCards, roadmapBadCard{string(r), n})
		}
	}
	return doc
}

// semverLess compares N.N.N numerically (0.9.0 < 0.10.0).
// Atoi errors are discarded: the callers only pass N.N.N matched by regex.
func semverLess(a, b string) bool {
	pa, pb := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < 3 && i < len(pa) && i < len(pb); i++ {
		x, _ := strconv.Atoi(pa[i])
		y, _ := strconv.Atoi(pb[i])
		if x != y {
			return x < y
		}
	}
	return false
}

func roadmapVersions(d roadmapDoc) []string {
	out := make([]string, len(d.cards))
	for i, c := range d.cards {
		out[i] = c.version
	}
	return out
}

func checkRoadmaps(root string, errors *[]string) {
	files := []struct{ rel, word, dateWord, toWord string }{
		{"docs/roadmap.html", "vivo", "data", "a"},
		{"docs/en/roadmap.html", "live", "date", "to"},
	}
	present := 0
	for _, f := range files {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(f.rel))); err == nil {
			present++
		}
	}
	if present == 0 {
		return
	}
	add := func(format string, args ...any) { *errors = append(*errors, fmt.Sprintf(format, args...)) }

	docs := make([]roadmapDoc, 0, 2)
	for _, f := range files {
		text, err := readText(filepath.Join(root, filepath.FromSlash(f.rel)))
		if err != nil {
			add("%s:1: expected the roadmap file to exist next to its pair, found it missing or unreadable", f.rel)
			return
		}
		docs = append(docs, parseRoadmap(f.rel, text))
	}

	changelog, err := readText(filepath.Join(root, "CHANGELOG.md"))
	if err != nil {
		add("CHANGELOG.md:1: expected CHANGELOG.md to exist because the roadmaps do, found it missing or unreadable")
		return
	}
	var versions []string
	dates := map[string]string{}
	for i, line := range strings.Split(changelog, "\n") {
		m := changelogHeading.FindStringSubmatch(line)
		switch {
		case m == nil && changelogAnyRel.MatchString(line):
			add("CHANGELOG.md:%d: expected a heading like \"## 0.1.0 — 2026-01-01\", found %q", i+1, line)
		case m == nil:
		case dates[m[1]] != "":
			add("CHANGELOG.md:%d: expected one heading per release, found %s again", i+1, m[1])
		default:
			versions = append(versions, m[1])
			dates[m[1]] = m[2]
		}
	}
	if len(versions) == 0 {
		add("CHANGELOG.md:1: expected a heading like \"## 0.1.0 — 2026-01-01\", found none")
		return
	}
	latest, latestDate := versions[0], dates[versions[0]]

	// clean[i] is true when file i passed R2-R4, so pt/en are compared on
	// those values only when both files are already right.
	clean := make([]bool, len(docs))
	for i, d := range docs {
		f := files[i]
		clean[i] = true
		check := func(name string, v roadmapValue, want string) {
			switch {
			case v.line == 0:
				add("%s:1: expected %s %q, found none", d.rel, name, want)
				clean[i] = false
			case v.text != want:
				add("%s:%d: expected %s %q, found %q", d.rel, v.line, name, want, v.text)
				clean[i] = false
			}
		}
		check("date chip", d.date, latestDate)
		check("status chip", d.status, latest+" · "+f.word)
		check("Sources range end", d.rangeEnd, latest)
		if d.date.line != 0 && d.date.word != f.dateWord {
			add("%s:%d: expected date chip word %q, found %q", d.rel, d.date.line, f.dateWord, d.date.word)
			clean[i] = false
		}
		if d.rangeEnd.line != 0 && d.rangeEnd.word != f.toWord {
			add("%s:%d: expected Sources range word %q, found %q", d.rel, d.rangeEnd.line, f.toWord, d.rangeEnd.word)
			clean[i] = false
		}
		if d.path.line == 0 {
			add("%s:1: expected path chip, found none", d.rel)
		}

		if d.entregueLine == 0 {
			add("%s:1: expected <h2 id=\"entregue\">, found none", d.rel)
			continue
		}
		for _, b := range d.badCards {
			add("%s:%d: expected \"<h3>N.N.N · \", found %q", d.rel, b.line, b.found)
		}
		count := map[string]int{}
		for _, c := range d.cards {
			count[c.version]++
			if count[c.version] == 2 {
				add("%s:%d: expected exactly one card for %s, found a duplicate", d.rel, c.line, c.version)
			}
		}
		inChangelog := map[string]bool{}
		for _, v := range versions {
			inChangelog[v] = true
			if count[v] == 0 {
				add("%s:%d: expected a Delivered card for %s, found none", d.rel, d.entregueLine, v)
			}
		}
		for _, c := range d.cards {
			if !inChangelog[c.version] {
				add("%s:%d: expected no card for %s (no CHANGELOG.md heading), found one", d.rel, c.line, c.version)
			}
		}
		for j := 1; j < len(d.cards); j++ {
			prev, cur := d.cards[j-1], d.cards[j]
			if semverLess(prev.version, cur.version) {
				add("%s:%d: expected newest first, found card %s after card %s (line %d)", d.rel, cur.line, cur.version, prev.version, prev.line)
			}
		}
	}

	// R7: the pair comparison skips any value one side lacks (already
	// reported above) so en is never blamed for what pt is missing.
	pt, en := docs[0], docs[1]
	if pt.entregueLine != 0 && en.entregueLine != 0 {
		ptV, enV := strings.Join(roadmapVersions(pt), ", "), strings.Join(roadmapVersions(en), ", ")
		if ptV != enV {
			add("%s:%d: expected the same version cards as docs/roadmap.html, found pt [%s] and en [%s]", en.rel, en.entregueLine, ptV, enV)
		}
	}
	same := func(name string, p, e roadmapValue) {
		if p.line != 0 && e.line != 0 && p.text != e.text {
			add("%s:%d: expected %s %q as in docs/roadmap.html, found %q", en.rel, e.line, name, p.text, e.text)
		}
	}
	same("path chip", pt.path, en.path)
	same("Sources range start", pt.rangeStart, en.rangeStart)
	// Date, status version and range end need no pt/en comparison: R2-R4 pin each
	// file to the CHANGELOG value, so equality across files follows transitively.
}
