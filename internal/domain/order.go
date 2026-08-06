package domain

import "time"

type OrderStatus string

const (
	OrderStatusUnspecified OrderStatus = "UNSPECIFIED"
	OrderStatusPending     OrderStatus = "PENDING"
	OrderStatusCompleted   OrderStatus = "COMPLETED"
	OrderStatusFailed      OrderStatus = "FAILED"
	OrderStatusCancelled   OrderStatus = "CANCELLED"
	OrderStatusRefunded    OrderStatus = "REFUNDED"
	OrderStatusExpired     OrderStatus = "EXPIRED"
)

type Order struct {
	ID           string
	ProductID    string
	ProductName  string
	ProductType  ProductType
	DurationDays int
	Amount       int64
	Currency     string
	Status       OrderStatus
	CreatedAt    time.Time
	PaidAt       time.Time
}

type OrderDetails struct {
	OrderID   string
	UserID    string
	ProductID string
	Status    OrderStatus
	Amount    int64
	Currency  string
}
