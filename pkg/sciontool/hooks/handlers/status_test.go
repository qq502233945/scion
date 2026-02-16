/*
Copyright 2025 The Scion Authors.
*/

package handlers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ptone/scion-agent/pkg/sciontool/hooks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusHandler_UpdateStatus(t *testing.T) {
	// Create temp dir
	tmpDir := t.TempDir()
	statusPath := filepath.Join(tmpDir, "agent-info.json")

	h := &StatusHandler{
		StatusPath: statusPath,
	}

	// Test updating status
	err := h.UpdateStatus(hooks.StateThinking, false)
	require.NoError(t, err)

	// Verify file contents
	data, err := os.ReadFile(statusPath)
	require.NoError(t, err)

	var info AgentInfo
	err = json.Unmarshal(data, &info)
	require.NoError(t, err)
	assert.Equal(t, "THINKING", info.Status)
	assert.Empty(t, info.SessionStatus)

	// Test updating session status
	err = h.UpdateStatus(hooks.StateWaitingForInput, true)
	require.NoError(t, err)

	data, err = os.ReadFile(statusPath)
	require.NoError(t, err)
	err = json.Unmarshal(data, &info)
	require.NoError(t, err)
	assert.Equal(t, "THINKING", info.Status) // Previous status preserved
	assert.Equal(t, "WAITING_FOR_INPUT", info.SessionStatus)
}

func TestStatusHandler_Handle(t *testing.T) {
	tmpDir := t.TempDir()
	statusPath := filepath.Join(tmpDir, "agent-info.json")

	h := &StatusHandler{
		StatusPath: statusPath,
	}

	tests := []struct {
		name       string
		event      *hooks.Event
		wantStatus hooks.AgentState
	}{
		{
			name:       "SessionStart sets STARTING",
			event:      &hooks.Event{Name: hooks.EventSessionStart},
			wantStatus: hooks.StateStarting,
		},
		{
			name:       "PreStart sets INITIALIZING",
			event:      &hooks.Event{Name: hooks.EventPreStart},
			wantStatus: hooks.StateInitializing,
		},
		{
			name:       "PostStart sets IDLE",
			event:      &hooks.Event{Name: hooks.EventPostStart},
			wantStatus: hooks.StateIdle,
		},
		{
			name:       "PreStop sets SHUTTING_DOWN",
			event:      &hooks.Event{Name: hooks.EventPreStop},
			wantStatus: hooks.StateShuttingDown,
		},
		{
			name:       "PromptSubmit sets THINKING",
			event:      &hooks.Event{Name: hooks.EventPromptSubmit},
			wantStatus: hooks.StateThinking,
		},
		{
			name:       "ToolStart sets EXECUTING",
			event:      &hooks.Event{Name: hooks.EventToolStart, Data: hooks.EventData{ToolName: "Bash"}},
			wantStatus: hooks.StateExecuting,
		},
		{
			name:       "ToolEnd sets IDLE",
			event:      &hooks.Event{Name: hooks.EventToolEnd},
			wantStatus: hooks.StateIdle,
		},
		{
			name:       "AgentEnd sets IDLE",
			event:      &hooks.Event{Name: hooks.EventAgentEnd},
			wantStatus: hooks.StateIdle,
		},
		{
			name:       "SessionEnd sets EXITED",
			event:      &hooks.Event{Name: hooks.EventSessionEnd},
			wantStatus: hooks.StateExited,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := h.Handle(tt.event)
			require.NoError(t, err)

			data, err := os.ReadFile(statusPath)
			require.NoError(t, err)

			var info AgentInfo
			err = json.Unmarshal(data, &info)
			require.NoError(t, err)
			assert.Equal(t, string(tt.wantStatus), info.Status)
		})
	}
}

func TestStatusHandler_ClearWaitingStatus(t *testing.T) {
	tmpDir := t.TempDir()
	statusPath := filepath.Join(tmpDir, "agent-info.json")

	h := &StatusHandler{StatusPath: statusPath}

	// Set sessionStatus to WAITING_FOR_INPUT
	err := h.UpdateStatus(hooks.StateWaitingForInput, true)
	require.NoError(t, err)

	// Verify it's set
	info := readAgentInfo(t, statusPath)
	assert.Equal(t, "WAITING_FOR_INPUT", info.SessionStatus)

	// Clear it
	err = h.ClearWaitingStatus()
	require.NoError(t, err)

	info = readAgentInfo(t, statusPath)
	assert.Empty(t, info.SessionStatus)
}

func TestStatusHandler_ClearWaitingStatus_NoOpWhenNotWaiting(t *testing.T) {
	tmpDir := t.TempDir()
	statusPath := filepath.Join(tmpDir, "agent-info.json")

	h := &StatusHandler{StatusPath: statusPath}

	// Set sessionStatus to COMPLETED
	err := h.UpdateStatus(hooks.StateCompleted, true)
	require.NoError(t, err)

	// ClearWaitingStatus should not clear COMPLETED
	err = h.ClearWaitingStatus()
	require.NoError(t, err)

	info := readAgentInfo(t, statusPath)
	assert.Equal(t, "COMPLETED", info.SessionStatus)
}

func TestStatusHandler_ClearWaitingStatus_NoOpWhenEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	statusPath := filepath.Join(tmpDir, "agent-info.json")

	h := &StatusHandler{StatusPath: statusPath}

	// Set operational status only
	err := h.UpdateStatus(hooks.StateThinking, false)
	require.NoError(t, err)

	// ClearWaitingStatus should be a no-op
	err = h.ClearWaitingStatus()
	require.NoError(t, err)

	info := readAgentInfo(t, statusPath)
	assert.Empty(t, info.SessionStatus)
}

func TestStatusHandler_Handle_ClearsWaitingOnActivity(t *testing.T) {
	activityEvents := []struct {
		name  string
		event *hooks.Event
	}{
		{
			name:  "ToolStart clears waiting",
			event: &hooks.Event{Name: hooks.EventToolStart, Data: hooks.EventData{ToolName: "Bash"}},
		},
		{
			name:  "PromptSubmit clears waiting",
			event: &hooks.Event{Name: hooks.EventPromptSubmit},
		},
		{
			name:  "AgentStart clears waiting",
			event: &hooks.Event{Name: hooks.EventAgentStart},
		},
	}

	for _, tt := range activityEvents {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			statusPath := filepath.Join(tmpDir, "agent-info.json")
			h := &StatusHandler{StatusPath: statusPath}

			// Pre-set sessionStatus to WAITING_FOR_INPUT
			err := h.UpdateStatus(hooks.StateWaitingForInput, true)
			require.NoError(t, err)

			// Handle the activity event
			err = h.Handle(tt.event)
			require.NoError(t, err)

			info := readAgentInfo(t, statusPath)
			assert.Empty(t, info.SessionStatus, "sessionStatus should be cleared")
		})
	}
}

func TestStatusHandler_Handle_DoesNotClearCompletedOnActivity(t *testing.T) {
	tmpDir := t.TempDir()
	statusPath := filepath.Join(tmpDir, "agent-info.json")
	h := &StatusHandler{StatusPath: statusPath}

	// Pre-set sessionStatus to COMPLETED
	err := h.UpdateStatus(hooks.StateCompleted, true)
	require.NoError(t, err)

	// Handle a tool-start event
	err = h.Handle(&hooks.Event{
		Name: hooks.EventToolStart,
		Data: hooks.EventData{ToolName: "Bash"},
	})
	require.NoError(t, err)

	info := readAgentInfo(t, statusPath)
	assert.Equal(t, "COMPLETED", info.SessionStatus, "COMPLETED should not be cleared")
}

func TestStatusHandler_Handle_ToolEndDoesNotClearWaiting(t *testing.T) {
	tmpDir := t.TempDir()
	statusPath := filepath.Join(tmpDir, "agent-info.json")
	h := &StatusHandler{StatusPath: statusPath}

	// Pre-set sessionStatus to WAITING_FOR_INPUT
	err := h.UpdateStatus(hooks.StateWaitingForInput, true)
	require.NoError(t, err)

	// Handle a tool-end event (should NOT clear)
	err = h.Handle(&hooks.Event{Name: hooks.EventToolEnd})
	require.NoError(t, err)

	info := readAgentInfo(t, statusPath)
	assert.Equal(t, "WAITING_FOR_INPUT", info.SessionStatus, "tool-end should not clear waiting")
}

func TestStatusHandler_Handle_ClaudeExitPlanMode(t *testing.T) {
	tmpDir := t.TempDir()
	statusPath := filepath.Join(tmpDir, "agent-info.json")
	h := &StatusHandler{StatusPath: statusPath}

	// Handle ExitPlanMode tool-start from Claude dialect
	err := h.Handle(&hooks.Event{
		Name:    hooks.EventToolStart,
		Dialect: "claude",
		Data:    hooks.EventData{ToolName: "ExitPlanMode"},
	})
	require.NoError(t, err)

	info := readAgentInfo(t, statusPath)
	assert.Equal(t, "EXECUTING", info.Status)
	assert.Equal(t, "WAITING_FOR_INPUT", info.SessionStatus)
}

func TestStatusHandler_Handle_ClaudeAskUserQuestion(t *testing.T) {
	tmpDir := t.TempDir()
	statusPath := filepath.Join(tmpDir, "agent-info.json")
	h := &StatusHandler{StatusPath: statusPath}

	// Pre-set sessionStatus to WAITING_FOR_INPUT (simulating sciontool status ask_user)
	err := h.UpdateStatus(hooks.StateWaitingForInput, true)
	require.NoError(t, err)

	// Handle AskUserQuestion tool-start from Claude dialect
	err = h.Handle(&hooks.Event{
		Name:    hooks.EventToolStart,
		Dialect: "claude",
		Data:    hooks.EventData{ToolName: "AskUserQuestion"},
	})
	require.NoError(t, err)

	info := readAgentInfo(t, statusPath)
	assert.Equal(t, "WAITING_FOR_INPUT", info.SessionStatus, "AskUserQuestion should maintain WAITING_FOR_INPUT")
}

func TestStatusHandler_Handle_NonClaudeExitPlanModeIgnored(t *testing.T) {
	tmpDir := t.TempDir()
	statusPath := filepath.Join(tmpDir, "agent-info.json")
	h := &StatusHandler{StatusPath: statusPath}

	// Handle ExitPlanMode from a non-claude dialect — should NOT set sessionStatus
	err := h.Handle(&hooks.Event{
		Name:    hooks.EventToolStart,
		Dialect: "gemini",
		Data:    hooks.EventData{ToolName: "ExitPlanMode"},
	})
	require.NoError(t, err)

	info := readAgentInfo(t, statusPath)
	assert.Empty(t, info.SessionStatus, "non-claude ExitPlanMode should not set sessionStatus")
}

func TestStatusHandler_Handle_ClaudeExitPlanModeThenActivity(t *testing.T) {
	tmpDir := t.TempDir()
	statusPath := filepath.Join(tmpDir, "agent-info.json")
	h := &StatusHandler{StatusPath: statusPath}

	// ExitPlanMode sets WAITING_FOR_INPUT
	err := h.Handle(&hooks.Event{
		Name:    hooks.EventToolStart,
		Dialect: "claude",
		Data:    hooks.EventData{ToolName: "ExitPlanMode"},
	})
	require.NoError(t, err)

	info := readAgentInfo(t, statusPath)
	assert.Equal(t, "WAITING_FOR_INPUT", info.SessionStatus)

	// Tool-end for ExitPlanMode should NOT clear it
	err = h.Handle(&hooks.Event{Name: hooks.EventToolEnd, Dialect: "claude"})
	require.NoError(t, err)

	info = readAgentInfo(t, statusPath)
	assert.Equal(t, "WAITING_FOR_INPUT", info.SessionStatus)

	// User approves plan, next tool starts — should clear
	err = h.Handle(&hooks.Event{
		Name:    hooks.EventToolStart,
		Dialect: "claude",
		Data:    hooks.EventData{ToolName: "Bash"},
	})
	require.NoError(t, err)

	info = readAgentInfo(t, statusPath)
	assert.Empty(t, info.SessionStatus, "activity after plan approval should clear WAITING_FOR_INPUT")
}

// readAgentInfo is a test helper that reads and parses agent-info.json.
func readAgentInfo(t *testing.T, path string) AgentInfo {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)

	var info AgentInfo
	err = json.Unmarshal(data, &info)
	require.NoError(t, err)
	return info
}
