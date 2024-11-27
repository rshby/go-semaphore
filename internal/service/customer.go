package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"go-semaphore/internal/entity"
	"go-semaphore/proto/pb"
	"go-semaphore/utils"
)

type customerService struct {
	customerRepository entity.CustomerRepository
	pb.UnimplementedCustomerServiceServer
}

// NewCustomerService is method to create instance
func NewCustomerService(customerRepository entity.CustomerRepository) pb.CustomerServiceServer {
	return &customerService{
		customerRepository: customerRepository,
	}
}

func (c *customerService) CreateCustomer(ctx context.Context, request *pb.CreateCustomerRequestDTO) (*pb.Customer, error) {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": utils.DumpIncomingContext(ctx),
		"request": utils.Dump(request),
	})

	// call method in repository to insert data customer
	newCustomer, err := c.customerRepository.Create(ctx, &entity.Customer{
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		Email:       request.Email,
		PhoneNumber: request.PhoneNumber,
	})
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return newCustomer.ToProto(), nil
}

func (c *customerService) GetCustomerByID(ctx context.Context, request *pb.GetCustomerByIDRequestDTO) (*pb.Customer, error) {
	//TODO implement me
	panic("implement me")
}
