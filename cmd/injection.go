package cmd

import (
	"go-semaphore/internal/repository"
	"go-semaphore/internal/service"
	"go-semaphore/proto/pb"
	"gorm.io/gorm"
)

func InitCustomerService(db *gorm.DB) pb.CustomerServiceServer {
	// register repository
	customerRepository := repository.NewCustomerRepository(db)
	return service.NewCustomerService(customerRepository)
}
