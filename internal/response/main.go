package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/isparth/httpfromtcp/internal/headers"
)

// StatusCode represents an HTTP status code
type StatusCode int

// Define the constants for the status codes we care about
const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	s := ""
	switch statusCode {
	case 200:
		s = "HTTP/1.1 200 OK"
	case 400:
		s = "HTTP/1.1 400 Bad Request"
	case 500:
		s = "HTTP/1.1 500 Internal Server Error"
	}

	line := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, s)
	_, err := w.Write([]byte(line))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {

	h := headers.Headers{}
	h["Content-Length"] = strconv.Itoa(contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, h headers.Headers) error {
	for key, value := range h {
		line := fmt.Sprintf("%s: %s\r\n", key, value)
		if _, err := w.Write([]byte(line)); err != nil {
			return err
		}
	}
	// The CRLF that separates headers from the body
	_, err := w.Write([]byte("\r\n"))
	return err
}
