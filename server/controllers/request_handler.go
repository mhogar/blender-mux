package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// RequestHandler is an interface that encapsulates all other handler interfaces
type RequestHandler interface {
	UserHandler
}

// UserHandler is an interface for handling requests to user routes
type UserHandler interface {
	// PostUser handles Post requests to "/user"
	PostUser(http.ResponseWriter, *http.Request, httprouter.Params)

	// DeleteUser handles Post requests to "/user"
	DeleteUser(http.ResponseWriter, *http.Request, httprouter.Params)
}

// RequestHandle is an implementation of RequestHandler that uses controllers to satisfy the interface's methods
type RequestHandle struct {
	UserController
}
