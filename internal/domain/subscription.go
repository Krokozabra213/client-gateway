package domain

import "time"

type Subscription struct {
	ID          string
	ProductID   string
	ProductName string
	Status      SubscriptionStatus
	ExpiresAt   time.Time
}

type SubscriptionStatus string

const (
	SubscriptionStatusUnspecified SubscriptionStatus = "UNSPECIFIED"
	SubscriptionStatusActive      SubscriptionStatus = "ACTIVE"
	SubscriptionStatusExpired     SubscriptionStatus = "EXPIRED"
	SubscriptionStatusCanceled    SubscriptionStatus = "CANCELED"
)
