package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/isparth/httpfromtcp/internal/request"
	"github.com/isparth/httpfromtcp/internal/response"
	"github.com/isparth/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	// 1. Define our handler logic
	handler := func(w io.Writer, req *request.Request) *server.HandlerError {
		target := req.RequestLine.RequestTarget

		fmt.Println(target)

		// Routing logic
		switch target {
		case "/yourproblem":
			return &server.HandlerError{
				StatusCode: response.StatusBadRequest,
				Message:    "Your problem is not my problem\n",
			}
		case "/myproblem":
			return &server.HandlerError{
				StatusCode: response.StatusInternalServerError,
				Message:    "Woopsie, my bad\n",
			}
		default:
			// Success case: Write to the buffer (w) and return nil error
			fmt.Fprint(w, "All good, frfr\n")
			return nil
		}
	}

	// 2. Pass the handler into Serve
	srv, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer srv.Close()

	log.Println("Server started on port", port)

	// Graceful shutdown logic
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
