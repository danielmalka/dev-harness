package kit

import (
	"path/filepath"
	"strings"
	"testing"
)

// fixedDocCountInventory stands in for real counts of 19 agents, 17 commands, 21 skills, 4 profiles and 8 templates (this repository's shape after ADR-011, with
// `templates/` in a single folder), without depending on the repository's real content.
func fixedDocCountInventory() Inventory {
	return Inventory{
		Agents:    make([]string, 19),
		Commands:  make([]string, 17),
		Skills:    make([]string, 21),
		Profiles:  make([]string, 4),
		Templates: make([]string, 8),
	}
}

func TestCheckDocumentedCounts(t *testing.T) {
	tests := []struct {
		name    string
		path    string // relative to the temp root
		content string
		want    []string // substrings every wanted error must contain
		noError bool
	}{
		{
			name:    "correct current-state enumeration, pt-br",
			path:    "docs/roadmap.html",
			content: "<p>Contagem atual: 19 agentes, 17 comandos, 21 skills, 4 perfis, 8 templates</p>",
			noError: true,
		},
		{
			name:    "o espelho em inglês (docs/en/) é ignorado, mesmo com contagem velha em pt-BR",
			path:    "docs/en/roadmap.html",
			content: "<p>Esperado: 18 agentes, 17 comandos, 19 skills, 3 perfis</p>",
			noError: true,
		},
		{
			name:    "stale expected-state enumeration, pt-br",
			path:    "docs/tutoriais/00-primeira-maquina.html",
			content: "<p>Esperado: 18 agentes, 17 comandos, 19 skills, 3 perfis</p>",
			want: []string{
				"a contagem documentada de agents é 18, a real é 19",
				"a contagem documentada de skills é 19, a real é 21",
				"a contagem documentada de profiles é 3, a real é 4",
			},
		},
		{
			name:    "o espelho em inglês, numa subpasta de docs/en/, também é ignorado",
			path:    "docs/en/tutorials/00-first-machine.html",
			content: "<p>Esperado: 18 agentes, 17 comandos, 19 skills, 3 perfis</p>",
			noError: true,
		},
		{
			name:    "per-language aside is a single-noun line, not compared even with an anchor word",
			path:    "docs/roadmap.html",
			content: "<p>Contagem atual: 8 templates espelhados por idioma.</p>",
			noError: true,
		},
		{
			name:    "contagem em inglês não é mais reconhecida (o kit é só pt-BR)",
			path:    "docs/roadmap.html",
			content: "<p>Current count: 18 agents, 17 commands, 19 skills, 3 profiles</p>",
			noError: true,
		},
		{
			name:    "changelog entry describing a past release has no anchor word",
			path:    "docs/roadmap.html",
			content: "<p><strong>Resumo:</strong> 18 agentes, 17 comandos e 19 skills como fontes canonicas</p>",
			noError: true,
		},
		{
			name:    "dated product-plan narrative has no anchor word",
			path:    "docs/plano-produto.html",
			content: "<p>Inspecao desta revisao: fontes canonicas com 18 agentes, 17 comandos e 19 skills</p>",
			noError: true,
		},
		{
			name:    "single mention with an anchor word is still only one noun",
			path:    "docs/roadmap.html",
			content: "<p>Atual: 5 agentes em algum lugar sem relação</p>",
			noError: true,
		},
		{
			name:    "README at root is scanned",
			path:    "README.md",
			content: "Esperado: 18 agentes, 17 comandos, 19 skills, 3 perfis",
			want:    []string{"README.md:1: a contagem documentada de agents é 18, a real é 19"},
		},
		{
			name:    "non-doc extension under docs/ is ignored",
			path:    "docs/data/counts.json",
			content: `{"note": "Esperado: 18 agentes, 17 comandos, 19 skills, 3 perfis"}`,
			noError: true,
		},
		{
			name:    "correct enumeration produces no error even with an anchor word",
			path:    "README.en.md",
			content: "Esperado: 18 agentes, 17 comandos, 19 skills, 3 perfis",
			noError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			writeFixture(t, filepath.Join(root, filepath.FromSlash(test.path)), test.content)

			var errs []string
			checkDocumentedCounts(root, fixedDocCountInventory(), &errs)
			joined := strings.Join(errs, "\n")

			if test.noError {
				if len(errs) != 0 {
					t.Fatalf("errors = %v, want none", errs)
				}
				return
			}
			for _, want := range test.want {
				if !strings.Contains(joined, want) {
					t.Fatalf("errors = %v, want substring %q", errs, want)
				}
			}
			// commands (17 == 17 in the stale cases) must not be reported.
			if strings.Contains(joined, "a contagem documentada de commands") {
				t.Fatalf("errors = %v, want the matching commands count left unreported", errs)
			}
		})
	}
}

// TestCheckDocumentedCountsMultipleFilesAndLines confirms the file path and
// 1-based line number in a report point at the right place, across more than
// one file.
func TestCheckDocumentedCountsMultipleFilesAndLines(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, filepath.Join(root, "docs", "a.html"), "linha um\nlinha dois\nContagem atual: 18 agentes, 17 comandos\n")
	writeFixture(t, filepath.Join(root, "docs", "sub", "b.md"), "Esperado: 3 perfis, 8 templates\n")

	var errs []string
	checkDocumentedCounts(root, fixedDocCountInventory(), &errs)
	joined := strings.Join(errs, "\n")

	if !strings.Contains(joined, "docs/a.html:3: a contagem documentada de agents é 18, a real é 19") {
		t.Fatalf("errors = %v, want a.html:3 agents mismatch", errs)
	}
	if !strings.Contains(joined, "docs/sub/b.md:1: a contagem documentada de profiles é 3, a real é 4") {
		t.Fatalf("errors = %v, want b.md:1 profiles mismatch", errs)
	}
}
