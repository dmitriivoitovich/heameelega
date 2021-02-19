package controller

import (
	"html/template"

	"github.com/labstack/echo/v4"
)

var (
	dashboardTmpl = template.Must(
		template.New("dashboard.gohtml").
			Funcs(funcMap).
			ParseFiles(
				"template/dashboard/dashboard.gohtml",
				"template/layout/header.gohtml",
				"template/layout/footer.gohtml",
				"template/layout/navbar.gohtml",
				"template/layout/sidebar.gohtml",
			),
	)
)

func GetDashboard(c echo.Context) error {
	return dashboardTmpl.Execute(c.Response(), nil)
}
