package util

type OrderSaveStatus int

const (
	OrderSavedNew OrderSaveStatus = iota
	OrderAlreadyUploadedByUser
	OrderUploadedByAnotherUser
)

// OrderInfo — структура для ответа GET /api/user/orders
type OrderInfo struct {
	Number     string   `json:"number"`            // номер заказа
	Status     string   `json:"status"`            // статус заказа
	Accrual    *float64 `json:"accrual,omitempty"` // начисление, может быть null
	UploadedAt string   `json:"uploaded_at"`       // RFC3339
}
