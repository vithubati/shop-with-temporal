package activities

import (
	"context"
	"github.com/vithubati/shop-with-temporal/services/shopper/data"
)

type CartActivity struct {
	CartStore data.CartStore
}

func (a *CartActivity) InsertCart(ctx context.Context, cart *data.Cart) error {
	_, err := a.CartStore.Persist(ctx, cart)
	if err != nil {
		return err
	}
	return nil
}
