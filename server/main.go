package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/blendermux/server/controllers"
	"github.com/blendermux/server/dependencies"

	"github.com/julienschmidt/httprouter"
)

func main() {
	resolver := dependencies.CreateDependencyResolver()
	router := configureRoutes(resolver)

	//create the server
	server := &http.Server{
		Addr:    ":8443",
		Handler: router,
	}

	fmt.Println("Server is running on port", server.Addr)

	//run the server
	log.Fatal(server.ListenAndServeTLS("server/cert/public.crt", "server/cert/private.key"))
}

func configureRoutes(resolver dependencies.DependencyResolver) *httprouter.Router {
	router := httprouter.New()

	accountCon := controllers.AccountController{
		UserCRUD: resolver.Database,
	}
	router.POST("/account", accountCon.PostAccount)

	sessionCon := controllers.SessionController{
		UserCRUD:    resolver.Database,
		SessionCRUD: resolver.Database,
	}
	router.POST("/login", sessionCon.PostLogin)
	router.POST("/logout", sessionCon.PostLogout)

	return router
}
