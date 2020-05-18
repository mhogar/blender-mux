package main

import (
	"fmt"
	"log"
	"net/http"

	"blendermux/server/config"
	"blendermux/server/controllers"
	"blendermux/server/dependencies"

	"github.com/julienschmidt/httprouter"
)

func main() {
	config.InitConfig()
	router := configureRoutes()

	//create the server
	server := &http.Server{
		Addr:    ":8443",
		Handler: router,
	}

	fmt.Println("Server is running on port", server.Addr)

	//run the server
	log.Fatal(server.ListenAndServeTLS("server/cert/public.crt", "server/cert/private.key"))
}

func configureRoutes() *httprouter.Router {
	router := httprouter.New()

	userCon := controllers.UserController{
		UserCRUD: dependencies.ResolveDatabase(),
	}
	router.POST("/user", controllers.CreateRequestHandler(userCon.PostUser))

	// sessionCon := controllers.SessionController{
	// 	UserCRUD:    dependencies.ResolveDatabase(),
	// 	SessionCRUD: dependencies.ResolveDatabase(),
	// }
	// router.POST("/login", sessionCon.PostLogin)
	// router.POST("/logout", sessionCon.PostLogout)

	return router
}
