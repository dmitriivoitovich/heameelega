package service

import (
	"github.com/dmitriivoitovich/heameelega/controller/request"
	"github.com/dmitriivoitovich/heameelega/dao"
	"github.com/dmitriivoitovich/heameelega/dao/db"
	"github.com/google/uuid"
	"math"
	"regexp"
	"strings"
)

const pageSize = 10

var (
	spacesRegExp = regexp.MustCompile(`\s+`)
)

func SearchCustomers(req request.SearchCustomersNormalized) ([]db.Customer, uint32, error) {
	s := spacesRegExp.ReplaceAllString(req.Filter, " ")

	filters := make([]string, 0)
	if s != "" {
		filters = strings.Split(s, " ")
	}

	customers, err := dao.Customers(req.Page, pageSize, req.Order, req.Direction, filters...)
	if err != nil {
		return nil, 0, err
	}

	total, err := dao.CustomersCount(filters...)
	if err != nil {
		return nil, 0, err
	}

	pages := uint32(math.Ceil(float64(total) / float64(pageSize)))

	return customers, pages, nil
}

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
