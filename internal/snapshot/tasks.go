package snapshot

import (
	"fmt"
	"strings"
)

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
