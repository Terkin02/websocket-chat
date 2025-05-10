package main

import (
	"bufio"
	"log"
	"os"

	"chat/internal/chat"

	"github.com/gorilla/websocket"
)

func main() {
	done := make(chan interface{})

	// Адрес WebSocket-сервера
	socketURL := "ws://localhost:8888/start"
	conn, _, err := websocket.DefaultDialer.Dial(socketURL, nil)
	if err != nil {
		log.Fatal("Error connecting to websocket server:", err)
	}
	defer conn.Close()

	// Канал для ввода пользователя
	inputChan := make(chan string)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			msg, err := reader.ReadString('\n')
			if err != nil {
				log.Println("Error reading stdin:", err)
				close(inputChan)
				return
			}
			inputChan <- msg
		}
	}()

	// Очистка экрана и запуск основного цикла чата
	chat.ClearScreen()
	chat.RunChat(conn, inputChan, done)

	done <- "Hi" // ждём завершения
}
