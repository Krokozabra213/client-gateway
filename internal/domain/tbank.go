package domain

type TbankInitPayload struct {
	IdempotencyKey string
	UserID         string
	Amount         int
	Desc           string
}

func NewTbankInit(idempKey, userID, desc string, amount int) *TbankInitPayload {
	return &TbankInitPayload{
		IdempotencyKey: idempKey,
		UserID:         userID,
		Desc:           desc,
		Amount:         amount,
	}
}

type TBankPaymentNotification struct {
	TerminalKey string `json:"TerminalKey"`
	OrderID     string `json:"OrderId"`
	Success     bool   `json:"Success"`
	Status      string `json:"Status"`
	PaymentID   int64  `json:"PaymentId"`
	ErrorCode   string `json:"ErrorCode"`
	Amount      int32    `json:"Amount"`
	Token       string `json:"Token"`

	Message  *string `json:"Message,omitempty"`
	Details  *string `json:"Details,omitempty"`
	Pan      *string `json:"Pan,omitempty"`
	ExpDate  *string `json:"ExpDate,omitempty"`
	RebillID *int64  `json:"RebillId,omitempty"`
	CardID   *int64  `json:"CardId,omitempty"`
}
