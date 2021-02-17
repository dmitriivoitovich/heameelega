package controller

import (
	"fmt"
	"github.com/dmitriivoitovich/heameelega/controller/request"
	"github.com/dmitriivoitovich/heameelega/dao"
	"github.com/dmitriivoitovich/heameelega/dao/db"
	"github.com/dmitriivoitovich/heameelega/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
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

var (
	funcMap = template.FuncMap{
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
		"link": func(filter, order, direction string, reverseDirection bool, page uint32) string {
			if reverseDirection {
				if direction == "asc" {
					direction = "desc"
				} else {
					direction = "asc"
				}
			}

			params := url.Values{
				"o": {order},
				"d": {direction},
				"p": {strconv.Itoa(int(page))},
			}

			if filter != "" {
				params.Add("s", filter)
			}

			return fmt.Sprintf("/customers?%s", params.Encode())
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
	Error500Tmpl = template.Must(
		template.New("500.gohtml").
			ParseFiles(
				"template/error/500.gohtml",
				"template/layout/header.gohtml",
				"template/layout/footer.gohtml",
			),
	)
	Error400Tmpl = template.Must(
		template.New("404.gohtml").
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

		return c.Redirect(http.StatusSeeOther, "/customers")
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
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "failed to find a customer by id")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load customer")
	}

	return viewCustomerTmpl.Execute(c.Response(), customer)
}
