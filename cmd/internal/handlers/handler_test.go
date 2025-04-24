package handlers_test

import (
	"didactic-giggle/cmd/internal/handlers"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestNewHandler(t *testing.T) {
	logger := zaptest.NewLogger(t)

	handler := handlers.NewHandler("http://localhost:8080/", logger, "dev")

	assert.Equal(t, "http://localhost:8080", handler.BaseURL())
	assert.Equal(t, logger, handler.Logger)
	assert.Equal(t, "dev", handler.Mode)
}
