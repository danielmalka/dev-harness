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
