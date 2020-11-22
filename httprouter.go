package api

import "github.com/julienschmidt/httprouter"

var router *httprouter.Router

// GetHTTPRouter returns the current instance of the httprouter
func GetHTTPRouter() *httprouter.Router {
	if router == nil {
		router = httprouter.New()
	}
	return router
}

// ResetHTTPRouter creates a new httprouter even if one exists
func ResetHTTPRouter() {
	router = nil
}
