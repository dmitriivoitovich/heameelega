package service

import (
	"github.com/dmitriivoitovich/wallester-test-assignment/controller/request"
	"github.com/dmitriivoitovich/wallester-test-assignment/dao"
	"github.com/dmitriivoitovich/wallester-test-assignment/dao/db"
	"github.com/google/uuid"
	"strings"
)

func CreateCustomer(req request.CreateCustomer) error {
	customer, err := convertReqToDB(req)
	if err != nil {
		return err
	}

	if err := dao.CreateCustomer(customer); err != nil {
		return err
	}

	return nil
}

func convertReqToDB(req request.CreateCustomer) (*db.Customer, error) {
	gender := db.GenderFemale
	if req.Gender == request.GenderMale {
		gender = db.GenderMale
	}

	customer := &db.Customer{
		ID:        uuid.New(),
		FirstName: strings.TrimSpace(req.FirstName),
		LastName:  strings.TrimSpace(req.LastName),
		Gender:    gender,
		Email:     strings.ToLower(req.Email),
	}

	date, err := req.BirthDateTime()
	if err != nil {
		return nil, err
	}

	customer.BirthDate = *date

	address := strings.TrimSpace(req.Address)
	if address != "" {
		customer.Address = &address
	}

	return customer, nil
}