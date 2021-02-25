package service

import (
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/dmitriivoitovich/heameelega/controller/request"
	"github.com/dmitriivoitovich/heameelega/dao"
	"github.com/dmitriivoitovich/heameelega/dao/db"
	"github.com/dmitriivoitovich/heameelega/util/apperror"
	"github.com/dmitriivoitovich/heameelega/util/i18n"
	"github.com/google/uuid"
)

const pageSize = 10

var spacesRegExp = regexp.MustCompile(`\s+`)

func SearchCustomers(userID uuid.UUID, req request.SearchCustomersNormalized) ([]db.Customer, uint32, *apperror.Error) {
	s := spacesRegExp.ReplaceAllString(req.Filter, " ")

	filters := make([]string, 0)
	if s != "" {
		filters = strings.Split(s, " ")
	}

	customers, err := dao.Customers(userID, req.Page, pageSize, req.Order, req.Direction, filters...)
	if err != nil {
		return nil, 0, apperror.Internal(err, "failed to load customers by search request")
	}

	total, err := dao.CustomersCount(userID, filters...)
	if err != nil {
		return nil, 0, apperror.Internal(err, "failed to count customers by search request")
	}

	pages := uint32(math.Ceil(float64(total) / float64(pageSize)))

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

func convertCreateReqToCustomer(userID uuid.UUID, req request.CreateCustomer) (*db.Customer, *apperror.Error) {
	gender := db.GenderMale
	if req.Gender == request.GenderFemale {
		gender = db.GenderFemale
	}

	customer := &db.Customer{
		ID:        uuid.New(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Gender:    gender,
		Email:     req.Email,
		UserID:    userID,
	}

	date, err := req.BirthDateTime()
	if err != nil {
		return nil, apperror.BadRequest(err, "failed to parse birth date from request")
	}

	customer.BirthDate = *date

	if req.Address != "" {
		customer.Address = &req.Address
	}

	return customer, nil
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
