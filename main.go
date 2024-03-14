package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kunal768/trustwallet/clients"
	"github.com/kunal768/trustwallet/parser"
)

func initEthClient() clients.EthereumClient {
	client := http.Client{}
	return clients.ClientService(&client)
}

func main() {
	ethClient := initEthClient()
	parserSvc := parser.NewParserService(ethClient)
	ctx := context.Background()

	// API endpoints
	http.HandleFunc("/currentBlock", func(w http.ResponseWriter, r *http.Request) {
		block, err := parserSvc.GetCurrentBlock(ctx)
		if err != nil {
			http.Error(w, "Failed to get current block", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Current block: %s", block)
	})

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		success := parserSvc.Subscribe(ctx, address)
		if success {
			fmt.Fprintf(w, "Subscribed to address: %s", address)
		} else {
			http.Error(w, "Address already subscribed", http.StatusConflict)
		}
	})

	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		txns, err := parserSvc.GetTransactions(ctx, address)
		if err != nil {
			http.Error(w, "Failed to get transactions", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(txns)
	})

	// Start a goroutine to periodically update subscriber transactions
	exit := make(chan struct{}) // Channel to signal goroutine exit
	go func() {
		defer close(exit) // Ensure channel is closed upon completion
		for {
			select {
			case <-ctx.Done(): // Exit if context is cancelled
				return
			case <-exit: // Exit if signal is received
				return
			default:
				parserSvc.UpdateSubscriberTxns(ctx)
				time.Sleep(time.Minute)
			}
		}
	}()

	http.ListenAndServe(":8080", nil)

	// stop the goroutine when the server shuts down
	<-exit
}
