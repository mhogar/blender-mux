package main

import (
	"blendermux/server/dependencies"
	"blendermux/server/router"
	"fmt"
	"log"
	"net/http"

	"blendermux/server/config"
)

func main() {
	config.InitConfig()

	//create the server
	server := &http.Server{
		Addr:    ":8443",
		Handler: router.CreateRouter(dependencies.ResolveRouteHandler()),
	}

	fmt.Println("Server is running on port", server.Addr)

	//run the server
	log.Fatal(server.ListenAndServeTLS("cert/public.crt", "cert/private.key"))
}
