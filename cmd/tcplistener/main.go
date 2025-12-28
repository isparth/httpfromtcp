package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {

	channel := make(chan string, 1)

	buffer := make([]byte, 8)
	line := ""

	go func() {
		defer f.Close()
		defer close(channel)
		for {
			n, err := f.Read(buffer)

			if n > 0 {

				// Convert the bytes we read to a string
				currentChunk := string(buffer[:n])

				// Check if there is a newline in THIS chunk
				i := strings.Index(currentChunk, "\n")

				if i != -1 {
					// Add everything before the \n to the line and print it.
					line += currentChunk[:i]
					channel <- line

					line = currentChunk[i+1:]
				} else {
					line += currentChunk
				}
			}

			if err != nil {
				if err == io.EOF {
					if line != "" {
						channel <- line
					}

					return
				}
				fmt.Printf("Unexpected error: %v\n", err)
				return
			}
		}

	}()

	return channel

}

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

		ch := getLinesChannel(conn)

		for line := range ch {
			fmt.Printf("%s \n", line)
		}

	}

}
