package request

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
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

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(r io.Reader) (*Request, error) {
	// Use a scanner or bufio.Reader to get only the first line

	reader := bufio.NewReader(r)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	// 2. Clean up the line (remove \r\n)
	line = strings.TrimSpace(line)
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, ErrMalformedRequest
	}

	// 3. Validate Method (Must be all caps)
	if parts[0] != strings.ToUpper(parts[0]) {
		return nil, ErrUnsupportedMethod
	}

	// 4. Validate Request Target (Must start with /)
	if !targetRegex.MatchString(parts[1]) {
		return nil, ErrInvalidTarget
	}

	// 5. Validate Protocol and Extract Version
	versionMatch := versionRegex.FindStringSubmatch(parts[2])
	if len(versionMatch) != 2 {
		return nil, ErrProtocolVersion
	}

	version := versionMatch[1]
	if version != "1.1" {
		return nil, fmt.Errorf("%w: expected 1.1, got %s", ErrProtocolVersion, version)
	}

	// 6. Return the Request
	return &Request{
		RequestLine: RequestLine{
			Method:        parts[0],
			RequestTarget: parts[1],
			HttpVersion:   version,
		},
	}, nil
}
