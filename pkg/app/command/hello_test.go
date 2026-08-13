package command

import (
	"bytes"
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestHelloExecute(t *testing.T) {
	var buf bytes.Buffer
	logger := zerolog.New(&buf)

	cmd := NewHello(&logger)
	exitCode := cmd.Execute(context.Background(), []string{"custom message"})

	assert.Equal(t, 0, exitCode)
	assert.Contains(t, buf.String(), "sample command executed")
	assert.Contains(t, buf.String(), "custom message")
}
