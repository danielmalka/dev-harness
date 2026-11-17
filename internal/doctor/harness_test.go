package doctor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func mk(t *testing.T, p, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// report runs the harness section for repo with an isolated home and returns the home too.
func report(t *testing.T, repo, home, settings string) string {
	t.Helper()
	t.Setenv("DH_HOME", home)
	var b bytes.Buffer
	printHarnessFor(&b, repo, home, settings)
	return b.String()
}

func TestHarnessRepoRegistered(t *testing.T) {
	repo, home := t.TempDir(), t.TempDir()
	mk(t, filepath.Join(repo, ".harness", "project.yaml"), "x\n")
	mk(t, filepath.Join(home, "config.yaml"), "projects:\n  \""+repo+"\": p\n")
	out := report(t, repo, home, filepath.Join(t.TempDir(), "none.json"))
	for _, want := range []string{"home: " + home, "mode: repo", "dir: " + repo + "/.harness", "project.yaml: presente", "MEMORY.md: ausente"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	for _, bad := range []string{"não está registrado", "existem os dois", "não está ignorado pelo git"} {
		if strings.Contains(out, bad) {
			t.Errorf("unexpected %q in:\n%s", bad, out)
		}
	}
}

func TestHarnessUnregisteredAndBoth(t *testing.T) {
	repo, home := t.TempDir(), t.TempDir()
	mk(t, filepath.Join(repo, ".harness", "project.yaml"), "x\n")
	out := report(t, repo, home, "")
	if !strings.Contains(out, "aviso: "+repo+" não está registrado em "+home+"/config.yaml: não aparece no dashboard até o dh link") {
		t.Errorf("no unregistered warning:\n%s", out)
	}
	mk(t, filepath.Join(home, "config.yaml"), "projects:\n  \""+repo+"\": p\n")
	os.MkdirAll(filepath.Join(home, "projects", "p"), 0o755)
	out = report(t, repo, home, "")
	if !strings.Contains(out, "mode: repo") || !strings.Contains(out, "aviso: "+repo+"/.harness e "+home+"/projects/p existem os dois: vale o interno") {
		t.Errorf("no both warning:\n%s", out)
	}
}

func TestHarnessGlobalAndNone(t *testing.T) {
	repo, home := t.TempDir(), t.TempDir()
	mk(t, filepath.Join(home, "config.yaml"), "projects:\n  \""+repo+"\": p\n")
	mk(t, filepath.Join(home, "projects", "p", "project.yaml"), "x\n")
	out := report(t, repo, home, "")
	if !strings.Contains(out, "mode: global") || !strings.Contains(out, "dir: "+home+"/projects/p") || strings.Contains(out, "aviso") {
		t.Errorf("global:\n%s", out)
	}
	moved := t.TempDir() // R6: repo moved, old path stays in the map
	out = report(t, moved, home, "")
	if !strings.Contains(out, "mode: none") || strings.Contains(out, "dir:") || !strings.Contains(out, "dica: não há harness para este diretório: rode /dh:setup (ou dh link)") {
		t.Errorf("none:\n%s", out)
	}
}

func TestHarnessLocalYamlIgnore(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	repo, home := t.TempDir(), t.TempDir()
	if err := exec.Command("git", "-C", repo, "init", "-q").Run(); err != nil {
		t.Skip("git init failed")
	}
	mk(t, filepath.Join(repo, ".harness", "local.yaml"), "x\n")
	mk(t, filepath.Join(home, "config.yaml"), "projects:\n  \""+repo+"\": p\n")
	const warn = "aviso: .harness/local.yaml não está ignorado pelo git"
	if out := report(t, repo, home, ""); !strings.Contains(out, warn) {
		t.Errorf("want warning:\n%s", out)
	}
	mk(t, filepath.Join(repo, ".gitignore"), ".harness/local.yaml\n")
	if out := report(t, repo, home, ""); strings.Contains(out, warn) {
		t.Errorf("unexpected warning:\n%s", out)
	}
	mk(t, filepath.Join(repo, ".gitignore"), ".harness/\n") // whole folder ignored: no warning either
	if out := report(t, repo, home, ""); strings.Contains(out, warn) || strings.Contains(out, "ignorado pelo git") {
		t.Errorf("unexpected warning:\n%s", out)
	}
}

func TestHarnessSettingsHint(t *testing.T) {
	repo, home, cfgDir := t.TempDir(), t.TempDir(), t.TempDir()
	mk(t, filepath.Join(repo, ".harness", "project.yaml"), "x\n")
	const hint = "dica: ~/.claude/settings.json não tem entrada em permissions.additionalDirectories que cubra " // + home
	settings := filepath.Join(cfgDir, "settings.json")
	if out := report(t, repo, home, settings); !strings.Contains(out, hint+home+`; acrescente "`+home+`" (o dh nunca escreve esse arquivo)`) {
		t.Errorf("missing file: want hint:\n%s", out)
	}
	mk(t, settings, `{"permissions":{"additionalDirectories":["/elsewhere"]}}`)
	if out := report(t, repo, home, settings); !strings.Contains(out, hint) {
		t.Errorf("no entry: want hint:\n%s", out)
	}
	mk(t, settings, `{not json`)
	if out := report(t, repo, home, settings); !strings.Contains(out, "não pôde ser lido") {
		t.Errorf("bad json: want parse hint:\n%s", out)
	}
	mk(t, settings, `{"permissions":{"additionalDirectories":["`+home+`"]}}`)
	if out := report(t, repo, home, settings); strings.Contains(out, "dica:") {
		t.Errorf("covered: unexpected hint:\n%s", out)
	}
	// home inside the repository: never a hint
	inside := filepath.Join(repo, "dh-home")
	os.MkdirAll(inside, 0o755)
	if out := report(t, repo, inside, filepath.Join(cfgDir, "absent.json")); strings.Contains(out, "settings.json") {
		t.Errorf("home inside repo: unexpected hint:\n%s", out)
	}
	// "~/..." entries expand against the user home
	uh := t.TempDir()
	t.Setenv("HOME", uh)
	mk(t, settings, `{"permissions":{"additionalDirectories":["~/.harness"]}}`)
	if out := report(t, repo, filepath.Join(uh, ".harness"), settings); strings.Contains(out, "dica:") {
		t.Errorf("tilde entry: unexpected hint:\n%s", out)
	}
}

func TestHarnessUnknownHome(t *testing.T) {
	repo := t.TempDir()
	mk(t, filepath.Join(repo, ".harness", "project.yaml"), "x\n")
	t.Setenv("DH_HOME", "")
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "")
	var b bytes.Buffer
	printHarnessFor(&b, repo, "", "")
	if out := b.String(); !strings.Contains(out, "home: desconhecido") || !strings.Contains(out, "defina DH_HOME") {
		t.Errorf("unknown home:\n%s", out)
	}
}

func TestSettingsHintSuggestsCustomHome(t *testing.T) {
	uh, custom := t.TempDir(), t.TempDir()
	t.Setenv("HOME", uh)
	if h := settingsHint(filepath.Join(uh, "none.json"), custom); !strings.Contains(h, `acrescente "`+custom+`"`) {
		t.Errorf("custom home: %s", h)
	}
	if h := settingsHint(filepath.Join(uh, "none.json"), filepath.Join(uh, ".harness")); !strings.Contains(h, `acrescente "~/.harness"`) {
		t.Errorf("default home: %s", h)
	}
}

func TestHarnessPdocsOnlyNoBothWarning(t *testing.T) {
	repo, home := t.TempDir(), t.TempDir()
	mk(t, filepath.Join(repo, ".harness", "project.yaml"), "x\n")
	mk(t, filepath.Join(home, "config.yaml"), "projects:\n  \""+repo+"\": p\n")
	os.MkdirAll(filepath.Join(home, "projects", "p", "pdocs"), 0o755)
	out := report(t, repo, home, "")
	if strings.Contains(out, "existem os dois") {
		t.Errorf("unexpected both warning:\n%s", out)
	}
	mk(t, filepath.Join(home, "projects", "p", "project.yaml"), "x\n")
	if out = report(t, repo, home, ""); !strings.Contains(out, "existem os dois") {
		t.Errorf("no both warning:\n%s", out)
	}
}

func TestHarnessSymlinkedCwd(t *testing.T) {
	real, home := t.TempDir(), t.TempDir()
	real, _ = filepath.EvalSymlinks(real)
	link := filepath.Join(t.TempDir(), "ln")
	if err := os.Symlink(real, link); err != nil {
		t.Skip("symlinks unavailable")
	}
	mk(t, filepath.Join(real, ".harness", "project.yaml"), "x\n")
	mk(t, filepath.Join(home, "config.yaml"), "projects:\n  \""+real+"\": p\n")
	os.MkdirAll(filepath.Join(home, "projects", "p"), 0o755)
	out := report(t, link, home, "")
	if strings.Contains(out, "não está registrado") || !strings.Contains(out, "existem os dois") {
		t.Errorf("symlinked cwd:\n%s", out)
	}
}

func TestHarnessMemorySize(t *testing.T) {
	for _, tc := range []struct {
		n    int
		warn bool
	}{{39999, false}, {40000, false}, {40001, true}} {
		repo, home := t.TempDir(), t.TempDir()
		mk(t, filepath.Join(repo, ".harness", "project.yaml"), "x\n")
		mk(t, filepath.Join(repo, ".harness", "MEMORY.md"), strings.Repeat("a", tc.n))
		out := report(t, repo, home, "")
		size := fmt.Sprintf("tamanho do MEMORY.md: %d bytes", tc.n)
		warn := fmt.Sprintf("aviso: o MEMORY.md tem %d bytes, acima do orçamento de 40000: rode /dh:consolidate-memory", tc.n)
		if !strings.Contains(out, size) || !strings.Contains(out, "MEMORY.md: presente") {
			t.Errorf("%d: missing size line in:\n%s", tc.n, out)
		}
		if strings.Contains(out, warn) != tc.warn {
			t.Errorf("%d: warning=%v want %v in:\n%s", tc.n, !tc.warn, tc.warn, out)
		}
	}
	repo, home := t.TempDir(), t.TempDir()
	mk(t, filepath.Join(repo, ".harness", "project.yaml"), "x\n")
	if out := report(t, repo, home, ""); strings.Contains(out, "tamanho do MEMORY.md") {
		t.Errorf("size line without MEMORY.md:\n%s", out)
	}
}

// ADR-011: the `language` option no longer exists. A legacy line in project.yaml or config.yaml is read, ignored, NEVER rewritten, and the doctor only warns (no failure).
func TestLegacyLanguageWarnsWithoutFailingAndNeverRewritesTheFiles(t *testing.T) {
	repo, home := t.TempDir(), t.TempDir()
	project := "language: en\nprofile: none\n"
	cfg := "language: en\nprojects:\n  \"" + repo + "\": p\n"
	mk(t, filepath.Join(repo, ".harness", "project.yaml"), project)
	mk(t, filepath.Join(home, "config.yaml"), cfg)
	var b bytes.Buffer
	t.Setenv("DH_HOME", home)
	if problems := printHarnessFor(&b, repo, home, ""); problems != 0 {
		t.Errorf("language antiga não pode reprovar: %d problemas\n%s", problems, b.String())
	}
	out := b.String()
	for _, file := range []string{filepath.Join(repo, ".harness", "project.yaml"), filepath.Join(home, "config.yaml")} {
		want := "aviso: `language` não existe mais e é ignorada (o kit é só pt-BR); pode apagar a linha de " + file
		if !strings.Contains(out, want) {
			t.Errorf("falta o aviso de %s:\n%s", file, out)
		}
	}
	for path, want := range map[string]string{filepath.Join(repo, ".harness", "project.yaml"): project, filepath.Join(home, "config.yaml"): cfg} {
		if got, _ := os.ReadFile(path); string(got) != want {
			t.Errorf("%s foi reescrito: %q", path, got)
		}
	}
}

func TestNoLegacyLanguageNoWarning(t *testing.T) {
	repo, home := t.TempDir(), t.TempDir()
	mk(t, filepath.Join(repo, ".harness", "project.yaml"), "profile: none\nconventions:\n  code_language: pt-br\n")
	mk(t, filepath.Join(home, "config.yaml"), "mode: repo\nprojects:\n  \""+repo+"\": p\n")
	out := report(t, repo, home, "")
	if strings.Contains(out, "language") || strings.Contains(out, "erro:") {
		t.Errorf("sem linha antiga e com code_language válido não há aviso nem erro:\n%s", out)
	}
}

// conventions.code_language: `en` and `pt-br` pass; any other value fails (the doctor answers with an error and exit status 1).
func TestCodeLanguageInvalidFailsTheDoctorAndValidPasses(t *testing.T) {
	for _, tc := range []struct {
		value string
		fails bool
	}{{"en", false}, {"pt-br", false}, {"xx", true}, {"pt-BR", true}} {
		repo, home := t.TempDir(), t.TempDir()
		mk(t, filepath.Join(repo, ".harness", "project.yaml"), "conventions:\n  code_language: "+tc.value+"\n")
		mk(t, filepath.Join(home, "config.yaml"), "projects:\n  \""+repo+"\": p\n")
		t.Setenv("DH_HOME", home)
		var b bytes.Buffer
		problems := printHarnessFor(&b, repo, home, "")
		if (problems > 0) != tc.fails {
			t.Errorf("%s: problemas=%d, reprova=%v\n%s", tc.value, problems, tc.fails, b.String())
		}
		wantLine := "erro: conventions.code_language " + `"` + tc.value + `"` + " inválido (en, pt-br) em " + filepath.Join(repo, ".harness", "project.yaml")
		if has := strings.Contains(b.String(), wantLine); has != tc.fails {
			t.Errorf("%s: linha de erro presente=%v, esperado %v\n%s", tc.value, has, tc.fails, b.String())
		}
	}
}

func TestDoctorRunExitsOneOnInvalidCodeLanguageOnly(t *testing.T) {
	repo, home := t.TempDir(), t.TempDir()
	mk(t, filepath.Join(repo, ".harness", "project.yaml"), "conventions:\n  code_language: xx\n")
	mk(t, filepath.Join(home, "config.yaml"), "projects:\n  \""+repo+"\": p\n")
	t.Setenv("DH_HOME", home)
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	if err := os.Chdir(repo); err != nil {
		t.Fatal(err)
	}
	var b bytes.Buffer
	if status := Run(t.TempDir(), &b); status != 1 || !strings.Contains(b.String(), "status: reprovado (erros do harness acima)") {
		t.Errorf("status %d\n%s", status, b.String())
	}
	mk(t, filepath.Join(repo, ".harness", "project.yaml"), "language: en\n")
	b.Reset()
	if status := Run(t.TempDir(), &b); status != 0 || !strings.Contains(b.String(), "status: aprovado") {
		t.Errorf("language antiga não falha: status %d\n%s", status, b.String())
	}
}
