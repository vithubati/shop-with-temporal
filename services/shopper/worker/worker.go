package worker

import (
	"github.com/vithubati/shop-with-temporal/services/shopper/data"
	"github.com/vithubati/shop-with-temporal/services/shopper/workflows"
	"github.com/vithubati/shop-with-temporal/services/shopper/workflows/activities"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
)

func Start(client client.Client, taskQueue string, orderStore data.OrderStore, cartStore data.CartStore) (func(), error) {
	w := worker.New(client, taskQueue, worker.Options{})
	w.RegisterWorkflow(workflows.CreateCartWorkflow)
	w.RegisterWorkflow(workflows.CreateOrderWorkflow)
	orderActivity := &activities.OrderActivity{OrderStore: orderStore}
	cartActivity := &activities.CartActivity{CartStore: cartStore}
	w.RegisterActivity(orderActivity)
	w.RegisterActivity(cartActivity)

	// start the worker and the web server
	err := w.Start()
	if err != nil {
		return func() {}, err
	}
	return func() {
		log.Println("stopping worker")
		w.Stop()
	}, nil
}
