package todo

import (
	"fmt"
	"sync"
	"time"
)

type Priority string
type Status string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"

	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Todo struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    Priority   `json:"priority"`
	Status      Status     `json:"status"`
	Tags        []string   `json:"tags"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type store struct {
	mu    sync.RWMutex
	todos map[string]*Todo
	order []string // insertion order
}

var s = &store{
	todos: make(map[string]*Todo),
}

func generateID() string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}

func GetAll() []*Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	todos := make([]*Todo, 0, len(s.order))
	for i := len(s.order) - 1; i >= 0; i-- {
		if t, ok := s.todos[s.order[i]]; ok {
			todos = append(todos, t)
		}
	}
	return todos
}

func Create(t *Todo) *Todo {
	s.mu.Lock()
	defer s.mu.Unlock()
	t.ID = generateID()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	if t.Status == "" {
		t.Status = StatusPending
	}
	if t.Priority == "" {
		t.Priority = PriorityMedium
	}
	if t.Tags == nil {
		t.Tags = []string{}
	}
	s.todos[t.ID] = t
	s.order = append(s.order, t.ID)
	return t
}

func GetByID(id string) *Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.todos[id]
}

func Update(id string, updates *Todo) *Todo {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.todos[id]
	if !ok {
		return nil
	}
	if updates.Title != "" {
		t.Title = updates.Title
	}
	if updates.Description != "" {
		t.Description = updates.Description
	}
	if updates.Priority != "" {
		t.Priority = updates.Priority
	}
	if updates.Status != "" {
		t.Status = updates.Status
	}
	if updates.Tags != nil {
		t.Tags = updates.Tags
	}
	if updates.DueDate != nil {
		t.DueDate = updates.DueDate
	}
	t.UpdatedAt = time.Now()
	return t
}

func Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.todos[id]; !ok {
		return false
	}
	delete(s.todos, id)
	for i, oid := range s.order {
		if oid == id {
			s.order = append(s.order[:i], s.order[i+1:]...)
			break
		}
	}
	return true
}
