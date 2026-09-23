package kit

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// docCountAnchor and docCountClaim implement checkDocumentedCounts.
//
// docCountClaim recognizes "<N> <noun>" for each of the five kit-content
// nouns the validator already counts, in both languages ("skills" and
// "templates" are spelled the same in pt-br and en). docCountAnchor
// recognizes the current/expected-state words this repository's prose
// actually uses around a real count claim ("Contagem atual: ...",
// "Current count: ...", "Esperado: ...", "Expected: ...").
var (
	docCountAnchor = regexp.MustCompile(`(?i)\b(atual|esperado|current|expected)\b`)
	docCountClaim  = regexp.MustCompile(`\b(\d+)\s+(agentes|agents|comandos|commands|skills|perfis|profiles|templates)\b`)
)

// canonicalCountNoun folds the pt-br and en spellings of a counted noun onto
// one name, so "19 agentes" and "19 agents" are recognized as the same claim.
func canonicalCountNoun(word string) string {
	switch word {
	case "agentes", "agents":
		return "agents"
	case "comandos", "commands":
		return "commands"
	case "perfis", "profiles":
		return "profiles"
	default:
		return word // skills, templates: identical in both languages
	}
}

// realDocumentedCount returns the real count the validator computed for a
// canonical noun, or -1 if the noun is not one it tracks.
func realDocumentedCount(inv Inventory, canonical string) int {
	switch canonical {
	case "agents":
		return len(inv.Agents)
	case "commands":
		return len(inv.Commands)
	case "skills":
		return len(inv.Skills)
	case "profiles":
		return len(inv.Profiles)
	case "templates":
		return len(inv.Templates)
	}
	return -1
}

// checkDocumentedCounts flags prose across docs/**/*.html, docs/**/*.md and
// README*.md that states a kit content count no longer matching the real
// count Validate just computed from disk (agents, commands, skills,
// profiles, templates). Recurring failure this repository keeps hitting:
// such prose silently falls behind reality, and nothing used to catch it.
//
// A line is compared only when it both (a) names at least two of the five
// nouns and (b) carries one of the current/expected-state anchor words above,
// in the same physical line. Both conditions are load-bearing for avoiding
// false positives, which are worse than a miss here:
//
//   - Without the two-noun requirement, a per-language aside such as
//     "8 templates espelhados por idioma" ("8 templates in pt-br") would be
//     compared against the total template count (16) and wrongly flagged —
//     it is a single-noun line describing a different quantity, not a total
//     count claim.
//   - Without the anchor-word requirement, a changelog entry describing a
//     past release ("Resumo: 18 agentes, 17 comandos e 19 skills" under a
//     0.1.0 card) or the product plan's dated "Inspeção desta revisão: ...
//     18 agentes" narrative would be flagged even though the numbers were
//     correct as of that past point in time, not stale documentation of the
//     current state. Neither carries "atual/current" or "esperado/expected".
//
// This narrows the check to exactly the shape this repository's real
// current-state claims use (verified against docs/roadmap.html and
// docs/en/roadmap.html, whose "Contagem atual"/"Current count" line is
// correct today, and against the historical lines above, which must not
// fire). A count claim that uses neither anchor word, or that appears alone
// on its line, is deliberately left unchecked.
func checkDocumentedCounts(root string, inv Inventory, errors *[]string) {
	for _, path := range docCountFiles(root) {
		text, err := readText(path)
		if err != nil {
			continue
		}
		for lineNumber, line := range strings.Split(text, "\n") {
			if !docCountAnchor.MatchString(line) {
				continue
			}
			matches := docCountClaim.FindAllStringSubmatch(line, -1)
			nouns := map[string]bool{}
			for _, match := range matches {
				nouns[canonicalCountNoun(match[2])] = true
			}
			if len(nouns) < 2 {
				continue
			}
			rel, relErr := filepath.Rel(root, path)
			if relErr != nil {
				rel = path
			}
			for _, match := range matches {
				claimed, convErr := strconv.Atoi(match[1])
				if convErr != nil {
					continue
				}
				canonical := canonicalCountNoun(match[2])
				real := realDocumentedCount(inv, canonical)
				if real >= 0 && claimed != real {
					*errors = append(*errors, fmt.Sprintf(
						"%s:%d: documented %s count is %d, real count is %d",
						filepath.ToSlash(rel), lineNumber+1, canonical, claimed, real,
					))
				}
			}
		}
	}
}

// docCountFiles lists every docs/**/*.html, docs/**/*.md and README*.md
// under root, in sorted order.
func docCountFiles(root string) []string {
	var paths []string
	docsRoot := filepath.Join(root, "docs")
	if isDir(docsRoot) {
		_ = filepath.WalkDir(docsRoot, func(path string, entry os.DirEntry, err error) error {
			if err != nil || entry.IsDir() {
				return nil
			}
			switch strings.ToLower(filepath.Ext(path)) {
			case ".html", ".md":
				paths = append(paths, path)
			}
			return nil
		})
	}
	readmes, _ := filepath.Glob(filepath.Join(root, "README*.md"))
	paths = append(paths, readmes...)
	sort.Strings(paths)
	return paths
}
