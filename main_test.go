package main

import (
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

// create a mock HTTP client to intercept the traffic
type MockClient struct {
	Requests []*http.Request
}

func (client *MockClient) Do(req *http.Request) (*http.Response, error) {
	client.Requests = append(client.Requests, req)
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("SUCCESS")),
	}, nil
}

func Test_NewServer(t *testing.T) {
	client := new(MockClient)
	server := NewServer("9001", "localhost:9002", client)
	go func() {
		server.Start()
	}()

	time.Sleep(250 * time.Millisecond)

	// Test the server as a client
	conn, err := net.Dial("tcp", ":9001")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}

	conn.Write([]byte("ConfrenceDualMode"))
	time.Sleep(750 * time.Millisecond)        // allow the server time to work
	conn.Write([]byte("ConfrenceSingleMode")) // send more data
	time.Sleep(750 * time.Millisecond)        // allow the server time to work

	conn.Close()

	if len(client.Requests) == 0 {
		t.Fatal("server did not make http request")
	}

	if len(client.Requests) != 2 {
		t.Fatal("server did not make both http requests")
	}

	if client.Requests[0].URL.Host != "localhost:9002" {
		t.Errorf("expected %s got %s", "localhost:9002", client.Requests[0].URL.Host)
	}
}
