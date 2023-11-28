package data

import (
	"context"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/vithubati/go-core/logger"
	"gorm.io/gorm"
)

type Order struct {
	ID             string `gorm:"primary_key"`
	WorkflowId     string
	Products       pq.StringArray `gorm:"type:text[]"`
	InvoiceStatus  string
	ShippingStatus string
	OrderStatus    string
}
type orderStore struct {
	db *gorm.DB
}

func NewOrderStore(db *gorm.DB) OrderStore {
	return &orderStore{db: db}
}

func (s orderStore) Migrate() error {
	err := s.db.AutoMigrate(&Order{})
	if err != nil {
		return err
	}
	return nil
}

// Persist - insert an order
func (s orderStore) Persist(ctx context.Context, order *Order) (*Order, error) {
	logger.From(ctx).Info("orderStore:Persist: creating an order")
	order.ID = uuid.New().String()

	records := s.db.Create(&order)
	if records.Error != nil {
		return nil, records.Error
	}

	return order, nil
}

// Get - get an order
func (s orderStore) Get(ctx context.Context, orderId string) (order *Order, err error) {
	logger.From(ctx).Info("orderStore:Get: getting an order")
	order = &Order{}
	records := s.db.Where("id = ?", orderId).First(&order)
	if records.Error != nil {
		return order, records.Error
	}

	return order, err
}

// UpdateInvoiceStatus - update order by id, invoiceStatus
func (s orderStore) UpdateInvoiceStatus(ctx context.Context, workflowId string, status string) (err error) {
	logger.From(ctx).Info("orderStore:UpdateInvoiceStatus: updating invoice status")
	records := s.db.Model(&Order{}).Where("workflow_id = ?", workflowId).Update("invoice_status", status)
	if records.Error != nil {
		return records.Error
	}

	return err
}

// UpdateShippingStatus - update order by id, invoiceStatus
func (s orderStore) UpdateShippingStatus(ctx context.Context, workflowId string, status string) (err error) {
	logger.From(ctx).Info("orderStore:UpdateShippingStatus: updating shipping status")
	records := s.db.Model(&Order{}).Where("workflow_id = ?", workflowId).Update("shipping_status", status)
	if records.Error != nil {
		return records.Error
	}

	return err
}

// UpdateStatus - update order by id, invoiceStatus
func (s orderStore) UpdateStatus(ctx context.Context, workflowId string, status string) (err error) {
	logger.From(ctx).Info("orderStore:UpdateStatus: updating status")
	records := s.db.Model(&Order{}).Where("workflow_id = ?", workflowId).Update("order_status", status)
	if records.Error != nil {
		return records.Error
	}

	return err
}
