package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type StringCountRequest struct {
	Text string `json:"text"`
}

type StringCountResponse struct {
	WordCount int `json:"word_count"`
}

func main() {
	conn, err := net.Dial("tcp", "185.102.139.169:8085")
	if err != nil {
		fmt.Println("error connect to serc:", err)
		return
	}
	defer conn.Close()

	fmt.Print("Input: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	request := StringCountRequest{
		Text: text,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("marhsal JSON error:", err)
		return
	}

	_, err = conn.Write(jsonData)
	if err != nil {
		fmt.Println("data trans error:", err)
		return
	}

	responseData := make([]byte, 1024)
	n, err := conn.Read(responseData)
	if err != nil {
		fmt.Println("get ans eeror:", err)
		return
	}

	var response StringCountResponse
	err = json.Unmarshal(responseData[:n], &response)
	if err != nil {
		fmt.Println("unmarshal error JSON:", err)
		return
	}

	fmt.Printf("ANS: %d\n", response.WordCount)
}
