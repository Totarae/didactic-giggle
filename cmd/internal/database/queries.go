package database

const (
	InsertOrderQuery = `
		INSERT INTO orders (order_number, user_id)
		VALUES ($1, $2)
		ON CONFLICT (order_number) DO NOTHING
	`

	SelectOrderOwnerQuery = `
		SELECT user_id FROM orders
		WHERE order_number = $1
	`

	GetOrdersQuery = `
		SELECT order_number, status, accrual, uploaded_at
		FROM orders
		WHERE user_id = $1
		ORDER BY uploaded_at
	`

	GetPendingOrdersQuery = `
		SELECT order_number
		FROM orders
		WHERE status IN ('NEW', 'PROCESSING')
	`

	UpdateOrderQuery = `
		UPDATE orders
		SET status = $2, accrual = $3
		WHERE order_number = $1
	`
	SelectUserQuery  = `SELECT EXISTS (SELECT 1 FROM users WHERE login = $1)`
	InsertUserQuery  = `INSERT INTO users (login, password, password_hash) VALUES ($1, $2, $3)`
	SelectUserHash   = `SELECT password_hash FROM users WHERE login = $1`
	SelectUserByID   = `SELECT id FROM users WHERE login = $1`
	SelectAccrualSum = `
		SELECT COALESCE(SUM(accrual), 0)
		FROM orders
		WHERE user_id = $1 AND status = 'PROCESSED'
	`
	SelectAmount = `
		SELECT COALESCE(SUM(amount), 0)
		FROM withdrawals
		WHERE user_id = $1
	`
	InsertWithdrawals = `
		INSERT INTO withdrawals (user_id, order_number, amount)
		VALUES ($1, $2, $3)
	`
	SelectWithdrawals = `
		SELECT order_number, amount, processed_at
		FROM withdrawals
		WHERE user_id = $1
		ORDER BY processed_at DESC
	`
)
