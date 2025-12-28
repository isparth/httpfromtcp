package main

import (
	"fmt"
	"net"

	"github.com/isparth/httpfromtcp/internal/request"
)

func main() {

	// 1. Listen on a port (TCP protocol)
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 42069...")

	for {
		// 2. Accept a new connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		output, err := request.RequestFromReader(conn)

		fmt.Println(output.RequestLine)

	}

}
