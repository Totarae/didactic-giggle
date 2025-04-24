package handlers

import (
	"bytes"
	"didactic-giggle/cmd/internal/handlers/mocks"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockUserService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewUserHandler(mockSvc, logger)

	login := "testuser"
	password := "password123"

	mockSvc.EXPECT().CreateUser(gomock.Any(), login, password).Return(nil)

	body := []byte(`{"login":"testuser","password":"password123"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.RegisterHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	cookies := rec.Result().Cookies()
	assert.Len(t, cookies, 1)
	assert.Equal(t, "auth", cookies[0].Name)
	assert.Equal(t, login, cookies[0].Value)
}

func TestRegisterHandler_UserExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockUserService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewUserHandler(mockSvc, logger)

	mockSvc.EXPECT().CreateUser(gomock.Any(), "user", "pass").Return(ErrUserExists)

	body := []byte(`{"login":"user","password":"pass"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.RegisterHandler(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockUserService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewUserHandler(mockSvc, logger)

	body := []byte(`{not-json}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.RegisterHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRegisterHandler_EmptyFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockUserService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewUserHandler(mockSvc, logger)

	body := []byte(`{"login":"","password":""}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.RegisterHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRegisterHandler_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockUserService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewUserHandler(mockSvc, logger)

	mockSvc.EXPECT().CreateUser(gomock.Any(), "fail", "fail").Return(errors.New("db error"))

	body := []byte(`{"login":"fail","password":"fail"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.RegisterHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestLoginHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockUserService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewUserHandler(mockSvc, logger)

	mockSvc.EXPECT().AuthenticateUser(gomock.Any(), "user", "pass").Return(true, nil)

	body := []byte(`{"login":"user","password":"pass"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.LoginHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockUserService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewUserHandler(mockSvc, logger)

	mockSvc.EXPECT().AuthenticateUser(gomock.Any(), "user", "wrongpass").Return(false, nil)

	body := []byte(`{"login":"user","password":"wrongpass"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.LoginHandler(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLoginHandler_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockUserService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewUserHandler(mockSvc, logger)

	mockSvc.EXPECT().AuthenticateUser(gomock.Any(), "user", "pass").Return(false, errors.New("db error"))

	body := []byte(`{"login":"user","password":"pass"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.LoginHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockUserService(ctrl)
	logger := zaptest.NewLogger(t)
	handler := NewUserHandler(mockSvc, logger)

	body := []byte(`{invalid}`)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.LoginHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
