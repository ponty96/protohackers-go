package main

import (
	"fmt"
	"io"
	"net"

	"go.uber.org/zap"
)

func main() {
	// Start listening on a TCP port
	listener, err := net.Listen("tcp", ":8080")
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
	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())
	defer conn.Close()

	// Create a reader to read data from the connection
	// reader := bufio.NewReader(conn)
	packet := make([]byte, 4096)
	tmp := make([]byte, 4096)

	for {
		// Read data from the connection
		_, err := conn.Read(tmp)

		if err != nil {
			if err != io.EOF {
				zap.L().Error("Read error", zap.Error(err))
			}
			zap.L().Info("End of file")
			break
		}

		zap.L().Info("Data received", zap.String("data", string(tmp[:])))

		// Append the data to the packet
		packet = append(packet, tmp...)
	}

	_, err := conn.Write(packet)
	if err != nil {
		zap.L().Error("Error writing to connection", zap.Error(err))
		return
	}
}
