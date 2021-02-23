package controller

import (
	"github.com/dmitriivoitovich/heameelega/controller/request"
	"github.com/dmitriivoitovich/heameelega/dao/db"
	"github.com/dmitriivoitovich/heameelega/helper"
	"github.com/dmitriivoitovich/heameelega/service"
	"github.com/dmitriivoitovich/heameelega/util/apperror"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type createCustomerTmplData struct {
	Request       request.CreateCustomer
	InvalidFields []string
	Error         *apperror.Error
}

type viewCustomerTmplData struct {
	Customer db.Customer
}

type editCustomerTmplData struct {
	Request       request.EditCustomer
	InvalidFields []string
	Error         *apperror.Error
}

type searchCustomerTmplData struct {
	Request       request.SearchCustomers
	InvalidFields []string
	Customers     []db.Customer
	Pages         uint32
	Error         *apperror.Error
}

func GetCreateCustomer(c echo.Context) error {
	return RenderTmpl(c, tmplCreateCustomer, createCustomerTmplData{})
}

func PostCreateCustomer(c echo.Context) error {
	ctx := c.(AppContext)

	// bind request
	tmplData := &createCustomerTmplData{}
	if err := c.Bind(&tmplData.Request); err != nil {
		return httpError(*apperror.BadRequest(err, "failed to bind create customer request"))
	}

	// validate request
	tmplData.InvalidFields = tmplData.Request.Validate()
	if len(tmplData.InvalidFields) > 0 {
		return RenderTmpl(c, tmplCreateCustomer, tmplData)
	}

	// create customer
	customer, appErr := service.CreateCustomer(ctx.User.ID, tmplData.Request.Sanitized())
	if appErr != nil {
		if appErr.IsValidation() {
			tmplData.Error = appErr

			return RenderTmpl(c, tmplCreateCustomer, tmplData)
		}

		return httpError(*appErr)
	}

	// redirect user to view customer page
	return redirect(c, helper.PageURLViewCustomer(customer.ID))
}

func GetSearchCustomers(c echo.Context) error {
	ctx := c.(AppContext)

	// bind request
	tmplData := &searchCustomerTmplData{}
	if err := c.Bind(&tmplData.Request); err != nil {
		return httpError(*apperror.BadRequest(err, "failed to bind search customers request"))
	}

	// validate request
	tmplData.InvalidFields = tmplData.Request.Validate()
	if len(tmplData.InvalidFields) > 0 {
		return RenderTmpl(c, tmplSearchCustomers, tmplData)
	}

	// search customers
	customers, pages, appErr := service.SearchCustomers(ctx.User.ID, tmplData.Request.Normalized())
	if appErr != nil {
		return httpError(*appErr)
	}

	tmplData.Customers = customers
	tmplData.Pages = pages

	// render template
	return RenderTmpl(c, tmplSearchCustomers, tmplData)
}

func GetViewCustomer(c echo.Context) error {
	ctx := c.(AppContext)

	// parse customer id from request
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpError(*apperror.BadRequest(err, "failed to parse customer if from request"))
	}

	// load customer
	customer, appErr := service.ViewCustomer(ctx.User.ID, id)
	if appErr != nil {
		return httpError(*appErr)
	}

	// render template
	return RenderTmpl(c, tmplViewCustomer, viewCustomerTmplData{Customer: *customer})
}

func GetEditCustomer(c echo.Context) error {
	ctx := c.(AppContext)

	// bind request
	tmplData := &editCustomerTmplData{}
	if err := c.Bind(&tmplData.Request); err != nil {
		return httpError(*apperror.BadRequest(err, "failed to bind edit customer request"))
	}

	// load customer
	customer, appErr := service.ViewCustomer(ctx.User.ID, tmplData.Request.ID)
	if appErr != nil {
		return httpError(*appErr)
	}

	// convert customer to edit customer request
	tmplData.Request = *service.ConvertCustomerToEditReq(*customer)

	// render template
	return RenderTmpl(c, tmplEditCustomer, tmplData)
}

func PostEditCustomer(c echo.Context) error {
	ctx := c.(AppContext)

	// bind request
	tmplData := &editCustomerTmplData{}
	if err := c.Bind(&tmplData.Request); err != nil {
		return httpError(*apperror.BadRequest(err, "failed to bind edit customer request"))
	}

	// validate request
	tmplData.InvalidFields = tmplData.Request.Validate()
	if len(tmplData.InvalidFields) > 0 {
		return RenderTmpl(c, tmplEditCustomer, tmplData)
	}

	// update customer
	if appErr := service.UpdateCustomer(ctx.User.ID, tmplData.Request.Sanitized()); appErr != nil {
		if appErr.IsValidation() {
			tmplData.Error = appErr

			return RenderTmpl(c, tmplEditCustomer, tmplData)
		}

		return httpError(*appErr)
	}

	// redirect user to view customer page
	return redirect(c, helper.PageURLViewCustomer(tmplData.Request.ID))
}
