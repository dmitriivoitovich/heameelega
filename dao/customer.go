package dao

import (
	"github.com/dmitriivoitovich/wallester-test-assignment/dao/db"
	"github.com/google/uuid"
)

func CreateCustomer(customer *db.Customer) error {
	return db.DB.Create(customer).Error
}

func Customer(id uuid.UUID) (*db.Customer, error) {
	customer := &db.Customer{}

	err := db.DB.Where("ID = ?", id).
		First(customer).
		Error

	if err != nil {
		return nil, err
	}

	return customer, nil
}
