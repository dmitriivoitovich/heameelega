package dao

import (
	"github.com/dmitriivoitovich/heameelega/dao/db"
	"github.com/google/uuid"
	"strings"
)

func CreateCustomer(customer *db.Customer) error {
	return db.DB.Create(customer).Error
}

func Customer(id uuid.UUID) (*db.Customer, error) {
	customer := &db.Customer{}

	err := db.DB.
		Where("ID = ?", id).
		First(customer).
		Error

	if err != nil {
		return nil, err
	}

	return customer, nil
}

func CustomerByEmail(email string) (*db.Customer, error) {
	customer := &db.Customer{}

	err := db.DB.
		Where("EMAIL = ?", email).
		First(customer).
		Error

	if err != nil {
		return nil, err
	}

	return customer, nil
}

func UpdateCustomer(customer db.Customer) error {
	return db.DB.
		Model(&db.Customer{}).
		Where("ID = ?", customer.ID).
		Updates(
			map[string]interface{}{
				"FirstName": customer.FirstName,
				"LastName":  customer.LastName,
				"BirthDate": customer.BirthDate,
				"Gender":    customer.Gender,
				"Email":     customer.Email,
				"Address":   customer.Address,
			},
		).
		Error
}

func Customers(page, pageSize uint32, order, direction string, filters ...string) ([]db.Customer, error) {
	customers := make([]db.Customer, 0)

	offset := (page - 1) * pageSize

	q := db.DB.Debug().
		Offset(int(offset)).
		Limit(int(pageSize))

	for i := range filters {
		filter := "%" + filters[i] + "%"
		q = q.Where("LOWER(first_name) LIKE LOWER(?) OR LOWER(last_name) LIKE LOWER(?)", filter, filter)
	}

	err := q.
		Where("deleted_at IS NULL").
		Order(order + " " + strings.ToUpper(direction)).
		Find(&customers).
		Error
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func CustomersCount(filters ...string) (uint32, error) {
	var count int64
	q := db.DB.Model(db.Customer{}).Debug()

	for i := range filters {
		filter := "%" + filters[i] + "%"
		q = q.Where("LOWER(first_name) LIKE LOWER(?) OR LOWER(last_name) LIKE LOWER(?)", filter, filter)
	}

	err := q.
		Where("deleted_at IS NULL").
		Count(&count).
		Error
	if err != nil {
		return 0, err
	}

	return uint32(count), nil
}
