package activities

import (
	"context"
	"github.com/vithubati/go-core/logger"
	"github.com/vithubati/shop-with-temporal/services/shopper/data"
	"go.temporal.io/sdk/temporal"
)

type OrderActivity struct {
	OrderStore data.OrderStore
}

func (a *OrderActivity) CreateOrder(ctx context.Context, products []string, workflowId string) (*data.Order, error) {
	logger.From(ctx).Infof("CreateOrder for %v", products)
	order, err := a.OrderStore.Persist(ctx, &data.Order{
		WorkflowId:  workflowId,
		Products:    products,
		OrderStatus: "pending",
	})
	if err != nil {
		return nil, temporal.NewNonRetryableApplicationError(err.Error(), "activity_error", err)
	}
	return order, nil
}

func (a *OrderActivity) CreateTransaction(ctx context.Context, orderID string) (err error) {
	// create transaction record
	logger.From(ctx).Infof("CreateTransaction for %s", orderID)
	return err
}

func (a *OrderActivity) CreatePackage(ctx context.Context, orderID string) (err error) {
	logger.From(ctx).Infof("CreatePackage for %s", orderID)
	return err
}

func (a *OrderActivity) ConfirmShipping(ctx context.Context, workflowId string, status string) (err error) {
	logger.From(ctx).Infof("ConfirmShipping for workflowId %s", workflowId)

	err = a.OrderStore.UpdateShippingStatus(ctx, workflowId, status)
	if err != nil {
		return temporal.NewNonRetryableApplicationError(err.Error(), "activity_error", err)
	}

	return err
}

func (a *OrderActivity) ConfirmOrder(ctx context.Context, workflowId string, status string) (err error) {
	logger.From(ctx).Infof("ConfirmOrder for workflowId %s", workflowId)

	err = a.OrderStore.UpdateStatus(ctx, workflowId, status)
	if err != nil {
		return temporal.NewNonRetryableApplicationError(err.Error(), "activity_error", err)
	}

	return err
}

func (a *OrderActivity) CreateInvoice(ctx context.Context, workflowId string) (err error) {
	// keep retying until shipping is confirmed
	logger.From(ctx).Infof("CreateInvoice for workflowId %s", workflowId)

	return err
}

func (a *OrderActivity) ConfirmInvoice(ctx context.Context, workflowId string, status string) (err error) {
	logger.From(ctx).Infof("ConfirmInvoice for workflowId %s", workflowId)
	err = a.OrderStore.UpdateInvoiceStatus(ctx, workflowId, status)
	if err != nil {
		return temporal.NewNonRetryableApplicationError(err.Error(), "activity_error", err)
	}

	return err
}
