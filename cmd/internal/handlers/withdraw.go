package handlers

import (
	"context"
	"didactic-giggle/cmd/internal/middleware"
	"didactic-giggle/cmd/internal/util"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
)

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type WithdrawService interface {
	Withdraw(ctx context.Context, userID int, orderNumber string, amount float64) error
	GetWithdrawals(ctx context.Context, userID int) ([]util.Withdrawal, error)
}

// структура хэндлера
type WithdrawHandler struct {
	withdrawService WithdrawService
	Logger          *zap.Logger
}

// конструктор
func NewWithdrawHandler(withdrawService WithdrawService, logger *zap.Logger) *WithdrawHandler {
	return &WithdrawHandler{
		Logger:          logger,
		withdrawService: withdrawService,
	}
}

func (h *WithdrawHandler) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Валидация Луна
	if !util.IsValidLuhn(req.Order) {
		http.Error(w, "invalid order number", http.StatusUnprocessableEntity)
		return
	}

	if err := h.withdrawService.Withdraw(r.Context(), userID, req.Order, req.Sum); err != nil {
		if errors.Is(err, util.ErrInsufficientFunds) {
			http.Error(w, "insufficient funds", http.StatusPaymentRequired)
			return
		}
		h.Logger.Error("ошибка при списании", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *WithdrawHandler) GetWithdrawalsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	withdrawals, err := h.withdrawService.GetWithdrawals(r.Context(), userID)
	if err != nil {
		h.Logger.Error("ошибка при получении списка списаний", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(withdrawals)
}
