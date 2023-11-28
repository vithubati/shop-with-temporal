package api

import (
	"encoding/json"
	"fmt"
	"github.com/vithubati/shop-with-temporal/services/shopper/workflows"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

func (a *ShopAPI) OrderHandler(w http.ResponseWriter, r *http.Request) {
	products := r.URL.Query().Get("products")
	productsArr := strings.Split(products, ",")

	_, err := a.client.ExecuteWorkflow(r.Context(), client.StartWorkflowOptions{
		ID:        "CreateOrderWorkflow_" + uuid.New().String(),
		TaskQueue: TQShoppingCart,
	}, workflows.CreateOrderWorkflow, productsArr)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to start workflow. err: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal("{status: 'ok'}")
	if err != nil {
		http.Error(w, "unable to marshal response", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func (a *ShopAPI) SignalOrderHandler(w http.ResponseWriter, r *http.Request) {
	orderId := r.URL.Query().Get("orderId")
	signalName := r.URL.Query().Get("signalName")
	orderStatus := r.URL.Query().Get("status")

	// get order from data layer
	order, err := a.orderStore.Get(r.Context(), orderId)
	if err != nil {
		http.Error(w, "unable to get order", http.StatusInternalServerError)
		return
	}
	fmt.Println(fmt.Sprintf("Oder is: %+v", order))
	fmt.Println(fmt.Sprintf("signalName: %s", signalName))
	fmt.Println(fmt.Sprintf("orderStatus: %s", orderStatus))
	workflowId := order.WorkflowId
	runId := "" // we did not store runId we can safely leave it empty
	err = a.client.SignalWorkflow(r.Context(), workflowId, runId, signalName, orderStatus)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to signal workflow. err: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal("{status: 'ok'}")
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to marshal response. Err: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}
