package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHowManyBytesLeft(t *testing.T) {
	buf := []byte{'a', 'p', 'p', 'l', 'e'}

	in := bytes.NewReader(buf)

	input := make([]byte, 3)

	n, err := in.Read(input)

	assert.Equal(t, nil, err)
	assert.Equal(t, 3, n)

    assert.Equal(t, 2, in.Len())

    n, err = in.Read(input)

	assert.Equal(t, nil, err)
	assert.Equal(t, 2, n)

    n, err = in.Read(input)

	assert.Equal(t, io.EOF, err)
}
