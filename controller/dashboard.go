package controller

import (
	"github.com/labstack/echo/v4"
)

func GetDashboard(c echo.Context) error {
	ctx := c.(AppContext)

	tmplData := tmplData{
		User: *ctx.User,
		NavbarData: NavbarData{
			User:   *ctx.User,
			Search: "",
		},
		SidebarData: SidebarData{
			User:       *ctx.User,
			ActivePage: "dashboard",
		},
	}

	return RenderTmpl(c, tmplDashboard, tmplData)
}
