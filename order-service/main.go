package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Order struct {
	ID     int     `json:"id"`
	Amount float64 `json:"amount"`
}

func main() {
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		orders := []Order{
			{ID: 1, Amount: 100.50},
			{ID: 2, Amount: 200.75},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	})

	log.Println("Order Service running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
