package database

import (
	"context"
	"database/sql"
	"didactic-giggle/cmd/internal/util"
	"time"
)

func (db *DB) SaveOrder(ctx context.Context, userID int, order string) (util.OrderSaveStatus, error) {
	cmdTag, err := db.Pool.Exec(ctx, InsertOrderQuery, order, userID)

	if err != nil {
		return 0, err
	}
	if cmdTag.RowsAffected() == 1 {
		return util.OrderSavedNew, nil
	}
	// Заказ уже есть — проверим владельца
	var existingUserID int
	err = db.Pool.QueryRow(ctx, SelectOrderOwnerQuery, order).Scan(&existingUserID)

	if err != nil {
		return 0, err
	}

	if existingUserID == userID {
		return util.OrderAlreadyUploadedByUser, nil
	}
	return util.OrderUploadedByAnotherUser, nil
}

func (db *DB) GetOrders(ctx context.Context, userID int) ([]util.OrderInfo, error) {
	rows, err := db.Pool.Query(ctx, GetOrdersQuery, userID)
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

func (db *DB) GetPendingOrders(ctx context.Context) ([]string, error) {
	rows, err := db.Pool.Query(ctx, GetPendingOrdersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var numbers []string // собеерм в массив номеров заказов
	for rows.Next() {
		var number string
		if err := rows.Scan(&number); err != nil {
			return nil, err
		}
		numbers = append(numbers, number)
	}
	return numbers, nil
}

func (db *DB) UpdateOrder(ctx context.Context, number string, status string, accrual float64) error {
	_, err := db.Pool.Exec(ctx, UpdateOrderQuery, number, status, accrual)
	return err
}
