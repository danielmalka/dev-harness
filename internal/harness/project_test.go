package harness

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeProject(t *testing.T, body string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "project.yaml")
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

// conventions.code_language (ADR-011): `en` and `pt-br` pass, any other value fails, and an absent key is valid (the default is English).
func TestCodeLanguageAcceptsEnAndPtBrAndRejectsTheRest(t *testing.T) {
	for _, v := range []string{"en", "pt-br", `"en"`, `'pt-br'`, "en # inglês", `"pt-br" # português`} {
		got, err := ReadProjectCheck(writeProject(t, "profile: none\nconventions:\n  code_language: "+v+"\nnotes: []\n"))
		if err != nil || len(got.Problems) != 0 || (got.CodeLanguage != "en" && got.CodeLanguage != "pt-br") {
			t.Errorf("%s deveria passar: %+v %v", v, got, err)
		}
	}
	for _, v := range []string{"xx", "pt-BR", "EN", "pt_br", "portuguese", `""`, ""} {
		got, err := ReadProjectCheck(writeProject(t, "conventions:\n  code_language: "+v+"\n"))
		if err != nil || len(got.Problems) != 1 || !strings.Contains(got.Problems[0], "conventions.code_language") || !strings.Contains(got.Problems[0], "(en, pt-br)") {
			t.Errorf("%q deveria reprovar com 1 problema: %+v %v", v, got, err)
		}
	}
}

func TestCodeLanguageOnlyCountsUnderTopLevelConventions(t *testing.T) {
	for name, body := range map[string]string{
		"sem a chave":               "profile: none\nconventions:\n  other: x\n",
		"outro bloco":               "notes:\n  code_language: xx\n",
		"conventions_applied":       "conventions_applied:\n  code_language: xx\n",
		"comentário":                "conventions:\n  # code_language: xx\n  other: 1\n",
		"conventions como lista":    "conventions:\n  - code_language: xx\n",
		"chave de topo igual":       "code_language: xx\n",
		"depois de outro bloco":     "conventions:\n  other: 1\nnotes:\n  code_language: xx\n",
		"arquivo vazio":             "",
		"só comentário e em branco": "# nada\n\n",
	} {
		got, err := ReadProjectCheck(writeProject(t, body))
		if err != nil || len(got.Problems) != 0 || got.CodeLanguage != "" {
			t.Errorf("%s: %+v %v", name, got, err)
		}
	}
}

func TestProjectLegacyLanguageIsDetectedAndNeverAProblem(t *testing.T) {
	for _, line := range []string{"language: en", "language: pt-br", "language: fr", "language:"} {
		got, err := ReadProjectCheck(writeProject(t, "# exemplo\n"+line+"\nprofile: none\n"))
		if err != nil || !got.LegacyLanguage || len(got.Problems) != 0 {
			t.Errorf("%q: %+v %v", line, got, err)
		}
	}
	for _, body := range []string{"# language: en\nprofile: none\n", "notes:\n  language: en\n", "languages: x\n", "profile: none\n"} {
		if got, _ := ReadProjectCheck(writeProject(t, body)); got.LegacyLanguage {
			t.Errorf("%q não é a opção antiga", body)
		}
	}
}

func TestReadProjectCheckAbsentAndInvalidFiles(t *testing.T) {
	if got, err := ReadProjectCheck(filepath.Join(t.TempDir(), "nao-existe.yaml")); err != nil || got.LegacyLanguage || len(got.Problems) != 0 {
		t.Errorf("ausente: %+v %v", got, err)
	}
	if _, err := ReadProjectCheck(t.TempDir()); err == nil {
		t.Error("um diretório no lugar do arquivo deveria dar erro")
	}
	big := writeProject(t, strings.Repeat("a: 1\n", maxConfigBytes/5+10))
	if _, err := ReadProjectCheck(big); err == nil {
		t.Error("arquivo acima do limite deveria dar erro")
	}
	// CRLF (arquivo vindo do Windows)
	if got, err := ReadProjectCheck(writeProject(t, "language: en\r\nconventions:\r\n  code_language: pt-br\r\n")); err != nil || !got.LegacyLanguage || got.CodeLanguage != "pt-br" || len(got.Problems) != 0 {
		t.Errorf("crlf: %+v %v", got, err)
	}
}
