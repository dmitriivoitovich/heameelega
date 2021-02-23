package controller

import "github.com/labstack/echo/v4"

func GetDashboard(c echo.Context) error {
	return RenderTmpl(c, tmplDashboard, nil)
}
