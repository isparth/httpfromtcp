package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/isparth/httpfromtcp/internal/request"
	"github.com/isparth/httpfromtcp/internal/response"
	"github.com/isparth/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	// 1. Define our handler logic
	handler := func(w *response.Writer, req *request.Request) {
		target := req.RequestLine.RequestTarget

		fmt.Println(target)

		if strings.HasPrefix(target, "/httpbin/") {
			path := strings.TrimPrefix(target, "/httpbin")
			resp, err := http.Get("https://httpbin.org" + path)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			headers := response.GetDefaultHeaders(0)
			delete(headers, "content-length")
			headers.Set("Transfer-Encoding", "chunked")
			if contentType := resp.Header.Get("Content-Type"); contentType != "" {
				headers.Set("Content-Type", contentType)
			}

			if err := w.WriteStatusLine(response.StatusOK); err != nil {
				return
			}
			if err := w.WriteHeaders(headers); err != nil {
				return
			}

			buf := make([]byte, 1024)
			for {
				n, readErr := resp.Body.Read(buf)
				if n > 0 {
					fmt.Println(n)
					if _, err := w.WriteChunkedBody(buf[:n]); err != nil {
						return
					}
				}
				if readErr != nil {
					if readErr == io.EOF {
						break
					}
					return
				}
			}

			_, _ = w.WriteChunkedBodyDone()
			return
		}

		var (
			status response.StatusCode
			body   string
		)

		// Routing logic
		switch target {
		case "/yourproblem":
			status = response.StatusBadRequest
			body = `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
		case "/myproblem":
			status = response.StatusInternalServerError
			body = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
		default:
			status = response.StatusOK
			body = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
		}

		bodyBytes := []byte(body)
		headers := response.GetDefaultHeaders(len(bodyBytes))
		headers.Set("Content-Type", "text/html")

		if err := w.WriteStatusLine(status); err != nil {
			return
		}
		if err := w.WriteHeaders(headers); err != nil {
			return
		}
		_, _ = w.WriteBody(bodyBytes)
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
