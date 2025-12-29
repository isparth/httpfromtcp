package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersLineParse(t *testing.T) {
	// Test: Valid single header
	headers := Headers{}
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = Headers{}
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	headers = Headers{}
	data = []byte("       H©st: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestHeadersMultipleValues(t *testing.T) {
	headers := Headers{}

	// First occurrence
	data1 := []byte("Set-Person: lane-loves-go\r\n")
	n1, done1, err1 := headers.Parse(data1)
	require.NoError(t, err1)
	assert.Equal(t, len(data1), n1)
	assert.False(t, done1)

	// Second occurrence (different casing to test case-insensitivity)
	data2 := []byte("set-person: prime-loves-zig\r\n")
	n2, done2, err2 := headers.Parse(data2)
	require.NoError(t, err2)
	assert.Equal(t, len(data2), n2)
	assert.False(t, done2)

	// Verify the combined result
	// Note: The key in the map should be lowercase ("set-person")
	expectedValue := "lane-loves-go, prime-loves-zig"
	assert.Equal(t, expectedValue, headers["set-person"])
}
