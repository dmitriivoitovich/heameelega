package controller

import (
	"github.com/dmitriivoitovich/heameelega/helper"
	"github.com/labstack/echo/v4"
)

func GetHome(c echo.Context) error {
	// redirect authorised user to dashboard
	if c.(AppContext).User != nil {
		return redirect(c, helper.PageURLDashboard())
	}

	// redirect unauthorised user to login page
	return redirect(c, helper.PageURLLogin())
}
