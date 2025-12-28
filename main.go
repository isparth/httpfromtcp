package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("Could not open file: %v\n", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, 8)
	line := ""

	for {
		n, err := file.Read(buffer)

		if n > 0 {

			// Convert the bytes we read to a string
			currentChunk := string(buffer[:n])

			// Check if there is a newline in THIS chunk
			i := strings.Index(currentChunk, "\n")

			if i != -1 {
				// Add everything before the \n to the line and print it.
				line += currentChunk[:i]
				fmt.Printf("Finished line: %s\n", line)

				line = currentChunk[i+1:]
			} else {
				line += currentChunk
			}
		}

		if err != nil {
			if err == io.EOF {
				if line != "" {
					fmt.Printf("Final line: %s\n", line)
				}
				fmt.Println("--- End of File reached ---")
				break
			}
			fmt.Printf("Unexpected error: %v\n", err)
			break
		}
	}
}
