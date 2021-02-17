package controller

import (
	"fmt"
	"github.com/dmitriivoitovich/wallester-test-assignment/controller/request"
	"github.com/dmitriivoitovich/wallester-test-assignment/dao/db"
	"github.com/dmitriivoitovich/wallester-test-assignment/service"
	"github.com/labstack/echo/v4"
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
				"templates/search.gohtml",
				"templates/layout/header.gohtml",
				"templates/layout/footer.gohtml",
				"templates/layout/navbar.gohtml",
				"templates/layout/sidebar.gohtml",
			),
	)
	createCustomerTmpl = template.Must(
		template.New("create.gohtml").
			Funcs(funcMap).
			ParseFiles(
				"templates/create.gohtml",
				"templates/layout/header.gohtml",
				"templates/layout/footer.gohtml",
				"templates/layout/navbar.gohtml",
				"templates/layout/sidebar.gohtml",
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

//func EditCustomer(c echo.Context) error {
//
// http.NotFound(w, r)
// return
//
//	id, err := uuid.Parse(c.Param("id"))
//	if err != nil {
//		return err
//	}
//
//	customer, err := dao.Customer(id)
//	if err != nil {
//		return err
//	}
//
//	t, _ := template.ParseFiles("edit.gohtml")
//
//	return t.Execute(c.Response(), customer)
//}
