package headers

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

var (
	ErrMalformedHeader = errors.New("malformed Header")
)

var headerRegex = regexp.MustCompile(`^[ \t]*([a-zA-Z0-9!#$%&'*+\-.^_` + "`" + `|~]+):[ \t]*(.*?)[ \t]*$`)

func (h *Headers) Parse(data []byte) (int, bool, error) {

	// 1. Find the next CRLF
	lineEnd := bytes.Index(data, []byte("\r\n"))

	// If no CRLF is found, we don't have a full line yet
	if lineEnd == -1 {
		return 0, false, nil
	}

	if lineEnd == 0 {

		return 2, true, nil
	}

	// Calculate actual position in the original data slice

	line := data[:lineEnd]

	// 3. Process the header line
	key, value, err := parseSingleHeader(string(line))
	if err != nil {
		return 0, false, err
	}

	existingValue, exists := (*h)[key]

	if exists {
		(*h)[key] = existingValue + ", " + value
	} else {
		(*h)[key] = value
	}

	return lineEnd + 2, false, nil

}

func (h *Headers) Get(key string) string {
	return (*h)[strings.ToLower(key)]
}

func (h Headers) String() string {
	if len(h) == 0 {
		return "Headers:\n- (none)"
	}

	var b strings.Builder
	b.WriteString("Headers:\n")

	for k, v := range h {
		fmt.Fprintf(&b, "- %s: %s\n", k, v)
	}

	return b.String()
}

func parseSingleHeader(line string) (string, string, error) {
	rawLine := strings.TrimSpace(line)
	if !headerRegex.MatchString(rawLine) {
		fmt.Println(rawLine)
		return "", "", fmt.Errorf("%v:, got %s", ErrMalformedHeader, rawLine)
	}

	parts := strings.SplitN(rawLine, ":", 2)

	return strings.ToLower(parts[0]), strings.TrimSpace(parts[1]), nil

}
