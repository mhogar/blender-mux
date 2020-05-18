package router

import (
	"blendermux/server/controllers"

	"github.com/julienschmidt/httprouter"
)

// CreateRouter creates a new router with the endpoints and panic handler configured
func CreateRouter(handler controllers.RequestHandler) *httprouter.Router {
	router := httprouter.New()

	router.PanicHandler = controllers.PanicHandler

	router.POST("/user", handler.PostUser)
	router.DELETE("/user/:id", handler.DeleteUser)

	return router
}
