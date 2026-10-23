package http

import (
	"errors"
	"time"
)

type TaskDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (t TaskDTO) ValidateForCreate() error {
	if t.Title == "" {
		return errors.New("title is empty")
	}
	if t.Description == "" {
		return errors.New("description is empty")
	}

	return nil
}

type ErrorDTO struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

type CompleteDTO struct {
	Completed bool `json:"completed"`
}
