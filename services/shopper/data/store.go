package data

import "context"

type Store interface {
	Migrate() error
}
type ProductStore interface {
	Store
	List(ctx context.Context) ([]Product, error)
}

type CartStore interface {
	Store
	Persist(ctx context.Context, cart *Cart) (*Cart, error)
	Get(ctx context.Context) (*Cart, error)
}

type OrderStore interface {
	Store
	Persist(ctx context.Context, order *Order) (*Order, error)
	Get(ctx context.Context, orderId string) (order *Order, err error)
	UpdateInvoiceStatus(ctx context.Context, workflowId string, status string) (err error)
	UpdateShippingStatus(ctx context.Context, workflowId string, status string) (err error)
	UpdateStatus(ctx context.Context, workflowId string, status string) (err error)
}
