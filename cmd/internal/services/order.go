package services

import (
	"context"
	"didactic-giggle/cmd/internal/database"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"time"
)

type OrderService struct {
	storage          *database.DB
	accrualSystemURL string
	logger           *zap.Logger
}

func NewOrderService(repo *database.DB, url string, logger *zap.Logger) *OrderService {
	return &OrderService{storage: repo, accrualSystemURL: url, logger: logger}
}

// запускаеv фоновый процесс
func (s *OrderService) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				s.logger.Info("Воркер статусов остановлен")
				return
			case <-ticker.C:
				s.syncPendingOrders(ctx)
			}
		}
	}()
}

// syncPendingOrders обходит все ожидающие заказы
func (s *OrderService) syncPendingOrders(ctx context.Context) {
	orders, err := s.storage.GetPendingOrders(ctx)
	if err != nil {
		s.logger.Error("failed to get pending orders", zap.Error(err))
		return
	}

	for _, number := range orders {
		if err := s.fetchAndStoreAccrual(ctx, number); err != nil {
			s.logger.Warn("fetchAccrual failed", zap.String("order", number), zap.Error(err))
		}
	}
}

// fetchAndStoreAccrual обращаемся к черному ящику
func (s *OrderService) fetchAndStoreAccrual(ctx context.Context, number string) error {
	u, err := url.Parse(s.accrualSystemURL)
	if err != nil {
		return fmt.Errorf("невалидный accrualSystemURL: %w", err)
	}
	u.Path = fmt.Sprintf("/api/orders/%s", number)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var result struct {
			Order   string   `json:"order"`
			Status  string   `json:"status"`  // REGISTERED, PROCESSING, INVALID, PROCESSED
			Accrual *float64 `json:"accrual"` // может быть 0
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		// Преобразуем во внутренний статус
		internalStatus := map[string]string{
			"REGISTERED": "PROCESSING",
			"PROCESSING": "PROCESSING",
			"INVALID":    "INVALID",
			"PROCESSED":  "PROCESSED",
		}[result.Status]

		// Перехватим null
		accrualValue := 0.0
		if result.Accrual != nil {
			accrualValue = *result.Accrual
		}

		return s.storage.UpdateOrder(ctx, number, internalStatus, accrualValue)

	case http.StatusTooManyRequests:
		// логируем, но не падаем
		s.logger.Warn("accrual: too many requests", zap.String("url", u.String()))
		return nil

	default:
		// другие ошибки считаем временными
		s.logger.Warn("accrual: unexpected status", zap.Int("status", resp.StatusCode), zap.String("url", u.String()))
		return nil
	}
}
