package main

import (
	"errors"
	"fmt"
)


var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrNotFound          = errors.New("task not found")
	ErrInvalidTransition = errors.New("invalid state transition")
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

type Result struct {
	Scheme      string
	Host        string
	Port        int 
	IsSecure    bool
	Path        string
	PathDepth   int 
	QueryParams int 
	HasFragment bool
}

type Task struct {
	ID     string
	Input  string
	Status Status
	Result Result 
	Error  string 
}


func NewTask(id, input string) *Task {
	return &Task{ID: id, Input: input, Status: StatusPending}
}


func canMove(from, to Status) bool {
	switch from {
	case StatusPending:
		return to == StatusProcessing
	case StatusProcessing:
		return to == StatusCompleted || to == StatusFailed
	}
	return false
}

func (t *Task) moveTo(next Status) error {
	if !canMove(t.Status, next) {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidTransition, t.Status, next)
	}
	t.Status = next
	return nil
}
