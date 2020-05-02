package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/blendermux/server/controllers"
	"github.com/blendermux/server/dependencies"

	"github.com/julienschmidt/httprouter"
)

func main() {
	resolver := dependencies.InitDependencyResolver()
	router := configureRoutes(resolver)

	//create the server
	server := &http.Server{
		Addr: ":8443",
		//TLSConfig: tlsConfig(),
		Handler: router,
	}

	fmt.Println("Server is running on port", server.Addr)

	//run the server
	log.Fatal(server.ListenAndServeTLS("server/cert/public.crt", "server/cert/private.key"))
}

func configureRoutes(resolver *dependencies.DependencyResolver) *httprouter.Router {
	router := httprouter.New()

	accountCon := controllers.AccountController{resolver.UserCRUD}
	router.POST("/account", accountCon.PostAccount)

	sessionCon := controllers.SessionController{resolver.UserCRUD}
	router.POST("/login", sessionCon.PostLogin)

	return router
}

func tlsConfig() *tls.Config {
	//load the client's certificate
	cert, err := ioutil.ReadFile("render_client/cert/public.crt")
	if err != nil {
		log.Fatal(err)
	}

	//create a certificate pool and add the client's to it
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)

	return &tls.Config{
		ClientCAs:  certPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
}
