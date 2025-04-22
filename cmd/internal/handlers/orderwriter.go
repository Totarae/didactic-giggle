package handlers

import (
	"context"
	"didactic-giggle/cmd/internal/middleware"
	"didactic-giggle/cmd/internal/util"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type OrderWriter interface {
	SaveOrder(ctx context.Context, userID int, order string) (util.OrderSaveStatus, error)
	UpdateOrder(ctx context.Context, number string, status string, accrual float64) error
}
type OrderWriterHandler struct {
	store  OrderWriter
	Logger *zap.Logger
}

func NewOrderWriterHandler(store OrderWriter, logger *zap.Logger) *OrderWriterHandler {
	return &OrderWriterHandler{
		store:  store,
		Logger: logger,
	}
}

func (h *OrderWriterHandler) UploadOrderHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	order := strings.TrimSpace(string(body))

	if _, err := strconv.ParseInt(order, 10, 64); err != nil {
		http.Error(w, "invalid order format", http.StatusUnprocessableEntity)
		return
	}

	if !util.IsValidLuhn(order) {
		http.Error(w, "invalid order number (Luhn check failed)", http.StatusUnprocessableEntity)
		return
	}

	status, err := h.store.SaveOrder(r.Context(), userID, order)
	if err != nil {
		h.Logger.Error("failed to save order", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	switch status {
	case util.OrderSavedNew:
		w.WriteHeader(http.StatusAccepted)
	case util.OrderAlreadyUploadedByUser:
		w.WriteHeader(http.StatusOK)
	case util.OrderUploadedByAnotherUser:
		http.Error(w, "order already uploaded by another user", http.StatusConflict)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
