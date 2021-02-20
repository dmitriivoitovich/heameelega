package controller

import (
	"html/template"

	"github.com/labstack/echo/v4"
)

var (
	loginTmpl = template.Must(
		template.New("login.gohtml").
			Funcs(funcMap).
			ParseFiles(
				"template/user/login.gohtml",
				"template/layout/header.gohtml",
				"template/layout/footer.gohtml",
			),
	)
	registerTmpl = template.Must(
		template.New("register.gohtml").
			Funcs(funcMap).
			ParseFiles(
				"template/user/register.gohtml",
				"template/layout/header.gohtml",
				"template/layout/footer.gohtml",
			),
	)
)

func GetLogin(c echo.Context) error {
	return loginTmpl.Execute(c.Response(), nil)
}

func GetRegister(c echo.Context) error {
	return registerTmpl.Execute(c.Response(), nil)
}
