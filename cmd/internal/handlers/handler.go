package handlers

import (
	"go.uber.org/zap"
	"strings"
)

type AppHandlers struct {
	User        *UserHandler
	OrderWriter *OrderWriterHandler
	OrderReader *OrderReaderHandler
	Balance     *BalanceHandler
	Withdraw    *WithdrawHandler
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
