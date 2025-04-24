package handlers

import (
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

func TestGetOrdersHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader := mocks.NewMockOrderReader(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewOrderReaderHandler(mockReader, logger)

	expectedOrders := []util.OrderInfo{
		{Number: "123", Status: "NEW"},
		{Number: "456", Status: "PROCESSED"},
	}

	mockReader.
		EXPECT().
		GetOrders(gomock.Any(), 42).
		Return(expectedOrders, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
	req = req.WithContext(middleware.WithUserID(req.Context(), 42))
	rr := httptest.NewRecorder()

	handler.GetOrdersHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var got []util.OrderInfo
	err := json.NewDecoder(rr.Body).Decode(&got)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrders, got)
}

func TestGetOrdersHandler_NoContent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader := mocks.NewMockOrderReader(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewOrderReaderHandler(mockReader, logger)

	mockReader.
		EXPECT().
		GetOrders(gomock.Any(), 7).
		Return([]util.OrderInfo{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
	req = req.WithContext(middleware.WithUserID(req.Context(), 7))
	rr := httptest.NewRecorder()

	handler.GetOrdersHandler(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Empty(t, rr.Body.String())
}

func TestGetOrdersHandler_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader := mocks.NewMockOrderReader(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewOrderReaderHandler(mockReader, logger)

	mockReader.
		EXPECT().
		GetOrders(gomock.Any(), 99).
		Return(nil, errors.New("db failure"))

	req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
	req = req.WithContext(middleware.WithUserID(req.Context(), 99))
	rr := httptest.NewRecorder()

	handler.GetOrdersHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetOrdersHandler_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader := mocks.NewMockOrderReader(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewOrderReaderHandler(mockReader, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
	rr := httptest.NewRecorder()

	handler.GetOrdersHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
