package request

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/isparth/httpfromtcp/internal/headers"
)

var (
	ErrMalformedRequest  = errors.New("malformed request line")
	ErrUnsupportedMethod = errors.New("invalid or non-uppercase method")
	ErrInvalidTarget     = errors.New("invalid request target path")
	ErrProtocolVersion   = errors.New("unsupported protocol version")
)

var (
	targetRegex  = regexp.MustCompile(`^/`)
	versionRegex = regexp.MustCompile(`^HTTP/(\d+\.\d+)$`)
)

// ParserState is our custom "enum" type
type ParserState int

const (
	// Initialized will be 0
	Initialized ParserState = iota
	// ParsingHeaders will be 1
	ParsingHeaders
	// Done will be 2
	Done
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	state       ParserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (rl RequestLine) String() string {
	return fmt.Sprintf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s",
		rl.Method,
		rl.RequestTarget,
		rl.HttpVersion,
	)
}

func (r *Request) parse(data []byte) (int, error) {

	if r.state == Initialized {
		requestLine, err, consumed := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if requestLine != nil {
			r.RequestLine = *requestLine
			r.state = ParsingHeaders
			return consumed, nil
		}

	}

	if r.state == ParsingHeaders {
		if r.Headers == nil {
			r.Headers = make(headers.Headers)
		}
		consumed, done, err := r.Headers.Parse(data)
		if err != nil {
			return consumed, err
		}
		if done {
			r.state = Done
		}
		return consumed, nil

	}

	return 0, nil
}

func RequestFromReader(r io.Reader) (*Request, error) {
	output := &Request{state: Initialized}
	// This buffer accumulates data across multiple Read calls
	var accumulated []byte
	// Temporary buffer for the current Read
	readBuf := make([]byte, 1024)

	for {
		n, err := r.Read(readBuf)

		if n > 0 {
			// Read
			accumulated = append(accumulated, readBuf[:n]...)

			// Parse
			consumed, parseErr := output.parse(accumulated)
			if parseErr != nil {
				return nil, parseErr
			}

			// Remove the bytes we successfully parsed from our accumulator
			if consumed > 0 {
				accumulated = accumulated[consumed:]
			}
		}

		if output.state == Done {
			break
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}

	return output, nil
}

func parseRequestLine(line string) (*RequestLine, error, int) {
	idx := strings.Index(line, "\r\n")
	if idx == -1 {
		// Not enough data yet, return 0 consumed and no error
		return nil, nil, 0
	}

	totalConsumed := idx + 2
	rawLine := line[:idx]

	parts := strings.Split(rawLine, " ")
	if len(parts) != 3 {
		return nil, ErrMalformedRequest, 0
	}

	method, target, proto := parts[0], parts[1], parts[2]

	if method != strings.ToUpper(method) {
		return nil, ErrUnsupportedMethod, 0
	}
	if !targetRegex.MatchString(target) {
		return nil, ErrInvalidTarget, 0
	}

	versionMatch := versionRegex.FindStringSubmatch(proto)
	if len(versionMatch) != 2 {
		return nil, ErrProtocolVersion, 0
	}

	version := versionMatch[1]
	if version != "1.1" {
		return nil, fmt.Errorf("%w: expected 1.1, got %s", ErrProtocolVersion, version), 0
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   version,
	}, nil, totalConsumed
}
