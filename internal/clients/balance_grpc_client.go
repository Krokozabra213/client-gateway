package clients

import (
	"context"
	"fmt"

	balancev1 "github.com/Krokozabra213/gargantua_common/gen/go/proto/balance/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BalanceClient struct {
	api  balancev1.BalanceService_APIClient
	conn *grpc.ClientConn
}

func NewBalance(ctx context.Context, addr string) (*BalanceClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("balance client: grpc new client: %w", err)
	}

	return &BalanceClient{
		api:  balancev1.NewBalanceService_APIClient(conn),
		conn: conn,
	}, nil
}

func (c *BalanceClient) Close() error {
	return c.conn.Close()
}

func (c *BalanceClient) GetAccount(ctx context.Context, userID string) (*balancev1.GetAccountResponse, error) {
	resp, err := c.api.GetAccount(ctx, &balancev1.GetAccountRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("balance client: get account: %w", err)
	}
	return resp, nil
}
