package shopper

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vithubati/go-core/logger"
	"github.com/vithubati/go-core/service"
	"github.com/vithubati/shop-with-temporal/pkg/config"
	"github.com/vithubati/shop-with-temporal/pkg/database"
	"github.com/vithubati/shop-with-temporal/services/shopper/api"
	"github.com/vithubati/shop-with-temporal/services/shopper/data"
	"github.com/vithubati/shop-with-temporal/services/shopper/worker"
	"go.temporal.io/sdk/client"
	"io/fs"
)

func NewApp(ctx context.Context, configs *config.Configuration) (service.Service, func(), error) {
	dbConf := database.DBConfiguration{
		Dialect:  database.Dialect(configs.Services.Shopper.DB.Dialect),
		Database: configs.Services.Shopper.DB.Database,
	}
	db, clean, err := database.Connection(dbConf)
	if err != nil {
		return nil, nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		clean()
		return nil, nil, err
	}
	// create a new temporal client
	c, err := client.Dial(client.Options{
		HostPort: configs.Temporal.HostPort,
	})
	if err != nil {
		clean()
		return nil, nil, err
	}
	cleanup := func() {
		clean()
		c.Close()
	}
	log := logger.New(logrus.Fields{"app-name": "shopper"}, logrus.DebugLevel)
	productStore := data.NewProductStore(db)
	cartStore := data.NewCartStore(db)
	orderStore := data.NewOrderStore(db)

	shopAPI := api.NewShopAPI(c, productStore, cartStore, orderStore)
	sApi, err := api.New(configs, shopAPI)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	sConfig := &service.Config{
		Address:  fmt.Sprintf("%s:%d", configs.Services.Shopper.Host, configs.Services.Shopper.Port),
		LogLevel: configs.Services.Shopper.LogLevel,
	}
	shopperService := service.NewService(sConfig, log, sqlDB, sApi)
	err = shopperService.RunMigration(func(db *sql.DB, schemaFs fs.FS, schema string) error {
		if err := productStore.Migrate(); err != nil {
			return err
		}
		if err := cartStore.Migrate(); err != nil {
			return err
		}
		if err := orderStore.Migrate(); err != nil {
			return err
		}
		return nil
	}, nil, "")

	if err != nil {
		cleanup()
		return nil, nil, err
	}
	cleanFunc, err := worker.Start(c, api.TQShoppingCart, orderStore, cartStore)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	return shopperService, func() {
		cleanFunc()
		cleanup()
	}, nil
}
