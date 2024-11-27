package entity

import (
	"context"
	"go-semaphore/proto/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Customer struct {
	ID          uint64    `gorm:"column:id;type:bigint;not null;primaryKey;autoIncrement" json:"id"`
	FirstName   string    `gorm:"column:first_name;type:varchar(255);not null" json:"first_name"`
	LastName    string    `gorm:"column:last_name;type:varchar(255)" json:"last_name"`
	Email       string    `gorm:"column:email;type:varchar(255);not nul;unique" json:"email"`
	PhoneNumber string    `gorm:"column:phone_number;type:varchar(15);unique" json:"phone_number"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;not null;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamp;not null;autoCreateTime;autoUpdateTime" json:"updated_at"`
}

func (c *Customer) TableName() string {
	return "customers"
}

func (c *Customer) IsIDExists() bool {
	return c.ID > 0
}

func (c *Customer) IsFirstNameExists() bool {
	return c.FirstName != ""
}

func (c *Customer) IsLastNameExists() bool {
	return c.LastName != ""
}

func (c *Customer) IsEmailExists() bool {
	return c.Email != ""
}

func (c *Customer) IsPhoneNumberExists() bool {
	return c.PhoneNumber != ""
}

func (c *Customer) ToProto() *pb.Customer {
	return &pb.Customer{
		Id:          c.ID,
		FirstName:   c.FirstName,
		LastName:    c.LastName,
		Email:       c.Email,
		PhoneNumber: c.PhoneNumber,
		CreatedAt:   timestamppb.New(c.CreatedAt),
		UpdatedAt:   timestamppb.New(c.UpdatedAt),
	}
}

type CustomerRepository interface {
	Create(ctx context.Context, input *Customer) (*Customer, error)
	GetByID(ctx context.Context, id uint64) (*Customer, error)
}
