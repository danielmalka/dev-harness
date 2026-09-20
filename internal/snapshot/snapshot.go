package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var ErrMissingSessionID = errors.New("missing session_id")

type Snapshot struct {
	Schema      int             `json:"schema"`
	SessionID   string          `json:"session_id"`
	SessionName string          `json:"session_name,omitempty"`
	CWD         string          `json:"cwd,omitempty"`
	Agent       string          `json:"agent,omitempty"`
	Model       *Model          `json:"model,omitempty"`
	State       string          `json:"state,omitempty"`
	StartedAt   string          `json:"started_at,omitempty"`
	UpdatedAt   string          `json:"updated_at,omitempty"`
	Cost        *Cost           `json:"cost,omitempty"`
	Context     *Context        `json:"context,omitempty"`
	RateLimits  json.RawMessage `json:"rate_limits,omitempty"`
	Tasks       []Task          `json:"tasks,omitempty"`
	Events      []EventRecord   `json:"events,omitempty"`
}

type Model struct {
	ID          string `json:"id,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

type Cost struct {
	TotalCostUSD    float64 `json:"total_cost_usd,omitempty"`
	TotalDurationMS int64   `json:"total_duration_ms,omitempty"`
}

type Context struct {
	UsedPercentage    float64 `json:"used_percentage,omitempty"`
	ContextWindowSize int64   `json:"context_window_size,omitempty"`
	InputTokens       int64   `json:"input_tokens,omitempty"`
	OutputTokens      int64   `json:"output_tokens,omitempty"`
}

type Task struct {
	ID                string `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	Type              string `json:"type,omitempty"`
	Status            string `json:"status,omitempty"`
	Model             string `json:"model,omitempty"`
	Effort            string `json:"effort,omitempty"`
	TokenCount        int64  `json:"tokenCount,omitempty"`
	ContextWindowSize int64  `json:"contextWindowSize,omitempty"`
	StartTime         string `json:"startTime,omitempty"`
}

type EventRecord struct {
	At        string `json:"at,omitempty"`
	Event     string `json:"event,omitempty"`
	AgentType string `json:"agent_type,omitempty"`
	AgentID   string `json:"agent_id,omitempty"`
}

type eventInput struct {
	SessionID     string  `json:"session_id"`
	HookEventName string  `json:"hook_event_name"`
	CWD           *string `json:"cwd"`
	AgentID       *string `json:"agent_id"`
	AgentType     *string `json:"agent_type"`
	Model         *string `json:"model"`
	SessionName   *string `json:"session_name"`
}

type taskInput struct {
	ID                *string `json:"id"`
	Name              *string `json:"name"`
	Type              *string `json:"type"`
	Status            *string `json:"status"`
	Description       *string `json:"description"`
	Label             *string `json:"label"`
	StartTime         *string `json:"startTime"`
	Model             *string `json:"model"`
	Effort            *string `json:"effort"`
	ContextWindowSize *int64  `json:"contextWindowSize"`
	TokenCount        *int64  `json:"tokenCount"`
}

type subagentsInput struct {
	SessionID string       `json:"session_id"`
	Tasks     *[]taskInput `json:"tasks"`
}

type statuslineInput struct {
	SessionID     string          `json:"session_id"`
	SessionName   *string         `json:"session_name"`
	Model         *modelInput     `json:"model"`
	Workspace     *workspaceInput `json:"workspace"`
	Cost          *costInput      `json:"cost"`
	ContextWindow *contextInput   `json:"context_window"`
	RateLimits    json.RawMessage `json:"rate_limits"`
	Agent         *agentInput     `json:"agent"`
}

type modelInput struct {
	ID          *string `json:"id"`
	DisplayName *string `json:"display_name"`
}

type workspaceInput struct {
	CurrentDir *string `json:"current_dir"`
}

type costInput struct {
	TotalCostUSD    *float64 `json:"total_cost_usd"`
	TotalDurationMS *int64   `json:"total_duration_ms"`
}

type contextInput struct {
	UsedPercentage    *float64 `json:"used_percentage"`
	ContextWindowSize *int64   `json:"context_window_size"`
	TotalInputTokens  *int64   `json:"total_input_tokens"`
	TotalOutputTokens *int64   `json:"total_output_tokens"`
}

type agentInput struct {
	Name *string `json:"name"`
}

type outputTask struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func Event(dir string, r io.Reader) error {
	var input eventInput
	if err := decodeJSON(r, &input); err != nil {
		return err
	}
	if input.SessionID == "" {
		return ErrMissingSessionID
	}

	snapshot, _, err := loadSnapshot(dir, input.SessionID)
	if err != nil {
		return err
	}
	if snapshot.Schema == 0 {
		snapshot.Schema = 1
	}
	now := timestamp()

	switch input.HookEventName {
	case "SessionStart":
		if input.SessionName != nil {
			snapshot.SessionName = *input.SessionName
		}
		if snapshot.StartedAt == "" {
			snapshot.StartedAt = now
		}
		snapshot.State = "idle"
		if input.CWD != nil {
			snapshot.CWD = *input.CWD
		}
		if input.AgentType != nil {
			snapshot.Agent = *input.AgentType
		}
		if input.Model != nil {
			if snapshot.Model == nil {
				snapshot.Model = &Model{}
			}
			snapshot.Model.ID = *input.Model
		}
	case "UserPromptSubmit":
		snapshot.State = "active"
	case "Stop":
		snapshot.State = "idle"
	case "SubagentStart", "SubagentStop":
		snapshot.State = "active"
		record := EventRecord{At: now, Event: input.HookEventName}
		if input.AgentType != nil {
			record.AgentType = *input.AgentType
		}
		if input.AgentID != nil {
			record.AgentID = *input.AgentID
		}
		snapshot.Events = append(snapshot.Events, record)
		if len(snapshot.Events) > 50 {
			snapshot.Events = snapshot.Events[len(snapshot.Events)-50:]
		}
	case "SessionEnd":
		snapshot.State = "closed"
	}
	snapshot.UpdatedAt = now
	return writeSnapshot(dir, *snapshot)
}

func Subagents(dir string, r io.Reader, w io.Writer) error {
	var input subagentsInput
	if err := decodeJSON(r, &input); err != nil {
		return err
	}
	if input.SessionID == "" {
		return ErrMissingSessionID
	}
	snapshot, _, err := loadSnapshot(dir, input.SessionID)
	if err != nil {
		return err
	}
	if input.Tasks == nil {
		snapshot.Tasks = nil
	} else {
		snapshot.Tasks = make([]Task, 0, len(*input.Tasks))
		for _, source := range *input.Tasks {
			task := taskFromInput(source)
			snapshot.Tasks = append(snapshot.Tasks, task)
			content := taskContent(source)
			line, err := json.Marshal(outputTask{ID: task.ID, Content: content})
			if err != nil {
				return err
			}
			if _, err := fmt.Fprintln(w, string(line)); err != nil {
				return err
			}
		}
	}
	if snapshot.Schema == 0 {
		snapshot.Schema = 1
	}
	snapshot.UpdatedAt = timestamp()
	return writeSnapshot(dir, *snapshot)
}

func Statusline(dir string, r io.Reader, w io.Writer) error {
	var input statuslineInput
	if err := decodeJSON(r, &input); err != nil {
		return err
	}
	if input.SessionID == "" {
		return ErrMissingSessionID
	}
	snapshot, existed, err := loadSnapshot(dir, input.SessionID)
	if err != nil {
		return err
	}
	if snapshot.Schema == 0 {
		snapshot.Schema = 1
	}
	if !existed {
		snapshot.State = "idle"
	}
	if input.SessionName != nil {
		snapshot.SessionName = *input.SessionName
	}
	if input.Model != nil {
		if snapshot.Model == nil {
			snapshot.Model = &Model{}
		}
		if input.Model.ID != nil {
			snapshot.Model.ID = *input.Model.ID
		}
		if input.Model.DisplayName != nil {
			snapshot.Model.DisplayName = *input.Model.DisplayName
		}
	}
	if input.Workspace != nil && input.Workspace.CurrentDir != nil {
		snapshot.CWD = *input.Workspace.CurrentDir
	}
	if input.Cost != nil {
		if snapshot.Cost == nil {
			snapshot.Cost = &Cost{}
		}
		if input.Cost.TotalCostUSD != nil {
			snapshot.Cost.TotalCostUSD = *input.Cost.TotalCostUSD
		}
		if input.Cost.TotalDurationMS != nil {
			snapshot.Cost.TotalDurationMS = *input.Cost.TotalDurationMS
		}
	}
	if input.ContextWindow != nil {
		if snapshot.Context == nil {
			snapshot.Context = &Context{}
		}
		if input.ContextWindow.UsedPercentage != nil {
			snapshot.Context.UsedPercentage = *input.ContextWindow.UsedPercentage
		}
		if input.ContextWindow.ContextWindowSize != nil {
			snapshot.Context.ContextWindowSize = *input.ContextWindow.ContextWindowSize
		}
		if input.ContextWindow.TotalInputTokens != nil {
			snapshot.Context.InputTokens = *input.ContextWindow.TotalInputTokens
		}
		if input.ContextWindow.TotalOutputTokens != nil {
			snapshot.Context.OutputTokens = *input.ContextWindow.TotalOutputTokens
		}
	}
	if input.RateLimits != nil {
		snapshot.RateLimits = append(snapshot.RateLimits[:0], input.RateLimits...)
	}
	if input.Agent != nil && input.Agent.Name != nil {
		snapshot.Agent = *input.Agent.Name
	}
	snapshot.UpdatedAt = timestamp()
	if err := writeSnapshot(dir, *snapshot); err != nil {
		return err
	}
	hasContextPercentage := snapshot.Context != nil && snapshot.Context.UsedPercentage != 0
	if input.ContextWindow != nil && input.ContextWindow.UsedPercentage != nil {
		hasContextPercentage = true
	}
	hasCost := snapshot.Cost != nil && snapshot.Cost.TotalCostUSD != 0
	if input.Cost != nil && input.Cost.TotalCostUSD != nil {
		hasCost = true
	}
	return writeStatusline(w, *snapshot, hasContextPercentage, hasCost)
}

func Prune(dir string, days int) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	cutoff := time.Now().Add(-time.Duration(days) * 24 * time.Hour)
	pruned := 0
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		data, readErr := os.ReadFile(path)
		var snapshot Snapshot
		parseErr := readErr
		old := false
		if parseErr == nil {
			parseErr = json.Unmarshal(data, &snapshot)
			if parseErr == nil {
				if snapshot.UpdatedAt == "" {
					parseErr = errors.New("missing updated_at")
				} else {
					var updatedAt time.Time
					updatedAt, parseErr = time.Parse(time.RFC3339, snapshot.UpdatedAt)
					if parseErr == nil {
						old = updatedAt.Before(cutoff)
					}
					if parseErr == nil && !old {
						continue
					}
				}
			}
		}
		if parseErr != nil || old {
			if err := os.Remove(path); err != nil {
				return pruned, err
			}
			pruned++
		}
	}
	return pruned, nil
}

func SnapshotDir() string {
	if dir := os.Getenv("DEV_HARNESS_SNAPSHOT_DIR"); dir != "" {
		return dir
	}
	if configDir := os.Getenv("CLAUDE_CONFIG_DIR"); configDir != "" {
		return filepath.Join(configDir, "dev-harness", "sessions")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".claude", "dev-harness", "sessions")
	}
	return filepath.Join(home, ".claude", "dev-harness", "sessions")
}

func decodeJSON(r io.Reader, dst any) error {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(dst); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return errors.New("input contains more than one JSON value")
		}
		return err
	}
	return nil
}

func loadSnapshot(dir, sessionID string) (*Snapshot, bool, error) {
	data, err := os.ReadFile(filepath.Join(dir, sessionID+".json"))
	if errors.Is(err, os.ErrNotExist) {
		return &Snapshot{Schema: 1, SessionID: sessionID}, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, false, err
	}
	if snapshot.SessionID == "" {
		snapshot.SessionID = sessionID
	}
	return &snapshot, true, nil
}

func writeSnapshot(dir string, snapshot Snapshot) error {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	data, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	temporary, err := os.CreateTemp(dir, ".snapshot-*.tmp")
	if err != nil {
		return err
	}
	temporaryName := temporary.Name()
	defer os.Remove(temporaryName)
	if err := temporary.Chmod(0600); err != nil {
		temporary.Close()
		return err
	}
	if _, err := temporary.Write(data); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return os.Rename(temporaryName, filepath.Join(dir, snapshot.SessionID+".json"))
}

func taskFromInput(source taskInput) Task {
	task := Task{}
	if source.ID != nil {
		task.ID = *source.ID
	}
	if source.Name != nil {
		task.Name = *source.Name
	}
	if source.Type != nil {
		task.Type = *source.Type
	}
	if source.Status != nil {
		task.Status = *source.Status
	}
	if source.Model != nil {
		task.Model = *source.Model
	}
	if source.Effort != nil {
		task.Effort = *source.Effort
	}
	if source.TokenCount != nil {
		task.TokenCount = *source.TokenCount
	}
	if source.ContextWindowSize != nil {
		task.ContextWindowSize = *source.ContextWindowSize
	}
	if source.StartTime != nil {
		task.StartTime = *source.StartTime
	}
	return task
}

func taskContent(source taskInput) string {
	parts := make([]string, 0, 4)
	if source.Name != nil {
		parts = append(parts, *source.Name)
	}
	if source.Model != nil {
		parts = append(parts, *source.Model)
	}
	if source.TokenCount != nil {
		parts = append(parts, fmt.Sprintf("%d tok", *source.TokenCount))
	}
	if source.Status != nil {
		parts = append(parts, *source.Status)
	}
	return strings.Join(parts, " · ")
}

func writeStatusline(w io.Writer, snapshot Snapshot, hasContextPercentage, hasCost bool) error {
	parts := make([]string, 0, 4)
	if snapshot.Model != nil {
		name := snapshot.Model.DisplayName
		if name == "" {
			name = snapshot.Model.ID
		}
		if name != "" {
			parts = append(parts, name)
		}
	}
	if snapshot.Context != nil && hasContextPercentage {
		parts = append(parts, fmt.Sprintf("ctx %.0f%%", math.Round(snapshot.Context.UsedPercentage)))
	}
	if snapshot.Cost != nil && hasCost {
		parts = append(parts, fmt.Sprintf("$%.2f", snapshot.Cost.TotalCostUSD))
	}
	if snapshot.Tasks != nil {
		active := 0
		for _, task := range snapshot.Tasks {
			switch task.Status {
			case "done", "completed", "failed", "cancelled":
			default:
				active++
			}
		}
		parts = append(parts, fmt.Sprintf("%d agents active", active))
	}
	if len(parts) == 0 {
		return nil
	}
	_, err := fmt.Fprintln(w, strings.Join(parts, " · "))
	return err
}

func timestamp() string {
	return time.Now().Format(time.RFC3339)
}
