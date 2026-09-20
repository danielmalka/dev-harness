package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func testSnapshotDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("DEV_HARNESS_SNAPSHOT_DIR", dir)
	return dir
}

func readTestSnapshot(t *testing.T, dir, sessionID string) Snapshot {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, sessionID+".json"))
	if err != nil {
		t.Fatal(err)
	}
	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		t.Fatal(err)
	}
	return snapshot
}

func TestEventLifecycleAndSessionStartTimestamp(t *testing.T) {
	dir := testSnapshotDir(t)
	start := `{"session_id":"s1","hook_event_name":"SessionStart","cwd":"/work","agent_type":"coordinator","model":"claude","session_name":"build"}`
	if err := Event(dir, strings.NewReader(start)); err != nil {
		t.Fatal(err)
	}
	first := readTestSnapshot(t, dir, "s1")
	if first.State != "idle" || first.StartedAt == "" || first.CWD != "/work" || first.Agent != "coordinator" || first.Model == nil || first.Model.ID != "claude" {
		t.Fatalf("unexpected SessionStart snapshot: %+v", first)
	}

	if err := Event(dir, strings.NewReader(start)); err != nil {
		t.Fatal(err)
	}
	second := readTestSnapshot(t, dir, "s1")
	if second.StartedAt != first.StartedAt {
		t.Fatalf("started_at changed from %q to %q", first.StartedAt, second.StartedAt)
	}

	for _, event := range []string{"UserPromptSubmit", "Stop", "SessionEnd"} {
		input := `{"session_id":"s1","hook_event_name":"` + event + `"}`
		if err := Event(dir, strings.NewReader(input)); err != nil {
			t.Fatal(err)
		}
	}
	final := readTestSnapshot(t, dir, "s1")
	if final.State != "closed" {
		t.Fatalf("state = %q, want closed", final.State)
	}

	if err := Event(dir, strings.NewReader(`{"session_id":"s1","hook_event_name":"Unknown"}`)); err != nil {
		t.Fatal(err)
	}
	unknown := readTestSnapshot(t, dir, "s1")
	if unknown.State != "closed" || unknown.StartedAt != final.StartedAt || unknown.SessionName != final.SessionName {
		t.Fatalf("unknown event changed snapshot fields: before=%+v after=%+v", final, unknown)
	}
}

func TestSubagentEventsKeepLastFifty(t *testing.T) {
	dir := testSnapshotDir(t)
	for i := 0; i < 60; i++ {
		input := `{"session_id":"s1","hook_event_name":"SubagentStart","agent_type":"qa","agent_id":"a` + string(rune('0'+i%10)) + `"}`
		if err := Event(dir, strings.NewReader(input)); err != nil {
			t.Fatal(err)
		}
	}
	snapshot := readTestSnapshot(t, dir, "s1")
	if len(snapshot.Events) != 50 {
		t.Fatalf("events = %d, want 50", len(snapshot.Events))
	}
	if snapshot.Events[0].At == "" || snapshot.Events[0].Event != "SubagentStart" || snapshot.State != "active" {
		t.Fatalf("unexpected first event or state: %+v", snapshot)
	}
}

func TestSubagentsReplacesTasksAndPrintsLines(t *testing.T) {
	dir := testSnapshotDir(t)
	input := `{"session_id":"s1","columns":120,"tasks":[{"id":"a","name":"qa","type":"ignored","status":"running","model":"sonnet","tokenCount":1234},{"id":"b","name":"build","status":"done"}]}`
	var output strings.Builder
	if err := Subagents(dir, strings.NewReader(input), &output); err != nil {
		t.Fatal(err)
	}
	want := "{\"id\":\"a\",\"content\":\"qa · sonnet · 1234 tok · running\"}\n" +
		"{\"id\":\"b\",\"content\":\"build · done\"}\n"
	if output.String() != want {
		t.Fatalf("output = %q, want %q", output.String(), want)
	}
	snapshot := readTestSnapshot(t, dir, "s1")
	if len(snapshot.Tasks) != 2 || snapshot.Tasks[0].Type != "ignored" || snapshot.Tasks[1].Name != "build" {
		t.Fatalf("unexpected tasks: %+v", snapshot.Tasks)
	}

	output.Reset()
	if err := Subagents(dir, strings.NewReader(`{"session_id":"s1","tasks":[]}`), &output); err != nil {
		t.Fatal(err)
	}
	if output.Len() != 0 || len(readTestSnapshot(t, dir, "s1").Tasks) != 0 {
		t.Fatalf("empty tasks were not handled: output=%q", output.String())
	}
}

func TestStatuslineMergesAndCountsActiveAgents(t *testing.T) {
	dir := testSnapshotDir(t)
	if err := Event(dir, strings.NewReader(`{"session_id":"s1","hook_event_name":"SessionStart"}`)); err != nil {
		t.Fatal(err)
	}
	if err := Subagents(dir, strings.NewReader(`{"session_id":"s1","tasks":[{"id":"a","status":"running"},{"id":"b","status":"done"},{"id":"c","status":"completed"},{"id":"d","status":"failed"},{"id":"e","status":"cancelled"}]}`), &strings.Builder{}); err != nil {
		t.Fatal(err)
	}
	var output strings.Builder
	input := `{"session_id":"s1","model":{"id":"m","display_name":"Sonnet"},"cost":{"total_cost_usd":0.1234,"total_duration_ms":42},"context_window":{"used_percentage":42.4,"context_window_size":1000,"total_input_tokens":10,"total_output_tokens":20},"workspace":{"current_dir":"/project"},"agent":{"name":"coordinator"},"rate_limits":{"remaining":3},"session_name":"run"}`
	if err := Statusline(dir, strings.NewReader(input), &output); err != nil {
		t.Fatal(err)
	}
	if output.String() != "Sonnet · ctx 42% · $0.12 · 1 agents active\n" {
		t.Fatalf("output = %q", output.String())
	}
	snapshot := readTestSnapshot(t, dir, "s1")
	if snapshot.State != "idle" || snapshot.CWD != "/project" || snapshot.Cost == nil || snapshot.Cost.TotalCostUSD != 0.1234 || snapshot.Context == nil || snapshot.Context.InputTokens != 10 || snapshot.Agent != "coordinator" {
		t.Fatalf("statusline changed unexpected fields: %+v", snapshot)
	}
}

func TestPruneDeletesOldAndInvalidFiles(t *testing.T) {
	dir := testSnapshotDir(t)
	old := Snapshot{Schema: 1, SessionID: "old", UpdatedAt: time.Now().Add(-8 * 24 * time.Hour).Format(time.RFC3339)}
	fresh := Snapshot{Schema: 1, SessionID: "fresh", UpdatedAt: time.Now().Format(time.RFC3339)}
	if err := writeSnapshot(dir, old); err != nil {
		t.Fatal(err)
	}
	if err := writeSnapshot(dir, fresh); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "invalid.json"), []byte("{"), 0600); err != nil {
		t.Fatal(err)
	}
	n, err := Prune(dir, 7)
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("pruned = %d, want 2", n)
	}
	if _, err := os.Stat(filepath.Join(dir, "fresh.json")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "old.json")); !os.IsNotExist(err) {
		t.Fatalf("old snapshot still exists: %v", err)
	}
}

func TestMissingSessionIDAndAtomicWrite(t *testing.T) {
	dir := testSnapshotDir(t)
	if err := Event(dir, strings.NewReader(`{"hook_event_name":"Stop"}`)); err != ErrMissingSessionID {
		t.Fatalf("error = %v, want %v", err, ErrMissingSessionID)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("missing session created files: %v", entries)
	}
	if err := Event(dir, strings.NewReader(`{"session_id":"s1","hook_event_name":"Stop"}`)); err != nil {
		t.Fatal(err)
	}
	tmp, err := filepath.Glob(filepath.Join(dir, "*.tmp"))
	if err != nil {
		t.Fatal(err)
	}
	if len(tmp) != 0 {
		t.Fatalf("temporary files remain: %v", tmp)
	}
}
