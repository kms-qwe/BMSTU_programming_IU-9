package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func randomNumberHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Ошибка при обновлении до WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// fmt.Println("Ошибка при чтении сообщения:", err)
			break
		}

		if string(msg) == "get_random" {
			rand.Seed(time.Now().UnixNano())
			randomNumber := rand.Intn(100)

			err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Random number: %d", randomNumber)))
			if err != nil {
				fmt.Println("Ошибка при отправке сообщения:", err)
				break
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", randomNumberHandler)

	port := "8081"
	fmt.Println("WebSocket сервер запущен на порту", port)
	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		panic("Ошибка запуска сервера: " + err.Error())
	}
}
