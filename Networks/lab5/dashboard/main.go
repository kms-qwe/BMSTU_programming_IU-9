package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan []TableRow)
	db        *sql.DB
)

type TableRow struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка при подключении:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true

	for {
		select {
		case msg := <-broadcast:
			jsonData, err := json.Marshal(msg)
			if err != nil {
				log.Println("Ошибка при маршаллинге данных:", err)
				continue
			}
			for client := range clients {
				if err := client.WriteMessage(websocket.TextMessage, jsonData); err != nil {
					log.Println("Ошибка при отправке сообщения:", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}

func generateTableStatus() ([]TableRow, error) {
	rows, err := db.Query("SELECT id, title, body FROM iu9Krasnobaev")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tableRows []TableRow
	for rows.Next() {
		var row TableRow
		if err := rows.Scan(&row.ID, &row.Title, &row.Body); err != nil {
			return nil, err
		}
		tableRows = append(tableRows, row)
	}
	return tableRows, nil
}

func handleUpdates() {
	for {
		time.Sleep(10 * time.Second)
		status, err := generateTableStatus()
		if err != nil {
			log.Printf("Ошибка генерации статуса таблицы: %v", err)
			continue
		}
		broadcast <- status
	}
}

func main() {
	var err error
	db, err = sql.Open("mysql", "iu9networkslabs:Je2dTYr6@tcp(students.yss.su:3306)/iu9networkslabs")
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	go handleUpdates()

	http.HandleFunc("/ws", handleConnections)

	fmt.Println("Сервер запущен на :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
