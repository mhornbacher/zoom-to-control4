package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

// Represents the TCP server that forwards requests
type Server struct {
	port    string
	timeout time.Duration
	target  string
	client  Client
}

// Returns a new Server
func NewServer(
	port string,
	timeout time.Duration,
	target string,
	client Client,
) *Server {
	return &Server{
		port:    port,
		timeout: timeout,
		target:  target,
		client:  client,
	}
}

// Starts the server
func (server *Server) Start() {
	listener, err := net.Listen("tcp", ":"+server.port)
	if err != nil {
		slog.Error("unable to start server", slog.Any("error", err))
		return
	}
	defer listener.Close()

	slog.Info("started server", "port", server.port)

	for {
		c, err := listener.Accept()
		if err != nil {
			slog.Error("unable to accept request", slog.Any("error", err))
			continue
		}
		go server.handleConnection(c)
	}
}

func (server *Server) handleConnection(conn net.Conn) {
	defer conn.Close() // always close the connection
	conn.SetReadDeadline(time.Now().Add(server.timeout))

	// create a logger with the IP for tracing concurrent requests
	defaultLogInfo := []slog.Attr{
		slog.String("remote address", conn.RemoteAddr().String()),
	}
	log := slog.New(slog.Default().Handler().WithAttrs(defaultLogInfo))

	log.Info("received request")

	message, err := bufio.NewReader(conn).ReadString('\n')
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		log.Error("connection timed out")
		return
	} else if err == io.EOF {
		log.Info("finised reading data", "message", message)
	} else if err != nil {
		log.Error("error reading data", "error", err)
		return
	}

	url := fmt.Sprintf("http://%s/%s",
		server.target, strings.TrimSpace(string(message)))
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("unable to create request", "error", err)
		return
	}

	log.Info("sending request", "uri", request.URL)

	response, err := server.client.Do(request)
	if err != nil {
		log.Error("request failed", "error", err)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error("unable to read body", "error", err)
	}

	log.Info("response",
		slog.Int("status code", response.StatusCode),
		slog.String("body", string(body)))
}
