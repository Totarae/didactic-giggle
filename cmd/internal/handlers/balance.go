package handlers

import (
	"didactic-giggle/cmd/internal/middleware"
	"didactic-giggle/cmd/internal/util"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
)

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func (h *Handler) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	current, withdrawn, err := h.store.GetUserBalance(r.Context(), userID)
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

func (h *Handler) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := h.store.Withdraw(r.Context(), userID, req.Order, req.Sum); err != nil {
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

func (h *Handler) GetWithdrawalsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	withdrawals, err := h.store.GetWithdrawals(r.Context(), userID)
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
