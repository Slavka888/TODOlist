package todo

import "errors"

var TaskNotFound = errors.New("Task not found")
var TaskAlreadyExists = errors.New("Task already exists")