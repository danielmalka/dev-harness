package doctor

import (
	"bytes"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/danielmalka/dev-harness/internal/kit"
)

func TestRunReportsRequiredSectionsAndMatchesValidator(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	for _, target := range []string{root, filepath.Join(root, "dist", "claude-code", "dev-harness")} {
		t.Run(target, func(t *testing.T) {
			var output bytes.Buffer
			status := Run(target, &output)
			report, err := kit.Validate(target, kit.Options{})
			if err != nil {
				t.Fatal(err)
			}
			wantStatus := 0
			if !report.Passed() {
				wantStatus = 1
			}
			if status != wantStatus {
				t.Fatalf("status = %d, want %d; output:\n%s", status, wantStatus, output.String())
			}
			for _, section := range []string{"## Mode", "## Kit", "## Environment", "## Harness records in current directory", "## Limits"} {
				if !strings.Contains(output.String(), section) {
					t.Errorf("output missing %q", section)
				}
			}
		})
	}
}
