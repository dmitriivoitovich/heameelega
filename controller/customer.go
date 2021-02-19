package controller

import (
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/dmitriivoitovich/heameelega/controller/request"
	"github.com/dmitriivoitovich/heameelega/dao"
	"github.com/dmitriivoitovich/heameelega/dao/db"
	"github.com/dmitriivoitovich/heameelega/helper"
	"github.com/dmitriivoitovich/heameelega/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type createCustomerTmplData struct {
	Request       request.CreateCustomer
	InvalidFields []string
	Error         string
}

type searchCustomerTmplData struct {
	Request       request.SearchCustomers
	InvalidFields []string
	Customers     []db.Customer
	Pages         uint32
	Error         string
}

type editCustomerTmplData struct {
	Request       request.EditCustomer
	InvalidFields []string
	Error         string
}

var (
	funcMap = template.FuncMap{
		"linkHome":            helper.PageURLHome,
		"linkLogout":          helper.PageURLLogout,
		"linkDashboard":       helper.PageURLDashboard,
		"linkCustomers":       helper.PageURLCustomers,
		"linkSearchCustomers": helper.PageURLSearchCustomers,
		"linkNewCustomer":     helper.PageURLNewCustomer,
		"linkViewCustomer":    helper.PageURLViewCustomer,
		"linkEditCustomer":    helper.PageURLEditCustomer,
		"inSlice": func(slice []string, key string) bool {
			for i := range slice {
				if slice[i] == key {
					return true
				}
			}

			return false
		},
		"seq": func(items uint32) []uint32 {
			res := make([]uint32, 0, items)

			for i := uint32(0); i < items; i++ {
				res = append(res, i)
			}

			return res
		},
		"inc": func(item uint32) uint32 {
			return item + 1
		},
		"dec": func(item uint32) uint32 {
			return item - 1
		},
	}

	searchCustomerTmpl = template.Must(
		template.New("search.gohtml").
			Funcs(funcMap).
			ParseFiles(
				"template/customer/search.gohtml",
				"template/layout/header.gohtml",
				"template/layout/footer.gohtml",
				"template/layout/navbar.gohtml",
				"template/layout/sidebar.gohtml",
			),
	)
	createCustomerTmpl = template.Must(
		template.New("create.gohtml").
			Funcs(funcMap).
			ParseFiles(
				"template/customer/create.gohtml",
				"template/layout/header.gohtml",
				"template/layout/footer.gohtml",
				"template/layout/navbar.gohtml",
				"template/layout/sidebar.gohtml",
			),
	)
	viewCustomerTmpl = template.Must(
		template.New("view.gohtml").
			Funcs(funcMap).
			ParseFiles(
				"template/customer/view.gohtml",
				"template/layout/header.gohtml",
				"template/layout/footer.gohtml",
				"template/layout/navbar.gohtml",
				"template/layout/sidebar.gohtml",
			),
	)
	editCustomerTmpl = template.Must(
		template.New("edit.gohtml").
			Funcs(funcMap).
			ParseFiles(
				"template/customer/edit.gohtml",
				"template/layout/header.gohtml",
				"template/layout/footer.gohtml",
				"template/layout/navbar.gohtml",
				"template/layout/sidebar.gohtml",
			),
	)
	Error500Tmpl = template.Must(
		template.New("500.gohtml").
			Funcs(funcMap).
			ParseFiles(
				"template/error/500.gohtml",
				"template/layout/header.gohtml",
				"template/layout/footer.gohtml",
			),
	)
	Error400Tmpl = template.Must(
		template.New("404.gohtml").
			Funcs(funcMap).
			ParseFiles(
				"template/error/404.gohtml",
				"template/layout/header.gohtml",
				"template/layout/footer.gohtml",
			),
	)
)

func GetCreateCustomer(c echo.Context) error {
	return createCustomerTmpl.Execute(c.Response(), createCustomerTmplData{})
}

func PostCreateCustomer(c echo.Context) error {
	tmplData := &createCustomerTmplData{}

	if err := c.Bind(&tmplData.Request); err != nil {
		c.Logger().Error(err)

		tmplData.Error = "Something is wrong with your request"

		return createCustomerTmpl.Execute(c.Response(), tmplData)
	}

	invalidFields := tmplData.Request.Validate()
	tmplData.InvalidFields = invalidFields

	if len(invalidFields) == 0 {
		if err := service.CreateCustomer(tmplData.Request); err != nil {
			c.Logger().Error(err)

			tmplData.Error = "Something went wrong"

			return createCustomerTmpl.Execute(c.Response(), tmplData)
		}

		return c.Redirect(http.StatusSeeOther, helper.PageURLCustomers())
	}

	return createCustomerTmpl.Execute(c.Response(), tmplData)
}

func GetSearchCustomers(c echo.Context) error {
	tmplData := &searchCustomerTmplData{}

	if err := c.Bind(&tmplData.Request); err != nil {
		c.Logger().Error(err)

		tmplData.Error = "Something is wrong with your request"
	}

	invalidFields := tmplData.Request.Validate()
	if len(invalidFields) > 0 {
		tmplData.Error = "Something is wrong with your request"
	}

	customers, pages, err := service.SearchCustomers(tmplData.Request.Normalized())
	if err != nil {
		c.Logger().Error(err)

		tmplData.Error = "Something went wrong"

		return searchCustomerTmpl.Execute(c.Response(), tmplData)
	}

	tmplData.Customers = customers
	tmplData.Pages = pages

	return searchCustomerTmpl.Execute(c.Response(), tmplData)
}

func GetViewCustomer(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "failed to parse customer id from request")
	}

	customer, err := dao.Customer(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "failed to find a customer by id")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load customer")
	}

	return viewCustomerTmpl.Execute(c.Response(), customer)
}

func GetEditCustomer(c echo.Context) error {
	tmplData := &editCustomerTmplData{}
	if err := c.Bind(&tmplData.Request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to bind request")
	}

	customer, err := dao.Customer(tmplData.Request.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "failed to find customer by id")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load customer")
	}

	tmplData.Request.FirstName = customer.FirstName
	tmplData.Request.LastName = customer.LastName
	tmplData.Request.BirthDate = customer.BirthDate.Format("2006-01-02")

	if customer.Gender {
		tmplData.Request.Gender = "male"
	} else {
		tmplData.Request.Gender = "female"
	}

	tmplData.Request.Email = customer.Email
	if customer.Address != nil {
		tmplData.Request.Address = *customer.Address
	}

	tmplData.Request.LoadedAt = time.Now()

	return editCustomerTmpl.Execute(c.Response(), tmplData)
}

func PostEditCustomer(c echo.Context) error {
	tmplData := &editCustomerTmplData{}
	if err := c.Bind(&tmplData.Request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to bind request")
	}

	tmplData.InvalidFields = tmplData.Request.Validate()
	if len(tmplData.InvalidFields) > 0 {
		return editCustomerTmpl.Execute(c.Response(), tmplData)
	}

	customer, err := dao.Customer(tmplData.Request.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "failed to find customer by id")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load customer")
	}

	if err := service.UpdateCustomer(tmplData.Request, *customer); err != nil {
		c.Logger().Error(err)

		tmplData.Error = "Something went wrong"

		return editCustomerTmpl.Execute(c.Response(), tmplData)
	}

	return c.Redirect(http.StatusSeeOther, helper.PageURLViewCustomer(customer.ID))
}
