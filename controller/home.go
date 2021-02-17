package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetHome(c echo.Context) error {
	return c.Redirect(http.StatusSeeOther, "/dashboard")
}
