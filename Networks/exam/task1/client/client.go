package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RandomNumberResponse struct {
	Number int `json:"number"`
}

func main() {
	url := "http://185.102.139.161:8080/random"

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d\n", resp.StatusCode)
		return
	}

	var response RandomNumberResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		panic(err)
	}

	fmt.Printf("Received random number: %d\n", response.Number)
}
