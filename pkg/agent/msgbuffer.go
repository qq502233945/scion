// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package agent

import (
	"strings"
	"sync"
	"time"

	"github.com/GoogleCloudPlatform/scion/pkg/util"
)

// MessageBuffer implements a debounce-style buffer for agent message delivery.
//
// Because messages are delivered to agents by effectively "typing" them into the
// agent's terminal input, a rapid burst of arriving messages can cause contention
// and garbled input. The buffer introduces a short delay window to coalesce
// consecutive messages into a single delivery.
//
// Behavior:
//   - When a message arrives for an agent, a 2-second timer starts.
//   - If additional messages arrive before the timer fires, they are appended
//     to the buffer and the timer is reset (debounce).
//   - When the timer finally fires (2 seconds after the LAST message), all
//     buffered messages are concatenated and delivered as a single string.
//   - Interrupt messages bypass the buffer entirely for immediate delivery.
type MessageBuffer struct {
	// bufferDelay is the debounce window duration. Each new message resets
	// the timer to this duration from the current time.
	bufferDelay time.Duration

	// deliverFunc is the callback that performs actual message delivery via tmux.
	// It receives the agent ID, the concatenated message text, and the interrupt flag.
	deliverFunc func(agentID string, message string, interrupt bool) error

	mu      sync.Mutex
	buffers map[string]*agentBuffer // keyed by agent ID
}

// agentBuffer holds the pending messages and timer for a single agent.
type agentBuffer struct {
	messages []string    // accumulated messages waiting for delivery
	timer    *time.Timer // debounce timer; fires to trigger delivery
}

// NewMessageBuffer creates a new MessageBuffer with the given debounce delay
// and delivery function. The deliverFunc is called asynchronously when the
// buffer flushes — it should perform the actual tmux send-keys delivery.
func NewMessageBuffer(delay time.Duration, deliverFunc func(agentID string, message string, interrupt bool) error) *MessageBuffer {
	return &MessageBuffer{
		bufferDelay: delay,
		deliverFunc: deliverFunc,
		buffers:     make(map[string]*agentBuffer),
	}
}

// Send queues a message for buffered delivery to the given agent.
// The message is added to the agent's buffer and the debounce timer is
// started (or reset if already running). The actual delivery occurs
// asynchronously once the timer fires.
func (mb *MessageBuffer) Send(agentID string, message string) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	buf, exists := mb.buffers[agentID]
	if !exists {
		// First message for this agent — create a new buffer entry.
		buf = &agentBuffer{}
		mb.buffers[agentID] = buf
	}

	// Append the message to the pending list.
	buf.messages = append(buf.messages, message)
	util.Debugf("msgbuffer: queued message for agent %s (%d pending)", agentID, len(buf.messages))

	// Reset or start the debounce timer. If a timer is already running,
	// stop it first so we can restart with a fresh delay window.
	if buf.timer != nil {
		buf.timer.Stop()
	}
	buf.timer = time.AfterFunc(mb.bufferDelay, func() {
		mb.flush(agentID)
	})
}

// flush delivers all buffered messages for the given agent as a single
// concatenated string. Called when the debounce timer fires.
func (mb *MessageBuffer) flush(agentID string) {
	mb.mu.Lock()
	buf, exists := mb.buffers[agentID]
	if !exists || len(buf.messages) == 0 {
		mb.mu.Unlock()
		return
	}

	// Take ownership of the pending messages and clean up the buffer entry.
	pending := buf.messages
	delete(mb.buffers, agentID)
	mb.mu.Unlock()

	// Concatenate all buffered messages with double-newline separators so
	// each original message remains visually distinct in the agent's input.
	combined := strings.Join(pending, "\n\n")
	util.Debugf("msgbuffer: flushing %d message(s) for agent %s", len(pending), agentID)

	if err := mb.deliverFunc(agentID, combined, false); err != nil {
		// Delivery errors are logged rather than returned since flush runs
		// asynchronously after the original Send() call has already returned.
		util.Debugf("msgbuffer: delivery error for agent %s: %v", agentID, err)
	}
}

// Close flushes all pending buffers immediately and stops all timers.
// Call this during shutdown to ensure no messages are lost.
func (mb *MessageBuffer) Close() {
	mb.mu.Lock()
	agentIDs := make([]string, 0, len(mb.buffers))
	for id, buf := range mb.buffers {
		if buf.timer != nil {
			buf.timer.Stop()
		}
		agentIDs = append(agentIDs, id)
	}
	mb.mu.Unlock()

	// Flush each agent's pending messages outside the lock.
	for _, id := range agentIDs {
		mb.flush(id)
	}
}
