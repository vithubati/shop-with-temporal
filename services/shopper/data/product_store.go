package data

import (
	"context"
	"github.com/vithubati/go-core/logger"
	"gorm.io/gorm"
)

type Product struct {
	Name  string
	Price float64
}
type productStore struct {
	db *gorm.DB
}

func NewProductStore(db *gorm.DB) ProductStore {
	return &productStore{db: db}
}

func (s productStore) List(ctx context.Context) ([]Product, error) {
	logger.From(ctx).Info("productStore:List: listing products")
	data := make([]Product, 0)
	records := s.db.Find(&data)
	if records.Error != nil {
		logger.From(ctx).WithError(records.Error).Error("listing products")
		return nil, records.Error
	}
	return data, nil
}

func (s productStore) Migrate() error {
	err := s.db.AutoMigrate(&Product{})
	if err != nil {
		return err
	}

	var count int64
	s.db.Model(&Product{}).Count(&count)
	if count == 0 {
		tx := s.db.Create([]Product{
			{Name: "VW", Price: 10.0},
			{Name: "Honda", Price: 20.0},
			{Name: "Lexus", Price: 30.0},
			{Name: "Toyota", Price: 30.0},
		})
		if tx.Error != nil {
			return tx.Error
		}
	}
	return nil
}
