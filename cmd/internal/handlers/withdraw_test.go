package handlers

import (
	"bytes"
	"didactic-giggle/cmd/internal/handlers/mocks"
	"didactic-giggle/cmd/internal/middleware"
	"didactic-giggle/cmd/internal/util"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithdrawHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockWithdrawService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewWithdrawHandler(mockSvc, logger)

	body := `{"order":"79927398713","sum":100.5}`

	mockSvc.EXPECT().
		Withdraw(gomock.Any(), 1, "79927398713", 100.5).
		Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", bytes.NewBufferString(body))
	req = req.WithContext(middleware.WithUserID(req.Context(), 1))
	rr := httptest.NewRecorder()

	handler.WithdrawHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestWithdrawHandler_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockWithdrawService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewWithdrawHandler(mockSvc, logger)

	req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", bytes.NewBufferString(`{invalid}`))
	req = req.WithContext(middleware.WithUserID(req.Context(), 1))
	rr := httptest.NewRecorder()

	handler.WithdrawHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestWithdrawHandler_InvalidLuhn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockWithdrawService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewWithdrawHandler(mockSvc, logger)

	body := `{"order":"123456789","sum":50}` // некорректен по Луну

	req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", bytes.NewBufferString(body))
	req = req.WithContext(middleware.WithUserID(req.Context(), 1))
	rr := httptest.NewRecorder()

	handler.WithdrawHandler(rr, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestWithdrawHandler_InsufficientFunds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockWithdrawService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewWithdrawHandler(mockSvc, logger)

	body := `{"order":"79927398713","sum":500}`

	mockSvc.EXPECT().
		Withdraw(gomock.Any(), 1, "79927398713", 500.0).
		Return(util.ErrInsufficientFunds)

	req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", bytes.NewBufferString(body))
	req = req.WithContext(middleware.WithUserID(req.Context(), 1))
	rr := httptest.NewRecorder()

	handler.WithdrawHandler(rr, req)

	assert.Equal(t, http.StatusPaymentRequired, rr.Code)
}

func TestWithdrawHandler_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockWithdrawService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewWithdrawHandler(mockSvc, logger)

	body := `{"order":"79927398713","sum":100}`

	mockSvc.EXPECT().
		Withdraw(gomock.Any(), 1, "79927398713", 100.0).
		Return(errors.New("db error"))

	req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", bytes.NewBufferString(body))
	req = req.WithContext(middleware.WithUserID(req.Context(), 1))
	rr := httptest.NewRecorder()

	handler.WithdrawHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestWithdrawHandler_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockWithdrawService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewWithdrawHandler(mockSvc, logger)

	req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", nil)
	rr := httptest.NewRecorder()

	handler.WithdrawHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetWithdrawalsHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockWithdrawService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewWithdrawHandler(mockSvc, logger)

	expected := []util.Withdrawal{
		{Order: "123", Sum: 50, ProcessedAt: "2025-04-25T10:00:00Z"},
	}

	mockSvc.EXPECT().
		GetWithdrawals(gomock.Any(), 1).
		Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
	req = req.WithContext(middleware.WithUserID(req.Context(), 1))
	rr := httptest.NewRecorder()

	handler.GetWithdrawalsHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got []util.Withdrawal
	err := json.NewDecoder(rr.Body).Decode(&got)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestGetWithdrawalsHandler_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockWithdrawService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewWithdrawHandler(mockSvc, logger)

	mockSvc.EXPECT().
		GetWithdrawals(gomock.Any(), 1).
		Return([]util.Withdrawal{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
	req = req.WithContext(middleware.WithUserID(req.Context(), 1))
	rr := httptest.NewRecorder()

	handler.GetWithdrawalsHandler(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestGetWithdrawalsHandler_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockWithdrawService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewWithdrawHandler(mockSvc, logger)

	mockSvc.EXPECT().
		GetWithdrawals(gomock.Any(), 1).
		Return(nil, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
	req = req.WithContext(middleware.WithUserID(req.Context(), 1))
	rr := httptest.NewRecorder()

	handler.GetWithdrawalsHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetWithdrawalsHandler_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockWithdrawService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewWithdrawHandler(mockSvc, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
	rr := httptest.NewRecorder()

	handler.GetWithdrawalsHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
