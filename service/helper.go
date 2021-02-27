package service

import (
	"github.com/dmitriivoitovich/heameelega/controller/request"
	"github.com/dmitriivoitovich/heameelega/dao/db"
	"github.com/dmitriivoitovich/heameelega/util/apperror"
	"github.com/google/uuid"
)

func convertCreateReqToCustomer(userID uuid.UUID, req request.CreateCustomer) (*db.Customer, *apperror.Error) {
	gender := db.GenderMale
	if req.Gender == request.GenderFemale {
		gender = db.GenderFemale
	}

	customer := &db.Customer{
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
