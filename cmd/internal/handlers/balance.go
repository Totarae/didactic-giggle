package handlers

import (
	"context"
	"didactic-giggle/cmd/internal/middleware"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type BalanceService interface {
	GetUserBalance(ctx context.Context, userID int) (float64, float64, error)
}

// структура хэндлера
type BalanceHandler struct {
	balanceService BalanceService
	Logger         *zap.Logger
}

// конструктор
func NewBalanceHandler(balanceService BalanceService, logger *zap.Logger) *BalanceHandler {
	return &BalanceHandler{
		Logger:         logger,
		balanceService: balanceService,
	}
}

func (h *BalanceHandler) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	current, withdrawn, err := h.balanceService.GetUserBalance(r.Context(), userID)
	if err != nil {
		h.Logger.Error("ошибка получения баланса", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := BalanceResponse{
		Current:   current,
		Withdrawn: withdrawn,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
