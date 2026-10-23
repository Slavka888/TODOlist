package main

import (
	"TODOlist/http"
	"TODOlist/todo"
	"fmt"
)

func main() {
	todoList := todo.NewList()
	httpHandlers := http.NewHTTPHandlers(todoList)
	httpServer := http.NewHTTPServer(httpHandlers)

	if err := httpServer.Start(); err != nil {
		fmt.Println("Failed to start HTTP server:", err)
	}
}
