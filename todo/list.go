package todo

import (
	"maps"
	"sync"
)

type List struct {
	tasks map[string]Task
	mtx   sync.RWMutex
}

func NewList() *List {
	return &List{
		tasks: make(map[string]Task),
	}
}

func (l *List) AddTask(task Task) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if _, exist := l.tasks[task.Title]; exist {
		return TaskAlreadyExists
	}

	l.tasks[task.Title] = task

	return nil
}

func (l *List) ListTasks() map[string]Task {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	return maps.Clone(l.tasks)
}

func (l *List) GetTask(title string) (Task, error) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	if task, exist := l.tasks[title]; exist {
		return task, nil
	}

	return Task{}, TaskNotFound
}

func (l *List) ListCompletedTasks() map[string]Task {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	tmp := make(map[string]Task)

	for k, v := range l.tasks {
		if v.IsDone {
			tmp[k] = v
		}
	}

	return tmp
}

func (l *List) ListNotCompletedTasks() map[string]Task {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	tmp := make(map[string]Task)

	for k, v := range l.tasks {
		if !v.IsDone {
			tmp[k] = v
		}
	}

	return tmp
}

func (l *List) CompleteTask(title string) (Task, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if task, exist := l.tasks[title]; exist {
		task.Done()
		l.tasks[title] = task
		return task, nil
	} else {
		return Task{}, TaskNotFound
	}
}

func (l *List) MakeActive(title string) (Task, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if task, exist := l.tasks[title]; exist {
		task.MakeActive()
		l.tasks[title] = task
		return task, nil
	} else {
		return Task{}, TaskNotFound
	}
}

func (l *List) DeleteTask(title string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	
	if _, exist := l.tasks[title]; exist {
		delete(l.tasks, title)
	} else {
		return TaskNotFound
	}

	return nil
}
