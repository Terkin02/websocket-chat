package chat

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)


func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}


func RunChat(conn *websocket.Conn, inputChan <-chan string, done chan interface{}) {
	// Горутина чтения входящих сообщений от сервера
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read from server error:", err)
				close(done)
				return
			}
			fmt.Printf("Received: %s", string(msg))
		}
	}()

	// Отправка вводимых пользователем строк на сервер
	for msg := range inputChan {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			log.Println("Write to server error:", err)
			close(done)
			return
		}
	}
}
