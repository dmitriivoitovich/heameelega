package controller

import (
	"strconv"

	"github.com/dmitriivoitovich/heameelega/service"
	"github.com/dmitriivoitovich/heameelega/util/i18n"
	"github.com/labstack/echo/v4"
)

type dashboardTmplData struct {
	tmplData
	Stats StatsData
}

type StatsData struct {
	Keys   []i18n.Key
	Values []uint32
}

func GetDashboard(c echo.Context) error {
	ctx := c.(AppContext)

	tmplData := dashboardTmplData{
		tmplData: tmplData{
			User: *ctx.User,
			NavbarData: NavbarData{
				User:   *ctx.User,
				Search: "",
			},
			SidebarData: SidebarData{
				User:       *ctx.User,
				ActivePage: "dashboard",
			},
		},
		Stats: StatsData{},
	}

	// load registration statistics
	stats, appErr := service.DashboardStats(ctx.User.ID)
	if appErr != nil {
		return httpError(*appErr)
	}

	// convert statistics into data for the template
	for i := range stats {
		tmplData.Stats.Keys = append(tmplData.Stats.Keys, i18n.Key("column.month-"+strconv.Itoa(int(stats[i].Month))))
		tmplData.Stats.Values = append(tmplData.Stats.Values, stats[i].Registrations)
	}

	// render template
	return RenderTmpl(c, tmplDashboard, tmplData)
}
