package main

import (
	"app/common"

	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	//create the server
	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig(),
	}

	//configure and run the server
	http.HandleFunc("/", handler)
	log.Fatal(server.ListenAndServeTLS("server/cert/public.crt", "server/cert/private.key"))
}

func handler(w http.ResponseWriter, req *http.Request) {
	//upgrade http server to use web sockets
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go readLoop(conn)

	//send the messages
	sendRequest(conn, common.RENDER, "This is a render message")

	//send a close message
	conn.WriteControl(websocket.CloseMessage, nil, time.Now().Add(time.Second))
}

func readLoop(conn *websocket.Conn) {
	//simple read loop to detect close messages
	for {
		_, _, err := conn.NextReader()
		if err != nil {
			conn.Close()
			break
		}
	}
}

func sendRequest(conn *websocket.Conn, requestType byte, message string) {
	length := len(message) + 1

	data := make([]byte, length, length)
	data[0] = requestType
	copy(data[1:length], []byte(message))

	err := conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		log.Println(err)
		return
	}
}

func writeMessage(conn *websocket.Conn, message string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println(err)
		return
	}
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
