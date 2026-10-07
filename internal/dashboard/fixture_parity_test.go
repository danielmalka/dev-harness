package dashboard

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

// The plugin test cannot read testdata/status-cases.json, so the mod keeps a .ts copy.
// This fails when the copy and the Go fixture stop describing the same cases.
func TestStatusCasesTSCopyMatches(t *testing.T) {
	js, err := os.ReadFile("testdata/status-cases.json")
	if err != nil {
		t.Fatal(err)
	}
	ts, err := os.ReadFile("../../adapters/claude-code/plugin/hooks/status-cases.fixture.ts")
	if err != nil {
		t.Fatal(err)
	}
	i := strings.Index(string(ts), "= [")
	if i < 0 {
		t.Fatal("no array literal in status-cases.fixture.ts")
	}
	var want, got []map[string]string
	if err := json.Unmarshal(js, &want); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(string(ts)[i+2:])), &got); err != nil {
		t.Fatalf("status-cases.fixture.ts array is not plain JSON: %v", err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatal("status-cases.fixture.ts differs from internal/dashboard/testdata/status-cases.json; regenerate the copy")
	}
}
