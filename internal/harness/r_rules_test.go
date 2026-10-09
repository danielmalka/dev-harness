package harness

import (
	"os"
	"strings"
	"testing"
)

// Rule-named entry points over the focused tests in harness_test.go.
func TestR2_ConfigRoundTripWindowsPathsAndNoNewDependency(t *testing.T) {
	TestConfigRoundTripPreservesBytes(t)
	TestSaveNewConfigQuotesWindowsPath(t)
	TestInvalidConfigIsReadableAndNotOverwritten(t)
	b, err := os.ReadFile("../../go.mod")
	if err != nil || strings.Contains(string(b), "require") {
		t.Errorf("go.mod must have no dependency (err=%v)", err)
	}
}

func TestR4_R6_LinkLifecycle(t *testing.T) {
	TestLinkCreateRelinkOff(t)
	TestLinkNameCollision(t)
	TestLinkOffOnStalePath(t)
}

func TestR5_ResolveTableAllSituations(t *testing.T) { TestResolveTable(t) }

func TestR17_ProjectsListCapAndOmissions(t *testing.T) { TestProjectsListSkipsAndCaps(t) }

// PRD-017 R2 (QA, T-1703-05): chained table plus each of the six records alone.
func TestR2_PdocsOnlyIsNotAHarness(t *testing.T) {
	TestIsHarnessAndResolvePdocsOnly(t)
	for _, it := range []string{"project.yaml", "MEMORY.md", "EPOCHAL.md", "RISKS.md", "tasks/", "prd/"} {
		d := t.TempDir()
		os.MkdirAll(d+"/pdocs", 0o755)
		if strings.HasSuffix(it, "/") {
			os.MkdirAll(d+"/"+it, 0o755)
		} else {
			os.WriteFile(d+"/"+it, nil, 0o644)
		}
		if !IsHarness(d) {
			t.Errorf("%s alone with pdocs/ should count as harness", it)
		}
	}
	TestProjectsListSkipsAndCaps(t)
}
