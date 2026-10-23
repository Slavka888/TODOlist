package todo

import "time"

type Task struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsDone      bool      `json:"is_done"`
	CreatedAt   time.Time `json:"created_at"`
	DoneAt      *time.Time `json:"done_at"`
}

func NewTask(title, description string) Task {
	return Task{
		Title:       title,
		Description: description,
		IsDone:      false,
		CreatedAt:   time.Now(),
		DoneAt:      nil,
	}
}

func (t *Task) Done() {
	completionTime := time.Now()
	t.DoneAt = &completionTime
	t.IsDone = true
}

func (t *Task) MakeActive() {
	t.IsDone = false
	t.DoneAt = nil
}