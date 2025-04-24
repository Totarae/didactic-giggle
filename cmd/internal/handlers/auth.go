package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

var ErrUserExists = errors.New("user already exists")

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserService interface {
	CreateUser(ctx context.Context, login, password string) error
	AuthenticateUser(ctx context.Context, login, password string) (bool, error)
	GetUserIDByLogin(ctx context.Context, login string) (int, error)
}

// структура хэндлера
type UserHandler struct {
	userService UserService
	Logger      *zap.Logger
}

// конструктор
func NewUserHandler(userService UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		Logger:      logger,
		userService: userService,
	}
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		h.Logger.Warn("ошибка парсинга запроса на регистрацию", zap.Error(err))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if creds.Login == "" || creds.Password == "" {
		http.Error(w, "login and password required", http.StatusBadRequest)
		return
	}

	err := h.userService.CreateUser(r.Context(), creds.Login, creds.Password)
	if err != nil {
		if errors.Is(err, ErrUserExists) {
			http.Error(w, "login already in use", http.StatusConflict)
		} else {
			h.Logger.Error("failed to register user", zap.Error(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	// Установка cookie для аутентификации
	setAuthCookie(w, creds.Login)
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		h.Logger.Warn("Ошибка парсинга логина", zap.Error(err))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	ok, err := h.userService.AuthenticateUser(r.Context(), creds.Login, creds.Password)
	if err != nil {
		h.Logger.Error("auth failed", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	setAuthCookie(w, creds.Login)
	w.WriteHeader(http.StatusOK)
}

func setAuthCookie(w http.ResponseWriter, login string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    login,
		Path:     "/",
		HttpOnly: true,
	})
}
