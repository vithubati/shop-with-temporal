package api

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/vithubati/shop-with-temporal/services/shopper/data"
	"github.com/vithubati/shop-with-temporal/services/shopper/workflows"
	"go.temporal.io/sdk/client"
	"net/http"
	"strings"
)

const (
	TQShoppingCart = "SHOPPING_CART"
)

func (a *ShopAPI) CartGetHandler(w http.ResponseWriter, r *http.Request) {
	cart, err := a.cartStore.Get(r.Context())
	jsonResponse, err := json.Marshal(cart)
	if err != nil {
		http.Error(w, "unable to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "unable to write response", http.StatusInternalServerError)
		return
	}
}

func (a *ShopAPI) CartInsertHandler(w http.ResponseWriter, r *http.Request) {
	productString := r.URL.Query().Get("products")
	// split productString
	stringArr := strings.Split(productString, ",")
	cart := &data.Cart{
		Products: stringArr,
	}
	we, err := a.client.ExecuteWorkflow(r.Context(), client.StartWorkflowOptions{
		ID:        "SetCartWorkflow_" + uuid.New().String(),
		TaskQueue: TQShoppingCart,
	}, workflows.CreateCartWorkflow, cart)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to start workflow %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if err := we.Get(r.Context(), nil); err != nil {
		http.Error(w, "unable to get workflow result", http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal("{status: 'ok'}")
	if err != nil {
		http.Error(w, "unable to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "unable to write response", http.StatusInternalServerError)
		return
	}
}
