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

	dialer := websocket.DefaultDialer
	dialer.TLSClientConfig = tlsConfig()

	conn, _, err := dialer.Dial(url.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		fmt.Println(string(message))
	}
}

func tlsConfig() *tls.Config {
	cert, err := ioutil.ReadFile("cert/public.crt")
	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)

	return &tls.Config{
		RootCAs: certPool,
	}
}
