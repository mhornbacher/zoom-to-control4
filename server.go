package main

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

// Represents the TCP server that forwards requests
type Server struct {
	port   string
	target string
	client Client
}

// Returns a new Server
func NewServer(
	port string,
	target string,
	client Client,
) *Server {
	return &Server{
		port:   port,
		target: target,
		client: client,
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

	// create a logger with the IP for tracing concurrent requests
	defaultLogInfo := []slog.Attr{
		slog.String("remote address", conn.RemoteAddr().String()),
	}
	log := slog.New(slog.Default().Handler().WithAttrs(defaultLogInfo))

	log.Info("received connection")

	buffer := make([]byte, 4096) // we don't expect any messages to be above 4kb

	for {
		n, err := conn.Read(buffer)
		// break the loop if the connection is closed
		if err != nil && err != io.EOF {
			log.Error("connection error, closing connection", "error", err)
			break
		} else if err == io.EOF {
			message := string(buffer[:n])
			log.Warn("EOF encountered, closing connection", "message", message)

			// if we got a message before it closed, honor the request
			if len(message) > 0 {
				go server.sendRequest(log, message)
			}
			break
			// only send data if we have any
		} else if n > 0 {
			message := string(buffer[:n])
			go server.sendRequest(log, message)
		}
	}

	log.Info("closed connection")
}

// sends a request to the target with the message in the URL
func (server *Server) sendRequest(log *slog.Logger, message string) {
	url := fmt.Sprintf("http://%s/%s", server.target, strings.TrimSpace(message))
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

	log.Info("received response",
		slog.Int("status code", response.StatusCode),
		slog.String("body", string(body)))
}
