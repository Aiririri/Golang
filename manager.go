package main

import "fmt"


type Analyzer interface {
	Validate(input string) error
	Analyze(input string) (Result, error)
}


type Manager struct {
	tasks    map[string]*Task
	analyzer Analyzer
	lastID   int
}

func NewManager(a Analyzer) *Manager {
	return &Manager{tasks: make(map[string]*Task), analyzer: a}
}

func (m *Manager) Create(input string) (*Task, error) {
	if err := m.analyzer.Validate(input); err != nil {
		return nil, err
	}
	m.lastID++
	t := NewTask(fmt.Sprintf("task-%d", m.lastID), input)
	m.tasks[t.ID] = t
	return t, nil
}


func (m *Manager) Get(id string) (*Task, error) {
	t, ok := m.tasks[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, id)
	}
	return t, nil
}

func (m *Manager) Process(id string) error {
	t, err := m.Get(id)
	if err != nil {
		return err
	}
	if err := t.moveTo(StatusProcessing); err != nil {
		return err
	}

	res, err := m.analyzer.Analyze(t.Input)
	if err != nil {
		t.Error = err.Error()
		_ = t.moveTo(StatusFailed) 
		return fmt.Errorf("processing %s: %w", id, err)
	}

	t.Result = res
	return t.moveTo(StatusCompleted)
}
