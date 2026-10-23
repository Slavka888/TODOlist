package http

import (
	"TODOlist/todo"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type HTTPHandlers struct {
	todoList *todo.List
}

func NewHTTPHandlers(todoList *todo.List) *HTTPHandlers {
	return &HTTPHandlers{
		todoList: todoList,
	}
}

/*
pattern: /tasks
method: POST
*/
func (h *HTTPHandlers) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	taskDTO := TaskDTO{}

	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		stringErrDTO, _ := json.Marshal(errDTO)
		http.Error(
			w,
			string(stringErrDTO),
			http.StatusBadRequest,
		)

		return
	}

	if err := taskDTO.ValidateForCreate(); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		stringErrDTO, _ := json.Marshal(errDTO)
		http.Error(
			w,
			string(stringErrDTO),
			http.StatusBadRequest,
		)

		return
	}

	todoTask := todo.NewTask(
		taskDTO.Title,
		taskDTO.Description,
	)

	if err := h.todoList.AddTask(todoTask); err != nil {
		ErrorHandler(w, err)
		return
	}

	jsonTask, _ := json.Marshal(todoTask)
	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write(jsonTask); err != nil {
		fmt.Println("Error writing response:", err)
		return
	}
}

/*
pattern: /tasks/{title}
method: GET
*/
func (h *HTTPHandlers) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	title := r.PathValue("title")

	task, err := h.todoList.GetTask(title)
	if err != nil {
		ErrorHandler(w, err)
		return
	}

	jsonTask, _ := json.Marshal(task)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonTask); err != nil {
		fmt.Println("Error writing response:", err)
		return
	}
}

/*
pattern: /tasks?completes=*
method: GET
*/
func (h *HTTPHandlers) HandleTasksRouting(w http.ResponseWriter, r *http.Request) {
	completed := r.URL.Query().Get("completed")

	switch completed {
	case "true":
		h.HandleGetCompletedTasks(w, r)
	case "false":
		h.HandleGetNotCompletedTasks(w, r)
	default:
		h.HandleGetAllTasks(w, r)
	}
}

func (h *HTTPHandlers) HandleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.todoList.ListTasks()
	jsonTasks, _ := json.Marshal(tasks)
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write((jsonTasks)); err != nil {
		fmt.Println("Error writing response:", err)
		return
	}
}

func (h *HTTPHandlers) HandleGetCompletedTasks(w http.ResponseWriter, r *http.Request) {
	completedTasks := h.todoList.ListCompletedTasks()
	jsonTasks, _ := json.Marshal(completedTasks)
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(jsonTasks); err != nil {
		fmt.Println("Error writing response:", err)
		return
	}
}

func (h *HTTPHandlers) HandleGetNotCompletedTasks(w http.ResponseWriter, r *http.Request) {
	notCompletedTasks := h.todoList.ListNotCompletedTasks()
	jsonTasks, _ := json.Marshal(notCompletedTasks)
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(jsonTasks); err != nil {
		fmt.Println("Error writing response:", err)
		return
	}
}

/*
pattern: /tasks/{title}
method: PATCH
*/
func (h *HTTPHandlers) HandleCompleteTask(w http.ResponseWriter, r *http.Request) {
	title := r.PathValue("title")
	var completeDTO CompleteDTO

	if err := json.NewDecoder(r.Body).Decode(&completeDTO); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		stringErrDTO, _ := json.Marshal(errDTO)
		http.Error(
			w,
			string(stringErrDTO),
			http.StatusBadRequest,
		)

		return
	}

	if completeDTO.Completed {
		if taskUpdated, err := h.todoList.CompleteTask(title); err != nil {
			ErrorHandler(w, err)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			jsonTask, _ := json.Marshal(taskUpdated)
			if _, err := w.Write(jsonTask); err != nil {
				fmt.Println("Error writing response:", err)
				return
			}
		}
	} else {
		if taskUpdated, err := h.todoList.MakeActive(title); err != nil {
			ErrorHandler(w, err)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			jsonTask, _ := json.Marshal(taskUpdated)
			if _, err := w.Write(jsonTask); err != nil {
				fmt.Println("Error writing response:", err)
				return
			}
		}
	}
}

/*
pattern: /tasks/{title}
method: DELETE
*/
func (h *HTTPHandlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	title := r.PathValue("title")

	if err := h.todoList.DeleteTask(title); err != nil {
		ErrorHandler(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func ErrorHandler(w http.ResponseWriter, err error) {
	errDTO := ErrorDTO{
		Message: err.Error(),
		Time:    time.Now(),
	}

	if errors.Is(err, todo.TaskNotFound) {
		stringErrDTO, _ := json.Marshal(errDTO)

		http.Error(
			w,
			string(stringErrDTO),
			http.StatusNotFound,
		)
	} else {
		stringErrDTO, _ := json.Marshal(errDTO)
		http.Error(
			w,
			string(stringErrDTO),
			http.StatusInternalServerError,
		)
	}
}
