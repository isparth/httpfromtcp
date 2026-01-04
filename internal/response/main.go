package response

import (
	"errors"
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

type Headers = headers.Headers

type writerState int

const (
	writerStateStart writerState = iota
	writerStateStatusWritten
	writerStateHeadersWritten
	writerStateBodyWritten
	writerStateChunked
	writerStateDone
)

var ErrInvalidWriterState = errors.New("response writer called out of order")

type Writer struct {
	w     io.Writer
	state writerState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w, state: writerStateStart}
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	line := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusText(statusCode))
	_, err := w.Write([]byte(line))
	return err
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != writerStateStart {
		return ErrInvalidWriterState
	}
	w.state = writerStateStatusWritten
	return WriteStatusLine(w.w, statusCode)
}

func GetDefaultHeaders(contentLen int) Headers {

	h := headers.Headers{}
	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, h Headers) error {
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

func (w *Writer) WriteHeaders(h Headers) error {
	if w.state != writerStateStatusWritten {
		return ErrInvalidWriterState
	}
	w.state = writerStateHeadersWritten
	return WriteHeaders(w.w, h)
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != writerStateHeadersWritten && w.state != writerStateBodyWritten {
		return 0, ErrInvalidWriterState
	}
	w.state = writerStateBodyWritten
	return w.w.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.state != writerStateHeadersWritten && w.state != writerStateChunked {
		return 0, ErrInvalidWriterState
	}
	w.state = writerStateChunked
	header := fmt.Sprintf("%x\r\n", len(p))
	if _, err := w.w.Write([]byte(header)); err != nil {
		return 0, err
	}
	if _, err := w.w.Write(p); err != nil {
		return 0, err
	}
	if _, err := w.w.Write([]byte("\r\n")); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.state != writerStateHeadersWritten && w.state != writerStateChunked {
		return 0, ErrInvalidWriterState
	}
	w.state = writerStateDone
	return w.w.Write([]byte("0\r\n\r\n"))
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.state != writerStateHeadersWritten && w.state != writerStateChunked {
		return ErrInvalidWriterState
	}
	w.state = writerStateDone
	if _, err := w.w.Write([]byte("0\r\n")); err != nil {
		return err
	}
	for key, value := range h {
		line := fmt.Sprintf("%s: %s\r\n", key, value)
		if _, err := w.w.Write([]byte(line)); err != nil {
			return err
		}
	}
	_, err := w.w.Write([]byte("\r\n"))
	return err
}

func statusText(statusCode StatusCode) string {
	switch statusCode {
	case StatusOK:
		return "OK"
	case StatusBadRequest:
		return "Bad Request"
	case StatusInternalServerError:
		return "Internal Server Error"
	default:
		return ""
	}
}
