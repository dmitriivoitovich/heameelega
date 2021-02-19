package controller

import (
	"net/http"

	"github.com/dmitriivoitovich/heameelega/helper"
	"github.com/labstack/echo/v4"
)

func GetHome(c echo.Context) error {
	return c.Redirect(http.StatusSeeOther, helper.PageURLDashboard())
}
