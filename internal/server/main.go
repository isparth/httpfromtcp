package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/isparth/httpfromtcp/internal/request"
	"github.com/isparth/httpfromtcp/internal/response"
)

var (
	ErrMissingListenner = errors.New("Closed a lisner that was nil")
)

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	listener net.Listener
	handler  Handler
	isClosed atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	addr := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", addr)

	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: l,
		handler:  handler,
	}

	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.isClosed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}

	return ErrMissingListenner
}

func (s *Server) listen() {

	for !s.isClosed.Load() {

		conn, err := s.listener.Accept()
		if err != nil {
			// Handle errors (e.g., if the listener is closed)
			log.Printf("Accept error: %v", err)
			continue
		}

		go s.handle(conn)

	}

}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("Error parsing request: %v", err)
		return
	}

	writer := response.NewWriter(conn)
	s.handler(writer, req)
}
