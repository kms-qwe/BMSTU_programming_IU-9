package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func main() {
	serverAddr := "ws://185.102.139.161:8081/ws"
	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		log.Fatal("Ошибка подключения к серверу:", err)
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte("get_random"))
	if err != nil {
		log.Fatal("Ошибка при отправке сообщения:", err)
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Fatal("Ошибка при чтении сообщения:", err)
	}

	fmt.Println("Ответ от сервера:", string(msg))
}
