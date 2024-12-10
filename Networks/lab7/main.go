package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

// структура для хранения данных о блоке
type BlockInfo struct {
	Number       uint64 `json:"number"`
	Time         uint64 `json:"time"`
	Difficulty   uint64 `json:"difficulty"`
	Hash         string `json:"hash"`
	Transactions int    `json:"transactions"`
}

// структура для хранения данных о транзакции
type TransactionInfo struct {
	Hash     string `json:"hash"`
	Value    string `json:"value"`
	Gas      uint64 `json:"gas"`
	GasPrice string `json:"gas_price"`
	To       string `json:"to"`
}

func main() {
	// установление соединенияи для получения данных о криптовалюте
	infuraURL := "https://polygon-mainnet.infura.io/v3/a8281f70daf446d88a7a5ad46ed1707f"
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to Infura: %v", err)
	}

	// цикл с обработкой информации
	for {
		err := monitorBlockchain(client)
		if err != nil {
			log.Printf("Error in monitoring: %v", err)
		}
		time.Sleep(1 * time.Second)
	}
}

func monitorBlockchain(client *ethclient.Client) error {
	ctx := context.Background()

	// получение информации
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to get latest block header: %w", err)
	}

	blockNumber := big.NewInt(header.Number.Int64())
	block, err := client.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return fmt.Errorf("failed to get block: %w", err)
	}

	blockInfo := BlockInfo{
		Number:       block.Number().Uint64(),
		Time:         block.Time(),
		Difficulty:   block.Difficulty().Uint64(),
		Hash:         block.Hash().Hex(),
		Transactions: len(block.Transactions()),
	}

	// подключение к бд
	databaseURL := "https://lab7-d706c-default-rtdb.europe-west1.firebasedatabase.app/"

	blockPath := fmt.Sprintf("blocks/%d.json", blockInfo.Number)
	err = sendToFirebase(databaseURL+blockPath, blockInfo)
	if err != nil {
		return fmt.Errorf("failed to save block info to Firebase: %w", err)
	}
	log.Printf("Block %d saved to Firebase", blockInfo.Number)

	// цикл для сохранения данных о транзакциях в блоке
	for _, tx := range block.Transactions() {
		txInfo := TransactionInfo{
			Hash:     tx.Hash().Hex(),
			Value:    tx.Value().String(),
			Gas:      tx.Gas(),
			GasPrice: tx.GasPrice().String(),
		}

		if tx.To() != nil {
			txInfo.To = tx.To().Hex()
		} else {
			txInfo.To = "Contract Creation"
		}

		txPath := fmt.Sprintf("blocks/%d/transactions/%s.json", blockInfo.Number, txInfo.Hash)
		err := sendToFirebase(databaseURL+txPath, txInfo)
		if err != nil {
			log.Printf("Failed to save transaction %s: %v", txInfo.Hash, err)
		}
	}

	return nil
}

// функция для отправки запроса к бд
func sendToFirebase(url string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request to Firebase: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to save data to Firebase, status: %s", resp.Status)
	}

	return nil
}
