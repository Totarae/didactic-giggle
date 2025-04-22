package handlers

import (
	"context"
	"didactic-giggle/cmd/internal/middleware"
	"didactic-giggle/cmd/internal/util"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type OrderReader interface {
	GetOrders(ctx context.Context, userID int) ([]util.OrderInfo, error)
	GetPendingOrders(ctx context.Context) ([]string, error)
}

type OrderReaderHandler struct {
	store  OrderReader
	Logger *zap.Logger
}

func NewOrderReaderHandler(store OrderReader, logger *zap.Logger) *OrderReaderHandler {
	return &OrderReaderHandler{
		store:  store,
		Logger: logger,
	}
}

func (h *OrderReaderHandler) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := h.store.GetOrders(r.Context(), userID)
	if err != nil {
		h.Logger.Error("failed to get orders", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(orders)
}
