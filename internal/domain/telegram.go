package domain

type TGInitPayload struct {
	IdempotencyKey string
	UserID         string
	Amount         int
	Desc           string
	Currency       string
	Title          string
}

func NewTgInit(idempKey, userID, desc, currency, title string, amount int) *TGInitPayload {
	return &TGInitPayload{
		IdempotencyKey: idempKey,
		UserID:         userID,
		Desc:           desc,
		Amount:         amount,
		Currency:       currency,
		Title:          title,
	}
}

type PreCheckoutPayload struct {
	ID             string `json:"id"`
	From           User   `json:"from"`
	Currency       string `json:"currency"`
	TotalAmount    int    `json:"total_amount"`
	InvoicePayload string `json:"invoice_payload"`
}

type User struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
}

type SuccessfulPayload struct {
	Currency                string `json:"currency"`
	InvoicePayload          string `json:"invoice_payload"`
	TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
	ProviderPaymentChargeID string `json:"provider_payment_charge_id"`
	TotalAmount             int    `json:"total_amount"`
}
