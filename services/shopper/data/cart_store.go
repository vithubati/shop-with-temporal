package data

import (
	"context"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/vithubati/go-core/logger"
	"gorm.io/gorm"
)

type Cart struct {
	ID       string         `gorm:"primary_key"`
	Products pq.StringArray `gorm:"type:text[]"`
}

type cartStore struct {
	db *gorm.DB
}

func NewCartStore(db *gorm.DB) CartStore {
	return &cartStore{db: db}
}

func (s cartStore) Migrate() error {
	err := s.db.AutoMigrate(&Cart{})
	if err != nil {
		return err
	}
	return nil
}

func (s cartStore) Persist(ctx context.Context, cart *Cart) (*Cart, error) {
	logger.From(ctx).Info("cartStore:Persist: inserting cart")
	cart.ID = uuid.NewString()
	records := s.db.Model(&Cart{}).Create(&cart)
	if records.Error != nil {
		logger.From(ctx).WithError(records.Error).Error("inserting cart")
		return nil, records.Error
	}
	return cart, nil
}

func (s cartStore) Get(ctx context.Context) (*Cart, error) {
	logger.From(ctx).Info("cartStore:List: getting cart")
	cart := &Cart{}
	// get all records and assign to data
	records := s.db.Last(&cart)
	if records.Error != nil {
		logger.From(ctx).WithError(records.Error).Error("getting last cart")
		return nil, records.Error
	}

	return cart, nil
}
