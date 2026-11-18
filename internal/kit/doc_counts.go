package kit

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// docCountAnchor and docCountClaim implement checkDocumentedCounts.
//
// docCountClaim recognizes "<N> <noun>" for each of the five kit-content nouns the validator already counts, in pt-BR (the kit is pt-BR only since ADR-011).
// docCountAnchor recognizes the current/expected-state words this repository's prose actually uses around a real count claim ("Contagem atual: ...", "Esperado: ...").
var (
	docCountAnchor = lazyRe(`(?i)\b(atual|esperado)\b`)
	docCountClaim  = lazyRe(`\b(\d+)\s+(agentes|comandos|skills|perfis|templates)\b`)
)

// canonicalCountNoun maps the pt-BR spelling of a counted noun to the canonical name the inventory uses.
func canonicalCountNoun(word string) string {
	switch word {
	case "agentes":
		return "agents"
	case "comandos":
		return "commands"
	case "perfis":
		return "profiles"
	default:
		return word // skills, templates: same name in the inventory
	}
}

// realDocumentedCount returns the real count the validator computed for a canonical noun, or -1 if the noun is not one it tracks.
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

// checkDocumentedCounts flags prose across docs/**/*.html, docs/**/*.md and README.md that states a kit content count no longer matching the real count Validate just
// computed from disk (agents, commands, skills, profiles, templates). Recurring failure this repository keeps hitting: such prose silently falls behind reality, and nothing
// used to catch it. Only pt-BR files are read (ADR-011); the English mirror, `docs/en/` and `README.en.md`, was removed in 0.45.0 and stays ignored in case it comes back.
//
// A line is compared only when it both (a) names at least two of the five nouns and (b) carries one of the current/expected-state anchor words above, in the same physical
// line. Both conditions are load-bearing for avoiding false positives, which are worse than a miss here:
//
//   - Without the two-noun requirement, a standalone aside such as "8 templates por idioma" would be compared against the total template count and wrongly flagged:
//     it is a single-noun line describing a different quantity, not a total count claim.
//   - Without the anchor-word requirement, a changelog entry describing a past release ("Resumo: 18 agentes, 17 comandos e 19 skills" under a 0.1.0 card) would be flagged
//     even though the numbers were correct as of that point in time. Neither carries "atual" or "esperado".
//
// This narrows the check to exactly the shape this repository's real current-state claims use. A count claim that uses neither anchor word, or that appears alone on its
// line, is deliberately left unchecked.
func checkDocumentedCounts(root string, inv Inventory, errors *[]string) {
	for _, path := range docCountFiles(root) {
		text, err := readText(path)
		if err != nil {
			continue
		}
		for lineNumber, line := range strings.Split(text, "\n") {
			if !docCountAnchor().MatchString(line) {
				continue
			}
			matches := docCountClaim().FindAllStringSubmatch(line, -1)
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
						"%s:%d: a contagem documentada de %s é %d, a real é %d",
						filepath.ToSlash(rel), lineNumber+1, canonical, claimed, real,
					))
				}
			}
		}
	}
}

// docCountFiles lists every docs/**/*.html, docs/**/*.md and README.md under root, in sorted order, without the English mirror (`docs/en/` and `README.en.md`, removed in 0.45.0).
func docCountFiles(root string) []string {
	var paths []string
	docsRoot := filepath.Join(root, "docs")
	if isDir(docsRoot) {
		_ = filepath.WalkDir(docsRoot, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if entry.IsDir() {
				if path == filepath.Join(docsRoot, "en") {
					return filepath.SkipDir
				}
				return nil
			}
			switch strings.ToLower(filepath.Ext(path)) {
			case ".html", ".md":
				paths = append(paths, path)
			}
			return nil
		})
	}
	if readme := filepath.Join(root, "README.md"); isFile(readme) {
		paths = append(paths, readme)
	}
	sort.Strings(paths)
	return paths
}
