package workflows

import (
	"github.com/vithubati/shop-with-temporal/services/shopper/data"
	"github.com/vithubati/shop-with-temporal/services/shopper/workflows/activities"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

// CreateCartWorkflow -  inserts products in a cart
func CreateCartWorkflow(ctx workflow.Context, cart *data.Cart) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var cartActivity *activities.CartActivity
	result := &data.Cart{}
	err := workflow.ExecuteActivity(ctx, cartActivity.InsertCart, cart).Get(ctx, result)
	if err != nil {
		return temporal.NewApplicationError(err.Error(), "error")
	}

	return nil
}
