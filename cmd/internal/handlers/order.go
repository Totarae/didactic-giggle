package handlers

import (
	"context"
	"didactic-giggle/cmd/internal/middleware"
	"didactic-giggle/cmd/internal/util"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type OrderService interface {
	SaveOrder(ctx context.Context, userID int, order string) (util.OrderSaveStatus, error)
	GetOrders(ctx context.Context, userID int) ([]util.OrderInfo, error)
	GetPendingOrders(ctx context.Context) ([]string, error)
	UpdateOrder(ctx context.Context, number string, status string, accrual float64) error
}

// структура хэндлера
type OrderHandler struct {
	store  OrderService
	Logger *zap.Logger
}

// конструктор
func NewOrderHandler(storage OrderService, logger *zap.Logger) *OrderHandler {
	return &OrderHandler{
		Logger: logger,
		store:  storage,
	}
}

func (h *OrderHandler) UploadOrderHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	order := strings.TrimSpace(string(body))

	// Проверяем: только цифры
	if _, err := strconv.ParseInt(order, 10, 64); err != nil {
		http.Error(w, "invalid order format", http.StatusUnprocessableEntity)
		return
	}

	// Проверка по алгоритму Луна
	if !util.IsValidLuhn(order) {
		http.Error(w, "invalid order number (Luhn check failed)", http.StatusUnprocessableEntity)
		return
	}

	// Сохраняем заказ
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

func (h *OrderHandler) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
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
