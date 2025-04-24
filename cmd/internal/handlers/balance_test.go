package handlers

import (
	"context"
	"didactic-giggle/cmd/internal/handlers/mocks"
	"didactic-giggle/cmd/internal/middleware"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

func TestGetBalanceHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockBalanceService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewBalanceHandler(mockService, logger)

	userID := 42
	current := 100.5
	withdrawn := 25.0

	mockService.EXPECT().
		GetUserBalance(gomock.Any(), userID).
		Return(current, withdrawn, nil)

	// создаем запрос с userID в контексте
	req := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
	ctx := middleware.WithUserID(context.Background(), userID)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	handler.GetBalanceHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	expected := `{"current":100.5,"withdrawn":25}`
	assert.JSONEq(t, expected, rec.Body.String())
}

func TestGetBalanceHandler_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockBalanceService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewBalanceHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
	rec := httptest.NewRecorder()

	handler.GetBalanceHandler(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetBalanceHandler_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockBalanceService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewBalanceHandler(mockService, logger)

	userID := 42

	mockService.EXPECT().
		GetUserBalance(gomock.Any(), userID).
		Return(0, 0, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
	ctx := middleware.WithUserID(context.Background(), userID)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	handler.GetBalanceHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
