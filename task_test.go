package main

import (
	"errors"
	"testing"
)

func TestMoveTo(t *testing.T) {
	tests := []struct {
		name    string
		from    Status
		to      Status
		wantErr bool
	}{
		{"pending to processing", StatusPending, StatusProcessing, false},
		{"processing to completed", StatusProcessing, StatusCompleted, false},
		{"processing to failed", StatusProcessing, StatusFailed, false},
		{"pending to completed", StatusPending, StatusCompleted, true},
		{"processing to processing", StatusProcessing, StatusProcessing, true},
		{"completed to processing", StatusCompleted, StatusProcessing, true},
		{"failed to processing", StatusFailed, StatusProcessing, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			task := NewTask("task-1", "http://a.b")
			task.Status = tc.from

			err := task.moveTo(tc.to)

			if tc.wantErr {
				if !errors.Is(err, ErrInvalidTransition) {
					t.Fatalf("err = %v, want ErrInvalidTransition", err)
				}
				if task.Status != tc.from {
					t.Fatalf("status changed to %s after a rejected transition", task.Status)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if task.Status != tc.to {
				t.Fatalf("status = %s, want %s", task.Status, tc.to)
			}
		})
	}
}
