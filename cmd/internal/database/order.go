package database

import (
	"context"
	"database/sql"
	"didactic-giggle/cmd/internal/util"
	"time"
)

func (db *DB) SaveOrder(ctx context.Context, userID int, order string) (util.OrderSaveStatus, error) {
	cmdTag, err := db.Pool.Exec(ctx, `
		INSERT INTO orders (order_number, user_id)
		VALUES ($1, $2)
		ON CONFLICT (order_number) DO NOTHING
	`, order, userID)

	if err != nil {
		return 0, err
	}
	if cmdTag.RowsAffected() == 1 {
		return util.OrderSavedNew, nil
	}
	// Заказ уже есть — проверим владельца
	var existingUserID int
	err = db.Pool.QueryRow(ctx, `
		SELECT user_id FROM orders WHERE order_number = $1
	`, order).Scan(&existingUserID)

	if err != nil {
		return 0, err
	}

	if existingUserID == userID {
		return util.OrderAlreadyUploadedByUser, nil
	}
	return util.OrderUploadedByAnotherUser, nil
}

func (db *DB) GetOrders(ctx context.Context, userID int) ([]util.OrderInfo, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT order_number, status, accrual, uploaded_at
		FROM orders
		WHERE user_id = $1
		ORDER BY uploaded_at
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []util.OrderInfo
	for rows.Next() {
		var o util.OrderInfo
		var accrual sql.NullFloat64
		var uploadedAt time.Time

		err := rows.Scan(&o.Number, &o.Status, &accrual, &uploadedAt)
		if err != nil {
			return nil, err
		}

		if accrual.Valid {
			o.Accrual = &accrual.Float64
		}

		o.UploadedAt = uploadedAt.Format(time.RFC3339)
		orders = append(orders, o)
	}
	return orders, nil
}
