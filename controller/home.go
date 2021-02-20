package controller

import (
	"net/http"

	"github.com/dmitriivoitovich/heameelega/helper"
	"github.com/labstack/echo/v4"
)

func GetHome(c echo.Context) error {
	if c.(AppContext).User != nil {
		return c.Redirect(http.StatusSeeOther, helper.PageURLDashboard())
	}

	return c.Redirect(http.StatusSeeOther, helper.PageURLLogin())
}
