package api

import (
	"github.com/vithubati/shop-with-temporal/services/shopper/data"
	"go.temporal.io/sdk/client"
)

type ShopAPI struct {
	productStore data.ProductStore
	cartStore    data.CartStore
	orderStore   data.OrderStore
	client       client.Client
}

func NewShopAPI(client client.Client, productStore data.ProductStore, cartStore data.CartStore, orderStore data.OrderStore) *ShopAPI {
	return &ShopAPI{
		client:       client,
		productStore: productStore,
		cartStore:    cartStore,
		orderStore:   orderStore,
	}
}
