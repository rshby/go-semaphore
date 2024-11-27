package repository

import (
	"context"
	"github.com/sirupsen/logrus"
	"go-semaphore/internal/entity"
	"go-semaphore/utils"
	"gorm.io/gorm"
	"time"
)

type customerRepository struct {
	db *gorm.DB
}

// NewCustomerRepository is function to create new instance customerRepository, implement from interface CustomerRepository
func NewCustomerRepository(db *gorm.DB) entity.CustomerRepository {
	return &customerRepository{
		db: db,
	}
}

// Create is method to create new customer
func (c *customerRepository) Create(ctx context.Context, input *entity.Customer) (*entity.Customer, error) {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": utils.DumpIncomingContext(ctx),
		"input":   utils.Dump(input),
	})

	// insert to database
	if err := c.db.Model(&entity.Customer{}).Create(input).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()
	return input, nil
}

// GetByID is method to get data customer by ID
func (c *customerRepository) GetByID(ctx context.Context, id uint64) (*entity.Customer, error) {
	//TODO implement me
	panic("implement me")
}
