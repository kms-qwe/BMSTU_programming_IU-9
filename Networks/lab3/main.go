package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

var (
	ADD      = "ADD"
	REMOVE   = "REMOVE"
	FIND     = "FIND"
	NodeIP   = "127.0.0.1"
	NodePort = "8091"
	NextNode = "127.0.0.1:8092"
	node     = NewNode(NodeIP, NodePort, NextNode)
)

type Node struct {
	IP       string
	Port     string
	NextNode string
	Data     map[string]string
	mux      sync.Mutex
}

type Message struct {
	IP        string `json:"IP"`
	Operation string `json:"operation"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

func NewNode(ip, port, nextNode string) *Node {
	return &Node{
		IP:       ip,
		Port:     port,
		NextNode: nextNode,
		Data:     make(map[string]string),
	}
}

func (node *Node) AddKeyValue(message Message) {
	node.mux.Lock()
	node.Data[message.Key] = message.Value
	node.mux.Unlock()

	node.ForwardMessage(message)
	fmt.Printf("Key %s with value %s added: %#v\n", message.Key, message.Value, node.Data)
	sendToWebSocket(fmt.Sprintf("Key %s with value %s added: %#v\n", message.Key, message.Value, node.Data))
}

func (node *Node) RemoveKey(message Message) {
	node.mux.Lock()
	delete(node.Data, message.Key)
	node.mux.Unlock()

	node.ForwardMessage(message)
	fmt.Printf("Key %s removed: %#v\n", message.Key, node.Data)
	sendToWebSocket(fmt.Sprintf("Key %s removed: %#v\n", message.Key, node.Data))
}

func (node *Node) FindValue(message Message) {
	node.mux.Lock()
	value, ok := node.Data[message.Key]
	node.mux.Unlock()

	if ok {
		fmt.Printf("key = %s: value = %s\n", message.Key, value)
		sendToWebSocket(fmt.Sprintf("key = %s: value = %s\n", message.Key, value))
	} else {
		fmt.Printf("key = %s: no such key\n", message.Key)
		sendToWebSocket(fmt.Sprintf("key = %s: no such key\n", message.Key))
	}

	node.ForwardMessage(message)
}

func (node *Node) ForwardMessage(message Message) {
	conn, err := net.Dial("tcp", node.NextNode)
	if err != nil {
		fmt.Println("Error connecting to peer:", err)
		return
	}
	defer conn.Close()

	data, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("marshal error %v\n", err)
	}
	conn.Write(data)
}

func main() {
	listener, err := net.Listen("tcp", node.IP+":"+node.Port)
	if err != nil {
		fmt.Println("Error starting node:", err)
		return
	}
	defer listener.Close()

	webSocketServ := chi.NewRouter()
	webSocketServ.Get("/webSocket", handleWebSocket)
	go func() {
		fmt.Println("WebSocket server started on port 8093")
		err := http.ListenAndServe("localhost:8093", webSocketServ)
		if err != nil {
			fmt.Println("Error starting WebSocket server:", err)
		}
	}()

	httpServ := chi.NewRouter()
	httpServ.Get("/command", handleCommandForm)
	httpServ.Post("/command", handleCommandForm)
	go func() {
		fmt.Println("HTTP server started on port 8100")
		err := http.ListenAndServe("localhost:8100", httpServ)
		if err != nil {
			fmt.Println("Error starting HTTP server:", err)
		}
	}()

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Println("Enter command:")
			command, _ := reader.ReadString('\n')
			readCommand(command)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go handleConnection(conn, node)
	}
}

func handleConnection(conn net.Conn, node *Node) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Data reading error:", err)
		return
	}
	message := Message{}
	err = json.Unmarshal(buffer[:n], &message)
	if err != nil {
		fmt.Printf("unmarshal error %v\n", err)
	}
	if message.IP == node.IP+":"+node.Port {
		fmt.Println("Operation synchronized")
		return
	}
	switch message.Operation {
	case ADD:
		node.AddKeyValue(message)
	case REMOVE:
		node.RemoveKey(message)
	case FIND:
		node.FindValue(message)
	default:
		fmt.Println("Invalid operation")
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var wsConn *websocket.Conn
var wsMux sync.Mutex

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	var err error
	wsMux.Lock()
	wsConn, err = upgrader.Upgrade(w, r, nil)
	wsMux.Unlock()
	if err != nil {
		fmt.Println("WebSocket connection error:", err)
		return
	}
	defer wsConn.Close()

	for {
		_, _, err := wsConn.ReadMessage()
		if err != nil {
			fmt.Println("Client disconnected")
			wsMux.Lock()
			wsConn = nil
			wsMux.Unlock()
			break
		}
	}
}

func sendToWebSocket(message string) {
	wsMux.Lock()
	defer wsMux.Unlock()
	if wsConn != nil {
		err := wsConn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Println("WebSocket send error:", err)
			wsConn = nil
		}
	}
}

func handleCommandForm(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Command Form</title>
	</head>
	<body>
	<h1>Send a Command</h1>
	<form method="POST" action="/command">
	<label for="command">Command:</label>
	<input type="text" id="command" name="command">
	<button type="submit">Submit</button>
	</form>
	</body>
	</html>`
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Form parsing error", http.StatusBadRequest)
			return
		}
		command := r.FormValue("command")
		fmt.Printf("Received command: %s\n", command)
		readCommand(command)
		w.Write([]byte(html))
	}
}

func readCommand(command string) {
	command = strings.TrimSpace(command)
	coms := strings.Fields(command)
	if len(coms) == 0 {
		fmt.Println("Empty command")
		return
	}
	message := Message{IP: NodeIP + ":" + NodePort, Operation: coms[0]}
	switch coms[0] {
	case ADD:
		if len(coms) == 3 {
			fmt.Println("Adding key and value")
			message.Key = coms[1]
			message.Value = coms[2]
			node.AddKeyValue(message)
		} else {
			fmt.Println("Invalid arguments for ADD. Use: ADD key value")
		}
	case REMOVE:
		if len(coms) == 2 {
			message.Key = coms[1]
			node.RemoveKey(message)
		} else {
			fmt.Println("Invalid arguments for REMOVE. Use: REMOVE key")
		}
	case FIND:
		if len(coms) == 2 {
			message.Key = coms[1]
			node.FindValue(message)
		} else {
			fmt.Println("Invalid arguments for FIND. Use: FIND key")
		}
	default:
		fmt.Println("Invalid operation")
	}
}
