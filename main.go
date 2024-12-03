package main

import (
	"net/http"
	"os"
	"time"
)

func main() {
	server := NewServer(
		os.Getenv("PORT"),
		2*time.Second,
		os.Getenv("TARGET"),
		http.DefaultClient,
	)
	server.Start()
}
