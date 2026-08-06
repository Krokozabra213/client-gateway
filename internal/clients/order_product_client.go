package clients

import (
	"context"
	"fmt"

	"github.com/GargantuaLabs/client-gateway/internal/domain"
	orderproductv1 "github.com/Krokozabra213/gargantua_common/gen/go/proto/order-product/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type OrderProductClient struct {
	api  orderproductv1.OrderProductService_APIClient
	conn *grpc.ClientConn
}

func NewOrderProduct(ctx context.Context, addr string) (*OrderProductClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("balance client: grpc new client: %w", err)
	}

	return &OrderProductClient{
		api:  orderproductv1.NewOrderProductService_APIClient(conn),
		conn: conn,
	}, nil
}

func (c *OrderProductClient) Close() error {
	return c.conn.Close()
}

func (c *OrderProductClient) ListProducts(ctx context.Context, onlyActive bool) ([]domain.Product, error) {
	resp, err := c.api.ListProducts(ctx, &orderproductv1.ListProductsRequest{
		OnlyActive: onlyActive,
	})
	if err != nil {
		return nil, fmt.Errorf("order-product client: list products: %w", err)
	}

	products := make([]domain.Product, 0, len(resp.Products))
	for _, p := range resp.Products {
		bundleItems := make([]domain.BundleItem, 0, len(p.BundleItems))
		for _, b := range p.BundleItems {
			bundleItems = append(bundleItems, domain.BundleItem{
				ProductID: b.ProductId,
				Name:      b.Name,
			})
		}

		productType, err := mapProductType(p.Type)
		if err != nil {
			return nil, fmt.Errorf("order-product client: map product type: %w", err)
		}

		products = append(products, domain.Product{
			ID:            p.Id,
			Type:          productType,
			Name:          p.Name,
			Description:   p.Description,
			DurationDays:  int(p.DurationDays),
			PriceAmount:   p.PriceAmount,
			PriceCurrency: p.PriceCurrency,
			BundleItems:   bundleItems,
		})
	}

	return products, nil
}

func (c *OrderProductClient) GetProduct(ctx context.Context, productID string) (*domain.Product, error) {
	resp, err := c.api.GetProduct(ctx, &orderproductv1.GetProductRequest{
		ProductId: productID,
	})
	if err != nil {
		return nil, fmt.Errorf("order-product client: list products: %w", err)
	}

	if err != nil {

		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("order-product client: get product: %w", err)
	}

	if resp.GetProduct() == nil {
		return nil, domain.ErrNotFound
	}

	p := resp.GetProduct()
	productType, err := mapProductType(p.Type)
	if err != nil {
		return nil, fmt.Errorf("order-product client: map product type: %w", err)
	}

	bundleItems := make([]domain.BundleItem, 0, len(p.BundleItems))
	for _, b := range p.BundleItems {
		bundleItems = append(bundleItems, domain.BundleItem{
			ProductID: b.ProductId,
			Name:      b.Name,
		})
	}

	product := &domain.Product{
		ID:            p.Id,
		Type:          productType,
		Name:          p.GetName(),
		Description:   p.GetDescription(),
		DurationDays:  int(p.GetDurationDays()),
		PriceAmount:   p.GetPriceAmount(),
		PriceCurrency: p.GetPriceCurrency(),
		BundleItems:   bundleItems,
	}

	return product, nil

}

func (c *OrderProductClient) ListUserOrders(ctx context.Context, userID string, limit, offset int32) ([]domain.Order, error) {
	resp, err := c.api.ListUserOrders(ctx, &orderproductv1.ListUserOrdersRequest{
		UserId: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("order-product client: list user orders: %w", err)
	}

	orders := make([]domain.Order, 0, len(resp.Orders))
	for _, o := range resp.Orders {

		productType, err := mapProductType(o.ProductType)
		if err != nil {
			return nil, fmt.Errorf("order-product client: map product type: %w", err)
		}

		orderStatus, err := mapOrderStatus(o.Status)
		if err != nil {
			return nil, fmt.Errorf("order-product client: map order status: %w", err)
		}

		order := domain.Order{
			ID:           o.Id,
			ProductID:    o.ProductId,
			ProductName:  o.ProductName,
			ProductType:  productType,
			DurationDays: int(o.DurationDays),
			Amount:       o.Amount,
			Currency:     o.Currency,
			Status:       orderStatus,
		}

		if o.CreatedAt != nil {
			order.CreatedAt = o.CreatedAt.AsTime()
		}
		if o.PaidAt != nil {
			order.PaidAt = o.PaidAt.AsTime()
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (c *OrderProductClient) ListUserSubscriptions(
	ctx context.Context, userID string,
) ([]domain.Subscription, error) {

	resp, err := c.api.ListUserSubscriptions(ctx, &orderproductv1.ListUserSubscriptionsRequest{
		UserId: userID,
		Status: orderproductv1.SubscriptionStatus_SUBSCRIPTION_STATUS_ACTIVE,
	})
	if err != nil {
		return nil, fmt.Errorf("order-product client: list user subscriptions: %w", err)
	}

	subs := make([]domain.Subscription, 0, len(resp.Subscriptions))
	for _, s := range resp.Subscriptions {
		subscriptionStatus, err := mapSubscriptionStatus(s.Status)
		if err != nil {
			return nil, fmt.Errorf("order-product client: map subscription status: %w", err)
		}

		sub := domain.Subscription{
			ID:          s.Id,
			ProductID:   s.ProductId,
			ProductName: s.ProductName,
			Status:      subscriptionStatus,
		}

		if s.ExpiresAt != nil {
			sub.ExpiresAt = s.ExpiresAt.AsTime()
		}

		subs = append(subs, sub)
	}

	return subs, nil
}

func (c *OrderProductClient) PurchaseProduct(
	ctx context.Context, userID, productID string,
) (
	orderID string, orderStatus domain.OrderStatus, err error,
) {
	resp, err := c.api.PurchaseProduct(ctx, &orderproductv1.PurchaseProductRequest{
		UserId:    userID,
		ProductId: productID,
	})
	if err != nil {
		return "", domain.OrderStatus(""), fmt.Errorf("order-product client: purchase product: %w", err)
	}

	status, err := mapOrderStatus(resp.Status)
	if err != nil {
		return "", domain.OrderStatus(""), fmt.Errorf("order-product client: purchase product: %w", err)
	}

	return resp.OrderId, status, nil
}

func (c *OrderProductClient) CancelSubscription(ctx context.Context, userID, productID string) (bool, error) {
	resp, err := c.api.CancelSubscription(ctx, &orderproductv1.CancelSubscriptionRequest{
		UserId:    userID,
		ProductId: productID,
	})
	if err != nil {
		return false, fmt.Errorf("order-product client: cancel subscription: %w", err)
	}

	return resp.Success, nil
}

func (c *OrderProductClient) GetOrder(ctx context.Context, orderID string) (*domain.OrderDetails, error) {
	resp, err := c.api.GetOrder(ctx, &orderproductv1.GetOrderRequest{
		OrderId: orderID,
	})
	if err != nil {
		return nil, fmt.Errorf("order-product client: get order: %w", err)
	}

	status, err := mapOrderStatus(resp.Status)
	if err != nil {
		return nil, fmt.Errorf("order-product client: get order %s: %w", orderID, err)
	}

	return &domain.OrderDetails{
		OrderID:   resp.OrderId,
		UserID:    resp.UserId,
		ProductID: resp.ProductId,
		Status:    status,
		Amount:    resp.Amount,
		Currency:  resp.Currency,
	}, nil
}

func mapOrderStatus(s orderproductv1.OrderStatus) (domain.OrderStatus, error) {
	switch s {
	case orderproductv1.OrderStatus_ORDER_STATUS_UNSPECIFIED:
		return domain.OrderStatusUnspecified, nil
	case orderproductv1.OrderStatus_ORDER_STATUS_PENDING:
		return domain.OrderStatusPending, nil
	case orderproductv1.OrderStatus_ORDER_STATUS_COMPLETED:
		return domain.OrderStatusCompleted, nil
	case orderproductv1.OrderStatus_ORDER_STATUS_FAILED:
		return domain.OrderStatusFailed, nil
	case orderproductv1.OrderStatus_ORDER_STATUS_CANCELLED:
		return domain.OrderStatusCancelled, nil
	case orderproductv1.OrderStatus_ORDER_STATUS_REFUNDED:
		return domain.OrderStatusRefunded, nil
	case orderproductv1.OrderStatus_ORDER_STATUS_EXPIRED:
		return domain.OrderStatusExpired, nil
	default:
		return "", fmt.Errorf("unknown order status: %d", s)
	}
}

func mapSubscriptionStatus(s orderproductv1.SubscriptionStatus) (domain.SubscriptionStatus, error) {
	switch s {
	case orderproductv1.SubscriptionStatus_SUBSCRIPTION_STATUS_UNSPECIFIED:
		return domain.SubscriptionStatusUnspecified, nil
	case orderproductv1.SubscriptionStatus_SUBSCRIPTION_STATUS_ACTIVE:
		return domain.SubscriptionStatusActive, nil
	case orderproductv1.SubscriptionStatus_SUBSCRIPTION_STATUS_EXPIRED:
		return domain.SubscriptionStatusExpired, nil
	case orderproductv1.SubscriptionStatus_SUBSCRIPTION_STATUS_CANCELED:
		return domain.SubscriptionStatusCanceled, nil
	default:
		return "", fmt.Errorf("unknown subscription status: %d", s)
	}
}
func mapProductType(t orderproductv1.ProductType) (domain.ProductType, error) {
	switch t {
	case orderproductv1.ProductType_PRODUCT_TYPE_UNSPECIFIED:
		return domain.ProductTypeUnspecified, nil
	case orderproductv1.ProductType_PRODUCT_TYPE_SUBSCRIPTION:
		return domain.ProductTypeSubscription, nil
	case orderproductv1.ProductType_PRODUCT_TYPE_BUNDLE:
		return domain.ProductTypeBundle, nil
	default:
		return "", fmt.Errorf("unknown product type: %d", t)
	}
}
