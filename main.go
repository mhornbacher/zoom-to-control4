package main

import (
	"net/http"
	"os"
)

func main() {
	server := NewServer(
		os.Getenv("PORT"),
		os.Getenv("TARGET"),
		http.DefaultClient,
	)
	server.Start()
}
