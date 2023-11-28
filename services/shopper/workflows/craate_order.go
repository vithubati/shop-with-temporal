package workflows

import (
	"github.com/vithubati/shop-with-temporal/services/shopper/data"
	"github.com/vithubati/shop-with-temporal/services/shopper/workflows/activities"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

func CreateOrderWorkflow(ctx workflow.Context, products []string) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 20,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Minute,
			BackoffCoefficient: 2,
			MaximumAttempts:    7,
			MaximumInterval:    time.Hour,
			NonRetryableErrorTypes: []string{
				"activity_error",
			},
		},
	}

	// start workflow with activity options
	ctx = workflow.WithActivityOptions(ctx, options)

	// configure logger
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting CreateOrderWorkflow")

	// get running workflowId
	workflowId := workflow.GetInfo(ctx).WorkflowExecution.ID

	// result can be string, struct, other data types.
	var confirmInvoiceResult string
	var confirmShippingResult string

	// workflow selector
	invoiceSelector := workflow.NewSelector(ctx)
	shippingSelector := workflow.NewSelector(ctx)

	// workflow named signal channel
	invoiceSignalChan := workflow.GetSignalChannel(ctx, "confirmInvoice")
	shippingSignalChan := workflow.GetSignalChannel(ctx, "confirmShipping")

	// implement selector reciever via signal channel
	invoiceSelector.AddReceive(invoiceSignalChan, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &confirmInvoiceResult)
	})

	shippingSelector.AddReceive(shippingSignalChan, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &confirmShippingResult)
	})
	// start the activites

	// Create Order
	var orderActivity *activities.OrderActivity
	order := &data.Order{}
	err := workflow.ExecuteActivity(ctx, orderActivity.CreateOrder, products, workflowId).Get(ctx, order)
	if err != nil {
		logger.Error("error executing CreateOrder")
		return err
	}
	// Create Transaction
	err = workflow.ExecuteActivity(ctx, orderActivity.CreateTransaction, order.ID).Get(ctx, nil)
	if err != nil {
		logger.Error("error executing CreateTransaction")
		return err
	}
	err = workflow.ExecuteActivity(ctx, orderActivity.CreateInvoice, order.ID).Get(ctx, nil)
	if err != nil {
		logger.Error("error executing CreatePackage")
		return err
	}
	invoiceSelector.Select(ctx)
	// Confirm Transaction invoice
	err = workflow.ExecuteActivity(ctx, orderActivity.ConfirmInvoice, workflowId, confirmInvoiceResult).Get(ctx, nil)
	if err != nil {
		logger.Error("error executing ConfirmInvoice")
		return err
	}

	err = workflow.ExecuteActivity(ctx, orderActivity.CreatePackage, order.ID).Get(ctx, nil)
	if err != nil {
		logger.Error("error executing CreatePackage")
		return err
	}

	shippingSelector.Select(ctx)

	// Confirm Shipping fulfillment
	err = workflow.ExecuteActivity(ctx, orderActivity.ConfirmShipping, workflowId, confirmShippingResult).Get(ctx, nil)
	if err != nil {
		logger.Error("error executing ConfirmShipping")
		return err
	}
	// Confirm Order
	err = workflow.ExecuteActivity(ctx, orderActivity.ConfirmOrder, workflowId, "done").Get(ctx, nil)
	if err != nil {
		logger.Error("error executing ConfirmOrder")
		return err
	}
	return nil
}
