package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"go.uber.org/zap"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Start listening on a TCP port
	listener, err := net.Listen("tcp", ":"+port)
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))

	if err != nil {
		zap.L().Fatal("Error starting TCP server", zap.Error(err))
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 8080")

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			zap.L().Error("Error accepting connection", zap.Error(err))
			continue
		}
		fmt.Println("Connection established")

		// Handle the connection in a new goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Printf("Serving at this address %s\n", conn.RemoteAddr().String())
	defer conn.Close()

	_, err := io.Copy(conn, conn)
	if err != nil {
		zap.L().Error("Error writing to connection", zap.Error(err))
		return
	}
}
