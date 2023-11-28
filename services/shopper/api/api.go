package api

import (
	"github.com/go-chi/chi/v5"
	middleware2 "github.com/go-chi/chi/v5/middleware"
	"github.com/vithubati/go-core/http/middleware"
	"github.com/vithubati/shop-with-temporal/pkg/config"
	"net/http"
)

func New(cfg *config.Configuration, shopperAPI *ShopAPI) (http.Handler, error) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware2.Logger)
	r.Use(middleware2.Recoverer)
	r.Use(middleware2.URLFormat)
	//r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware2.Heartbeat("/health"))
	// game Routes
	r.Mount("/shopping", shopAPIRouter(shopperAPI))

	return r, nil
}

func shopAPIRouter(shopperAPI *ShopAPI) chi.Router {
	r := chi.NewRouter()
	r.Get("/carts", shopperAPI.CartGetHandler)     // curl -X GET http://localhost:8086/shopping/carts
	r.Post("/carts", shopperAPI.CartInsertHandler) // curl -X POST http://localhost:8086/shopping/carts\?products\=1,2,3
	r.Post("/orders", shopperAPI.OrderHandler)     // curl -X POST http://localhost:8086/shopping/orders?products=1,2,3
	// curl -X POST http://localhost:8086/shopping/orders/signal?orderId={orderID}&signalName=confirmInvoice&status=confirmed
	// curl -X POST http://localhost:8086/shopping/orders/signal?orderId={orderID}&signalName=confirmShipping&status=confirmed
	r.Post("/orders/signal", shopperAPI.SignalOrderHandler)
	return r
}
