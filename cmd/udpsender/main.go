package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	address := "localhost:42069"

	// Resolve the string address into a *net.UDPAddr
	udpAddr, err := net.ResolveUDPAddr("udp", address)

	if err != nil {
		fmt.Printf("Resolution failed: %s\n", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error:", err)
		}

		n, err := conn.Write([]byte(line))

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf("Written %v Bytes", n)

	}

}
