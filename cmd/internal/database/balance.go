package database

import (
	"context"
	"didactic-giggle/cmd/internal/util"
	"time"
)

func (db *DB) GetUserBalance(ctx context.Context, userID int) (float64, float64, error) {
	var current float64
	var withdrawn float64

	// Сумма начислений
	err := db.Pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(accrual), 0)
		FROM orders
		WHERE user_id = $1 AND status = 'PROCESSED'
	`, userID).Scan(&current)
	if err != nil {
		return 0, 0, err
	}

	// Сумма списаний
	err = db.Pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(amount), 0)
		FROM withdrawals
		WHERE user_id = $1
	`, userID).Scan(&withdrawn)
	if err != nil {
		return 0, 0, err
	}

	return current - withdrawn, withdrawn, nil
}

func (db *DB) Withdraw(ctx context.Context, userID int, orderNumber string, amount float64) error {
	// Проверяем текущий баланс
	var current, withdrawn float64

	err := db.Pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(accrual), 0)
		FROM orders
		WHERE user_id = $1 AND status = 'PROCESSED'
	`, userID).Scan(&current)
	if err != nil {
		return err
	}

	err = db.Pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(amount), 0)
		FROM withdrawals
		WHERE user_id = $1
	`, userID).Scan(&withdrawn)
	if err != nil {
		return err
	}

	available := current - withdrawn
	if amount > available {
		return util.ErrInsufficientFunds
	}

	_, err = db.Pool.Exec(ctx, `
		INSERT INTO withdrawals (user_id, order_number, amount)
		VALUES ($1, $2, $3)
	`, userID, orderNumber, amount)
	return err
}

func (db *DB) GetWithdrawals(ctx context.Context, userID int) ([]util.Withdrawal, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT order_number, amount, processed_at
		FROM withdrawals
		WHERE user_id = $1
		ORDER BY processed_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []util.Withdrawal
	for rows.Next() {
		var w util.Withdrawal
		var processedAt time.Time

		if err := rows.Scan(&w.Order, &w.Sum, &processedAt); err != nil {
			return nil, err
		}
		w.ProcessedAt = processedAt.Format(time.RFC3339)
		result = append(result, w)
	}

	return result, nil
}
