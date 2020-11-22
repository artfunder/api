package api

import (
	"errors"
)

// ErrorNotFound ...
var ErrorNotFound = errors.New("Object not found")

// ErrorNoHandler ...
var ErrorNoHandler = errors.New("No handler for that route and method")
