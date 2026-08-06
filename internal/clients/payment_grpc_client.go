package clients

import (
	"context"
	"fmt"

	"github.com/GargantuaLabs/client-gateway/internal/domain"
	paymentv1 "github.com/Krokozabra213/gargantua_common/gen/go/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentClient struct {
	api  paymentv1.PaymentService_APIClient
	conn *grpc.ClientConn
}

func NewPayment(ctx context.Context, addr string) (*PaymentClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("payment client: grpc new client: %w", err)
	}

	return &PaymentClient{
		api:  paymentv1.NewPaymentService_APIClient(conn),
		conn: conn,
	}, nil
}

func (c *PaymentClient) Close() error {
	return c.conn.Close()
}

func (c *PaymentClient) TbankInit(ctx context.Context, payload *domain.TbankInitPayload) (*paymentv1.TbankInitResponse, error) {
	resp, err := c.api.TbankInit(ctx, &paymentv1.TbankInitRequest{
		IdempotencyKey: payload.IdempotencyKey,
		UserId:         payload.UserID,
		Amount:         int32(payload.Amount),
		Description:    payload.Desc,
	})
	if err != nil {
		return nil, fmt.Errorf("payment client: tbank init: %w", err)
	}
	return resp, nil
}

func (c *PaymentClient) TGInit(ctx context.Context, payload *domain.TGInitPayload) (*paymentv1.TGInitResponse, error) {
	resp, err := c.api.TGInit(ctx, &paymentv1.TGInitRequest{
		IdempotencyKey: payload.IdempotencyKey,
		UserId:         payload.UserID,
		Amount:         int32(payload.Amount),
		Description:    payload.Desc,
		Currency:       payload.Currency,
		Title:          payload.Title,
	})
	if err != nil {
		return nil, fmt.Errorf("payment client: tg init: %w", err)
	}
	return resp, nil
}

func (c *PaymentClient) PaymentCancel(ctx context.Context, paymentID string) error {
	_, err := c.api.PaymentCancel(ctx, &paymentv1.PaymentCancelRequest{
		PaymentId: paymentID,
	})
	if err != nil {
		return fmt.Errorf("payment client: cancel: %w", err)
	}
	return nil
}

func (c *PaymentClient) TbankWebhook(ctx context.Context, payload *domain.TBankPaymentNotification) error {
	_, err := c.api.TbankWebhook(ctx, &paymentv1.TbankWebhookRequest{
		TerminalKey: payload.TerminalKey,
		Amount:      payload.Amount,
		OrderId:     payload.OrderID,
		Success:     payload.Success,
		Status:      payload.Status,
		PaymentId:   payload.PaymentID,
		ErrorCode:   payload.ErrorCode,
		Message:     payload.Message,
		Details:     payload.Details,
		RebillId:    payload.RebillID,
		CardId:      payload.CardID,
		Pan:         payload.Pan,
		ExpDate:     payload.ExpDate,
		Token:       payload.Token,
	})
	if err != nil {
		return fmt.Errorf("payment client: tbank webhook: %w", err)
	}
	return nil
}

func (c *PaymentClient) TGPreCheckout(ctx context.Context, payload *domain.PreCheckoutPayload) error {
	_, err := c.api.TGPreCheckout(ctx, &paymentv1.TGPreCheckoutRequest{
		Id:             payload.ID,
		Currency:       payload.Currency,
		TotalAmount:    int64(payload.TotalAmount),
		InvoicePayload: payload.InvoicePayload,
		From: &paymentv1.TGUser{
			Id:        payload.From.ID,
			IsBot:     payload.From.IsBot,
			FirstName: payload.From.FirstName,
		},
	})
	if err != nil {
		return fmt.Errorf("payment client: tg pre-checkout: %w", err)
	}
	return nil
}

func (c *PaymentClient) TGSuccessful(ctx context.Context, payload *domain.SuccessfulPayload) error {
	_, err := c.api.TGSuccessful(ctx, &paymentv1.TGSuccessfulRequest{
		Currency:                payload.Currency,
		InvoicePayload:          payload.InvoicePayload,
		TelegramPaymentChargeId: payload.TelegramPaymentChargeID,
		ProviderPaymentChargeId: payload.ProviderPaymentChargeID,
		TotalAmount:             int64(payload.TotalAmount),
	})
	if err != nil {
		return fmt.Errorf("payment client: tg successful: %w", err)
	}
	return nil
}
