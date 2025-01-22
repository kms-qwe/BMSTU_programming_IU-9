package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type RandomNumberResponse struct {
	Number int `json:"number"`
}

func randomHandler(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(100)

	response := RandomNumberResponse{Number: randomNumber}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/random", randomHandler)
	port := "8080"
	println("serv is strating " + port)
	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		panic(err)
	}
}
