package http

import (
	"errors"
	"net/http"
)

type HTTPServer struct {
	httpHandlers *HTTPHandlers
}

func NewHTTPServer(httpHandlers *HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: httpHandlers,
	}
}

func (s *HTTPServer) Start() error {
	router := http.NewServeMux()
	router.HandleFunc("POST /tasks", s.httpHandlers.HandleCreateTask)
	router.HandleFunc("GET /tasks/{title}", s.httpHandlers.HandleGetTask)
	router.HandleFunc("GET /tasks", s.httpHandlers.HandleTasksRouting)
	router.HandleFunc("PATCH /tasks/{title}", s.httpHandlers.HandleCompleteTask)
	router.HandleFunc("DELETE /tasks/{title}", s.httpHandlers.HandleDeleteTask)

	if err := http.ListenAndServe(":8080", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	}

	return nil
}

