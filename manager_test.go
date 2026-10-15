package main

import (
	"errors"
	"testing"
)


type failingAnalyzer struct{}

func (failingAnalyzer) Validate(string) error { return nil }
func (failingAnalyzer) Analyze(string) (Result, error) {
	return Result{}, errors.New("boom")
}

func newManager() *Manager { return NewManager(URLAnalyzer{}) }

func TestCreateSuccess(t *testing.T) {
	m := newManager()
	task, err := m.Create("https://example.com/a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task.ID == "" {
		t.Error("ID is empty")
	}
	if task.Status != StatusPending {
		t.Errorf("status = %s, want pending", task.Status)
	}
}

func TestCreateInvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"whitespace", "   "},
		{"no scheme", "example.com"},
		{"bad scheme", "ftp://example.com"},
		{"bad port", "http://example.com:99999"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := newManager()
			_, err := m.Create(tc.input)
			if !errors.Is(err, ErrInvalidInput) {
				t.Fatalf("err = %v, want ErrInvalidInput", err)
			}
			if len(m.tasks) != 0 {
				t.Fatalf("invalid task was stored, tasks = %d", len(m.tasks))
			}
		})
	}
}

func TestGet(t *testing.T) {
	m := newManager()
	created, _ := m.Create("https://example.com")

	got, err := m.Get(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != created.ID || got.Input != "https://example.com" {
		t.Fatalf("got %+v, want task %s", got, created.ID)
	}
}

func TestGetNotFound(t *testing.T) {
	_, err := newManager().Get("task-999")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestProcessSuccess(t *testing.T) {
	m := newManager()
	created, _ := m.Create("https://Example.com:8443/docs/go?lang=en&v=1")

	if err := m.Process(created.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	task, _ := m.Get(created.ID)
	if task.Status != StatusCompleted {
		t.Fatalf("status = %s, want completed", task.Status)
	}
	want := Result{
		Scheme: "https", Host: "example.com", Port: 8443, IsSecure: true,
		Path: "/docs/go", PathDepth: 2, QueryParams: 2,
	}
	if task.Result != want {
		t.Fatalf("result = %+v, want %+v", task.Result, want)
	}
}

func TestProcessNotFound(t *testing.T) {
	err := newManager().Process("task-999")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestProcessTwiceIsInvalidTransition(t *testing.T) {
	m := newManager()
	created, _ := m.Create("https://example.com")
	if err := m.Process(created.ID); err != nil {
		t.Fatal(err)
	}

	err := m.Process(created.ID)
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("err = %v, want ErrInvalidTransition", err)
	}
	task, _ := m.Get(created.ID)
	if task.Status != StatusCompleted {
		t.Fatalf("status = %s, completed task must stay completed", task.Status)
	}
}

func TestProcessFailure(t *testing.T) {
	m := NewManager(failingAnalyzer{})
	created, _ := m.Create("anything")

	if err := m.Process(created.ID); err == nil {
		t.Fatal("expected an error")
	}
	task, _ := m.Get(created.ID)
	if task.Status != StatusFailed || task.Error != "boom" {
		t.Fatalf("got status=%s error=%q, want failed / boom", task.Status, task.Error)
	}
}

func TestMultipleIndependentTasks(t *testing.T) {
	m := newManager()
	a, _ := m.Create("https://a.example.com/one")
	b, _ := m.Create("http://b.example.org:8080/x/y/z")
	c, _ := m.Create("https://c.example.net")

	if a.ID == b.ID || b.ID == c.ID || a.ID == c.ID {
		t.Fatalf("IDs are not unique: %s %s %s", a.ID, b.ID, c.ID)
	}

	
	if err := m.Process(a.ID); err != nil {
		t.Fatal(err)
	}
	if err := m.Process(b.ID); err != nil {
		t.Fatal(err)
	}

	if a.Result.Host != "a.example.com" || a.Result.PathDepth != 1 {
		t.Errorf("task a has wrong result: %+v", a.Result)
	}
	if b.Result.Host != "b.example.org" || b.Result.Port != 8080 || b.Result.PathDepth != 3 {
		t.Errorf("task b has wrong result: %+v", b.Result)
	}
	if c.Status != StatusPending || c.Result != (Result{}) {
		t.Errorf("task c must be pending with an empty result: %+v", c)
	}
}
