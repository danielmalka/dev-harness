package kit

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// checkRoadmaps keeps docs/roadmap.html in step with CHANGELOG.md (PRD-009, R1-R6; ADR-011 dropped the English mirror). It reads the HTML with regular expressions
// anchored on the exact markup that file uses; if the markup changes, the check fails loudly (chip or range "not found") instead of passing silently.
//
// Read from the roadmap: the date chip, the status chip ("<v> · vivo"), the path chip, the "CHANGELOG.md — releases 0.1.0 a X" line, and the <h3> cards between
// <h2 id="entregue"> and the next <h2>. A card whose title starts with "N.N.N · " is a version card; any other card is a detail card and is ignored. The <h3> title may
// contain <code>, so only the prefix is matched. An <h3> in Entregue whose title (tags stripped, joined across lines) starts like a version but not in the strict
// "<h3>N.N.N · " form (attribute, wrapped title, leading space, NBSP, em dash) is an error, never silently ignored.
//
// Skipped only when the roadmap is absent (generated package, consumer project). Every failure is "<file>:<line>: esperado ..., achei ...".
var (
	changelogHeading = lazyRe(`^## (\d+\.\d+\.\d+) — (\d{4}-\d{2}-\d{2})\s*$`)
	changelogAnyRel  = lazyRe(`^##\s*\[?v?\d`)
	roadmapDateChip  = lazyRe(`<span class="chip">data <b>([^<]*)</b>`)
	roadmapStatus    = lazyRe(`<span class="chip">status <b>([^<]*)</b>`)
	roadmapPathChip  = lazyRe(`<span class="chip">path <b>([^<]*)</b>`)
	roadmapRange     = lazyRe(`<code>CHANGELOG\.md</code> — releases (\d+\.\d+\.\d+) a (\d+\.\d+\.\d+)`)
	roadmapCard      = lazyRe(`<h3>(\d+\.\d+\.\d+) · `)
	roadmapH3Open    = lazyRe(`<h3\b`)
	roadmapTag       = lazyRe(`<[^>]*>`)
	roadmapVersionAt = lazyRe(`^\d+\.\d+\.\d+`)
	roadmapH2        = lazyRe(`<h2[ >]`)
	roadmapEntregue  = lazyRe(`<h2 id="entregue"`)
)

type roadmapValue struct {
	text string
	line int // 0 when not found
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
			if m := roadmapDateChip().FindStringSubmatch(line); m != nil {
				doc.date = roadmapValue{text: m[1], line: n}
			}
		}
		if doc.status.line == 0 {
			if m := roadmapStatus().FindStringSubmatch(line); m != nil {
				doc.status = roadmapValue{text: m[1], line: n}
			}
		}
		if doc.path.line == 0 {
			if m := roadmapPathChip().FindStringSubmatch(line); m != nil {
				doc.path = roadmapValue{text: m[1], line: n}
			}
		}
		if doc.rangeEnd.line == 0 {
			if m := roadmapRange().FindStringSubmatch(line); m != nil {
				doc.rangeStart = roadmapValue{text: m[1], line: n}
				doc.rangeEnd = roadmapValue{text: m[2], line: n}
			}
		}
		if roadmapEntregue().MatchString(line) {
			doc.entregueLine, inEntregue = n, true
			continue
		}
		if inEntregue && roadmapH2().MatchString(line) {
			inEntregue = false
		}
		if !inEntregue || !roadmapH3Open().MatchString(line) {
			continue
		}
		if m := roadmapCard().FindStringSubmatch(line); m != nil {
			doc.cards = append(doc.cards, roadmapCardRef{m[1], n})
			continue
		}
		// ponytail: wrapped titles are read for at most 8 lines, raise if a title ever wraps further.
		// Not the strict form: read the whole title (up to </h3>, at most 8 lines), strip tags, and flag it if it still starts like a version.
		title := line[roadmapH3Open().FindStringIndex(line)[0]:]
		for j := i; !strings.Contains(lines[j], "</h3>") && j+1 < len(lines) && j < i+8; j++ {
			title += " " + lines[j+1]
		}
		if k := strings.Index(title, ">"); k >= 0 {
			title = title[k+1:]
		}
		if k := strings.Index(title, "</h3>"); k >= 0 {
			title = title[:k]
		}
		title = roadmapTag().ReplaceAllString(title, "")
		probe := strings.NewReplacer("&nbsp;", " ", "\u00a0", " ").Replace(title)
		if roadmapVersionAt().MatchString(strings.TrimSpace(probe)) {
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

func checkRoadmaps(root string, errors *[]string) {
	const rel = "docs/roadmap.html"
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); err != nil {
		return
	}
	add := func(format string, args ...any) { *errors = append(*errors, fmt.Sprintf(format, args...)) }

	text, err := readText(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		add("%s:1: esperado ler o arquivo do roadmap, achei ilegível", rel)
		return
	}
	d := parseRoadmap(rel, text)

	changelog, err := readText(filepath.Join(root, "CHANGELOG.md"))
	if err != nil {
		add("CHANGELOG.md:1: esperado o CHANGELOG.md porque o roadmap existe, achei ausente ou ilegível")
		return
	}
	var versions []string
	dates := map[string]string{}
	for i, line := range strings.Split(changelog, "\n") {
		m := changelogHeading().FindStringSubmatch(line)
		switch {
		case m == nil && changelogAnyRel().MatchString(line):
			add("CHANGELOG.md:%d: esperado um título como \"## 0.1.0 — 2026-01-01\", achei %q", i+1, line)
		case m == nil:
		case dates[m[1]] != "":
			add("CHANGELOG.md:%d: esperado um título por release, achei %s de novo", i+1, m[1])
		default:
			versions = append(versions, m[1])
			dates[m[1]] = m[2]
		}
	}
	if len(versions) == 0 {
		add("CHANGELOG.md:1: esperado um título como \"## 0.1.0 — 2026-01-01\", achei nenhum")
		return
	}
	latest, latestDate := versions[0], dates[versions[0]]

	check := func(name string, v roadmapValue, want string) {
		switch {
		case v.line == 0:
			add("%s:1: esperado %s %q, achei nenhum", d.rel, name, want)
		case v.text != want:
			add("%s:%d: esperado %s %q, achei %q", d.rel, v.line, name, want, v.text)
		}
	}
	check("chip de data", d.date, latestDate)
	check("chip de status", d.status, latest+" · vivo")
	check("fim do intervalo de fontes", d.rangeEnd, latest)
	if d.path.line == 0 {
		add("%s:1: esperado o chip de caminho, achei nenhum", d.rel)
	}
	if d.entregueLine == 0 {
		add("%s:1: esperado <h2 id=\"entregue\">, achei nenhum", d.rel)
		return
	}
	for _, b := range d.badCards {
		add("%s:%d: esperado \"<h3>N.N.N · \", achei %q", d.rel, b.line, b.found)
	}
	count := map[string]int{}
	for _, c := range d.cards {
		count[c.version]++
		if count[c.version] == 2 {
			add("%s:%d: esperado exatamente um card para %s, achei um duplicado", d.rel, c.line, c.version)
		}
	}
	inChangelog := map[string]bool{}
	for _, v := range versions {
		inChangelog[v] = true
		if count[v] == 0 {
			add("%s:%d: esperado um card de Entregue para %s, achei nenhum", d.rel, d.entregueLine, v)
		}
	}
	for _, c := range d.cards {
		if !inChangelog[c.version] {
			add("%s:%d: esperado nenhum card para %s (sem título no CHANGELOG.md), achei um", d.rel, c.line, c.version)
		}
	}
	for j := 1; j < len(d.cards); j++ {
		prev, cur := d.cards[j-1], d.cards[j]
		if semverLess(prev.version, cur.version) {
			add("%s:%d: esperado o mais novo primeiro, achei o card %s depois do card %s (linha %d)", d.rel, cur.line, cur.version, prev.version, prev.line)
		}
	}
}
