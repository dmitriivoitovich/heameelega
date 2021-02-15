package controller

import (
	"github.com/dmitriivoitovich/wallester-test-assignment/controller/request"
	"github.com/dmitriivoitovich/wallester-test-assignment/service"
	"github.com/labstack/echo/v4"
	"html/template"
	"net/http"
)

type TmplData struct {
	Customer      request.CreateCustomer
	InvalidFields []string
	Error         string
}

var createCustomerTmpl = template.Must(
	template.New("create.gohtml").Funcs(funcMap).ParseFiles("templates/create.gohtml"),
)

var funcMap = template.FuncMap{
	"inSlice": func(slice []string, key string) bool {
		for i := range slice {
			if slice[i] == key {
				return true
			}
		}

		return false
	},
}

func GetCreateCustomer(c echo.Context) error {
	return createCustomerTmpl.Execute(c.Response(), TmplData{})
}

func PostCreateCustomer(c echo.Context) error {
	tmplData := &TmplData{}

	if err := c.Bind(&tmplData.Customer); err != nil {
		c.Logger().Error(err)
		tmplData.Error = "Something is wrong with your request"

		return createCustomerTmpl.Execute(c.Response(), tmplData)
	}

	reqValid, invalidFields := tmplData.Customer.Validate()
	tmplData.InvalidFields = invalidFields

	if reqValid {
		if err := service.CreateCustomer(tmplData.Customer); err != nil {
			c.Logger().Error(err)
			tmplData.Error = "Something went wrong"

			return createCustomerTmpl.Execute(c.Response(), tmplData)
		}

		return c.Redirect(http.StatusTemporaryRedirect, "/customers")
	}

	return createCustomerTmpl.Execute(c.Response(), tmplData)
}

//func ListCustomers(c echo.Context) error {
//	return nil
//}
//
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
