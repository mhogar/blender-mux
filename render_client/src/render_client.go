package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	url := url.URL{
		Scheme: "wss",
		Host:   "localhost:8443",
		Path:   "/",
	}

	//add the tls config to the dialer
	dialer := websocket.DefaultDialer
	dialer.TLSClientConfig = tlsConfig()

	//dial the server and start the connection
	conn, _, err := dialer.Dial(url.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	//main read loop: recieve messages and print them
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		fmt.Println(string(message))
	}
}

func tlsConfig() *tls.Config {
	//create the client's certificate
	clientCert, err := tls.LoadX509KeyPair("render_client/cert/public.crt", "render_client/cert/private.key")
	if err != nil {
		log.Fatal(err)
	}

	//load the server's certificate
	serverCert, err := ioutil.ReadFile("server/cert/public.crt")
	if err != nil {
		log.Fatal(err)
	}

	//create a certificate pool and add the server's to it
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(serverCert)

	return &tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{clientCert},
	}
}
