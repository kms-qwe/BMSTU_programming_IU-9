package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

type StringCountRequest struct {
	Text string `json:"text"`
}

type StringCountResponse struct {
	WordCount int `json:"word_count"`
}

func main() {
	listener, err := net.Listen("tcp", ":8085")
	if err != nil {
		fmt.Println("Error go serv:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Serv alive...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("coonect error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	requestData := make([]byte, 1024)
	n, err := conn.Read(requestData)
	if err != nil {
		fmt.Println("read error:", err)
		return
	}

	var request StringCountRequest
	err = json.Unmarshal(requestData[:n], &request)
	if err != nil {
		fmt.Println("decode error JSON:", err)
		return
	}

	wordCount := len(strings.Fields(request.Text))

	response := StringCountResponse{
		WordCount: wordCount,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println("coder error json:", err)
		return
	}

	_, err = conn.Write(jsonResponse)
	if err != nil {
		fmt.Println("go data error:", err)
		return
	}
}
