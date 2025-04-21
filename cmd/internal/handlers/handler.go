package handlers

import (
	"context"
	"go.uber.org/zap"
	"strings"
)

type UserStorage interface {
	GetUserBalance(ctx context.Context, userID int) (float64, float64, error)
}

type Handler struct {
	baseURL string
	Logger  *zap.Logger
	Mode    string
}

func NewHandler(baseURL string, logger *zap.Logger, mode string) *Handler {
	return &Handler{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		Logger:  logger,
		Mode:    mode,
	}
}
