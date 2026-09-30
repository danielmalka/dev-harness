package kit

import (
	"bytes"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

// R1: the real tree passes the roadmap check.
func TestRealTreeRoadmapsClean(t *testing.T) {
	report, err := Validate(repoRoot(t), Options{SourceOnly: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range report.Errors {
		if strings.Contains(e, "roadmap.html") || strings.HasPrefix(e, "CHANGELOG.md:") {
			t.Errorf("roadmap error on the real tree: %s", e)
		}
	}
}

func copySourceTree(t *testing.T, src, dst string) {
	t.Helper()
	err := filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, p)
		if d.IsDir() && (rel == "dist" || rel == ".git") {
			return filepath.SkipDir
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if !d.Type().IsRegular() {
			return nil
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
}

// R2, R3, R4, R8: a copy of the source tree passes unmutated and fails, with the
// expected file:line message and a non-zero exit, once per single-field mutation.
func TestCopiedTreeMutations(t *testing.T) {
	if testing.Short() {
		t.Skip("copies the source tree")
	}
	src := repoRoot(t)
	bin := filepath.Join(t.TempDir(), "dh")
	build := exec.Command("go", "build", "-o", bin, "./cmd/dh")
	build.Dir = src
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build dh: %v\n%s", err, out)
	}
	copyDir := t.TempDir()
	copySourceTree(t, src, copyDir)

	run := func() (int, string) {
		cmd := exec.Command(bin, "validate", "--source-only", copyDir)
		var buf bytes.Buffer
		cmd.Stdout, cmd.Stderr = &buf, &buf
		err := cmd.Run()
		code := 0
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else if err != nil {
			t.Fatal(err)
		}
		return code, buf.String()
	}

	if code, out := run(); code != 0 {
		t.Fatalf("unmutated copy exit %d:\n%s", code, out)
	}

	pt := filepath.Join(copyDir, "docs", "roadmap.html")
	orig, err := os.ReadFile(pt)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(string(orig), "\n")
	find := func(sub string) int {
		for i, l := range lines {
			if strings.Contains(l, sub) {
				return i
			}
		}
		t.Fatalf("no line with %q", sub)
		return -1
	}
	muts := []struct{ name, marker, from, to string }{
		{"date chip", `<span class="chip">data <b>`, "", ""},
		{"status chip", `<span class="chip">status <b>`, "", ""},
		{"fontes line", `<code>CHANGELOG.md</code> — releases 0.1.0 a `, "", ""},
	}
	for _, m := range muts {
		t.Run(m.name, func(t *testing.T) {
			i := find(m.marker)
			mut := append([]string(nil), lines...)
			switch m.name {
			case "date chip":
				mut[i] = strings.Replace(mut[i], "<b>", "<b>1999-01-01 x", 1)
			case "status chip":
				mut[i] = strings.Replace(mut[i], "<b>", "<b>9.9.9 · ", 1)
			default:
				mut[i] = strings.Replace(mut[i], " a ", " a 9.9.9 x ", 1)
			}
			if mut[i] == lines[i] {
				t.Fatal("mutation did not change the line")
			}
			if err := os.WriteFile(pt, []byte(strings.Join(mut, "\n")), 0o644); err != nil {
				t.Fatal(err)
			}
			defer os.WriteFile(pt, orig, 0o644)
			code, out := run()
			want := "docs/roadmap.html:" + strconv.Itoa(i+1) + ":"
			if code == 0 || !strings.Contains(out, want) {
				t.Fatalf("exit %d, want non-zero with %q in:\n%s", code, want, out)
			}
		})
	}
}
