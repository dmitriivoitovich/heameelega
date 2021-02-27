package service

import (
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dmitriivoitovich/heameelega/controller/request"
	"github.com/dmitriivoitovich/heameelega/dao"
	"github.com/dmitriivoitovich/heameelega/dao/db"
	"github.com/dmitriivoitovich/heameelega/util/apperror"
	"github.com/dmitriivoitovich/heameelega/util/i18n"
	"github.com/google/uuid"
)

const (
	PaginationPageSize   = 25
	PaginationPagesLimit = 10
)

var spacesRegExp = regexp.MustCompile(`\s+`)

func SearchCustomers(userID uuid.UUID, req request.SearchCustomersNormalized) ([]db.Customer, uint32, *apperror.Error) {
	s := spacesRegExp.ReplaceAllString(req.Filter, " ")

	filters := make([]string, 0)
	if s != "" {
		filters = strings.Split(s, " ")
	}

	customers, err := dao.Customers(userID, req.Page, PaginationPageSize, req.Order, req.Direction, filters...)
	if err != nil {
		return nil, 0, apperror.Internal(err, "failed to load customers by search request")
	}

	total, err := dao.CustomersCount(userID, filters...)
	if err != nil {
		return nil, 0, apperror.Internal(err, "failed to count customers by search request")
	}

	pages := uint32(math.Ceil(float64(total) / float64(PaginationPageSize)))

	return customers, pages, nil
}

func CreateCustomer(userID uuid.UUID, req request.CreateCustomer) (*db.Customer, *apperror.Error) {
	customer, err := convertCreateReqToCustomer(userID, req)
	if err != nil {
		return nil, err
	}

	if err := dao.CreateCustomer(customer); err != nil {
		if dao.IsErrDuplicateKey(err) {
			return nil, apperror.Validation(i18n.KeyErrorEmailTaken)
		}

		return nil, apperror.Internal(err, "failed to create new customer")
	}

	return customer, nil
}

func ViewCustomer(userID, customerID uuid.UUID) (*db.Customer, *apperror.Error) {
	customer, err := dao.CustomerByIDAndUserID(customerID, userID)
	if err != nil {
		if dao.IsErrRecordNotFound(err) {
			return nil, apperror.NotFound("customer not found by provided id")
		}

		return nil, apperror.Internal(err, "failed to load customer")
	}

	return customer, nil
}

func UpdateCustomer(userID uuid.UUID, req request.EditCustomer) *apperror.Error {
	customer, appErr := ViewCustomer(userID, req.ID)
	if appErr != nil {
		return appErr
	}

	if customer.UpdatedAt.After(req.LoadedAt) {
		return apperror.Validation(i18n.KeyErrorDataCollision)
	}

	customer.FirstName = req.FirstName
	customer.LastName = req.LastName

	birthDate, err := req.BirthDateTime()
	if err != nil {
		return apperror.Internal(err, "failed to parse birth date from request")
	}

	customer.BirthDate = birthDate
	customer.Gender = req.Gender == request.GenderMale
	customer.Email = req.Email

	customer.Address = nil
	if req.Address != "" {
		customer.Address = &req.Address
	}

	if err := dao.UpdateCustomer(customer); err != nil {
		if dao.IsErrDuplicateKey(err) {
			return apperror.Validation(i18n.KeyErrorEmailTaken)
		}

		return apperror.Internal(err, "failed to update customer")
	}

	return nil
}

func ConvertCustomerToEditReq(customer db.Customer) *request.EditCustomer {
	req := &request.EditCustomer{
		ID:        customer.ID,
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		BirthDate: customer.BirthDate.Format(request.DateFormat),
		Gender:    request.GenderMale,
		Email:     customer.Email,
		Address:   "",
		LoadedAt:  time.Now(),
	}

	if !customer.Gender {
		req.Gender = request.GenderFemale
	}

	if customer.Address != nil {
		req.Address = *customer.Address
	}

	return req
}

func CreateFakeCustomers(userID uuid.UUID, amount uint32) *apperror.Error {
	for i := uint32(0); i < amount; i++ {
		now := time.Now()
		seed := now.Nanosecond()
		faker := gofakeit.New(int64(seed))

		gender := db.GenderMale
		if faker.Gender() == "female" {
			gender = db.GenderFemale
		}

		address := faker.Address().Address
		birthDate := faker.DateRange(now.AddDate(-60, 0, 0), now.AddDate(-18, 0, 0))
		registerDate := faker.DateRange(now.AddDate(-1, 0, 0), now)

		customer := &db.Customer{
			BaseModel: db.BaseModel{
				CreatedAt: registerDate,
				UpdatedAt: registerDate,
			},
			FirstName: faker.FirstName(),
			LastName:  faker.LastName(),
			BirthDate: birthDate,
			Gender:    gender,
			Email:     faker.Email(),
			Address:   &address,
			UserID:    userID,
		}

		if err := dao.CreateCustomer(customer); err != nil {
			if dao.IsErrDuplicateKey(err) {
				// ignore duplicate key error
				continue
			}

			return apperror.Internal(err, "failed to create fake customer")
		}
	}

	return nil
}

func DashboardStats(userID uuid.UUID) ([]dao.CustomerMonthlyRegistrations, *apperror.Error) {
	end := time.Now()
	currentYear, currentMonth, _ := end.Date()
	start := time.Date(currentYear-1, currentMonth+1, 1, 0, 0, 0, 0, end.Location())

	stats, err := dao.CustomersMonthlyRegistrations(userID, start, end)
	if err != nil {
		return nil, apperror.Internal(err, "failed to load customer registrations data")
	}

	return stats, nil
}

func DeleteCustomer(userID, customerID uuid.UUID) *apperror.Error {
	// load customer
	customer, appErr := ViewCustomer(userID, customerID)
	if appErr != nil {
		return appErr
	}

	// delete customer
	if err := dao.DeleteCustomer(customer); err != nil {
		return apperror.Internal(err, "failed to delete customer")
	}

	return nil
}
