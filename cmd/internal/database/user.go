package database

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserExists = errors.New("user already exists")

func (db *DB) CreateUser(ctx context.Context, login, password string) error {
	var exists bool
	err := db.Pool.QueryRow(ctx, SelectUserQuery, login).Scan(&exists)
	if err != nil {
		db.Logger.Error("ошибка при проверке существования пользователя", zap.Error(err))
		return err
	}
	if exists {
		return ErrUserExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		db.Logger.Error("ошибка хеширования пароля", zap.Error(err))
		return err
	}

	_, err = db.Pool.Exec(ctx, InsertUserQuery, login, password, string(hashed))
	if err != nil {
		db.Logger.Error("ошибка при вставке нового пользователя", zap.Error(err))
	}
	return err
}

func (db *DB) AuthenticateUser(ctx context.Context, login, password string) (bool, error) {
	var hashed string
	err := db.Pool.QueryRow(ctx, SelectUserHash, login).Scan(&hashed)
	if err != nil {
		db.Logger.Warn("пользователь не найден или ошибка при запросе", zap.Error(err))
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		db.Logger.Warn("неверный пароль при попытке входа", zap.String("login", login))
		return false, nil // пароль неверный
	}

	return true, nil
}

func (db *DB) GetUserIDByLogin(ctx context.Context, login string) (int, error) {
	var id int
	err := db.Pool.QueryRow(ctx, SelectUserByID, login).Scan(&id)
	return id, err
}
