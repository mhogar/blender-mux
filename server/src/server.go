package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func handler(w http.ResponseWriter, req *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	writeMessage(conn, "Hello client!")
	time.Sleep(time.Second * 1)
	writeMessage(conn, "Hello again!")
	time.Sleep(time.Second * 1)
	writeMessage(conn, "Goodbye!")
	conn.WriteControl(websocket.CloseMessage, nil, time.Now().Add(time.Second))
}

func writeMessage(conn *websocket.Conn, message string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServeTLS(":8443", "server/cert/public.crt", "server/cert/private.key", nil)
}
