package handlers

import (
	"context"
	"didactic-giggle/cmd/internal/util"
	"go.uber.org/zap"
	"strings"
)

type UserStorage interface {
	CreateUser(ctx context.Context, login, password string) error
	AuthenticateUser(ctx context.Context, login, password string) (bool, error)
	SaveOrder(ctx context.Context, userID int, order string) (util.OrderSaveStatus, error)
	GetUserIDByLogin(ctx context.Context, login string) (int, error)
	GetOrders(ctx context.Context, userID int) ([]util.OrderInfo, error)
}

type Handler struct {
	store   UserStorage // Use the new URLStore for thread safety
	baseURL string
	//Repo    repositories.URLRepositoryInterface
	Logger *zap.Logger
	Mode   string
	//Auth    *auth.Auth
}

func NewHandler(baseURL string, logger *zap.Logger, mode string, storage UserStorage) *Handler {
	return &Handler{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		Logger:  logger,
		Mode:    mode,
		store:   storage,
	}
}
