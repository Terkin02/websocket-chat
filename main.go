package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	connMu      sync.Mutex
	connections []*websocket.Conn
)

func main() {
	http.HandleFunc("/start", SocketHandler)
	log.Println("Server listening on :8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	addConn(conn)
	defer func() {
		delConn(conn)
		conn.Close()
	}()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}
		broadcast(msgType, msg, conn)
	}
}

func addConn(c *websocket.Conn) {
	connMu.Lock()
	connections = append(connections, c)
	connMu.Unlock()
}

func delConn(c *websocket.Conn) {
	connMu.Lock()
	defer connMu.Unlock()
	for i, v := range connections {
		if v == c {
			connections[i] = connections[len(connections)-1]
			connections = connections[:len(connections)-1]
			return
		}
	}
}

func broadcast(msgType int, msg []byte, except *websocket.Conn) {
	connMu.Lock()
	defer connMu.Unlock()
	for _, c := range connections {
		if c == except {
			continue
		}
		if err := c.WriteMessage(msgType, msg); err != nil {
			log.Println("Write error:", err)
		}
	}
}
