package main

import (
	"log"
	"net/http"
)

func main() {
	router := NewRouter()
	log.Println("API Gateway running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
