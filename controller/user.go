package controller

import (
	"github.com/dmitriivoitovich/heameelega/controller/request"
	"github.com/dmitriivoitovich/heameelega/helper"
	"github.com/dmitriivoitovich/heameelega/service"
	"github.com/dmitriivoitovich/heameelega/util/apperror"
	"github.com/labstack/echo/v4"
)

type loginUserTmplData struct {
	Request       request.LoginUser
	InvalidFields []string
	Error         *apperror.Error
}

type registerUserTmplData struct {
	Request       request.RegisterUser
	InvalidFields []string
	Error         *apperror.Error
}

type userSettingsTmplData struct {
	tmplData
	Request       request.EditUser
	InvalidFields []string
	Error         *apperror.Error
}

func GetLogin(c echo.Context) error {
	return RenderTmpl(c, tmplLoginUser, loginUserTmplData{})
}

func PostLogin(c echo.Context) error {
	// bind request
	tmplData := &loginUserTmplData{}
	if err := c.Bind(&tmplData.Request); err != nil {
		return httpError(*apperror.BadRequest(err, "failed to bind login user request"))
	}

	// validate request
	tmplData.InvalidFields = tmplData.Request.Validate()
	if len(tmplData.InvalidFields) > 0 {
		return RenderTmpl(c, tmplLoginUser, tmplData)
	}

	// check request and update user's access token if successful
	user, appErr := service.UserLogin(tmplData.Request)
	if appErr != nil {
		if appErr.IsValidation() {
			tmplData.Error = appErr

			return RenderTmpl(c, tmplLoginUser, tmplData)
		}

		return httpError(*appErr)
	}

	// update session cookie
	setSessionCookie(c, *user.AccessToken)

	return redirect(c, helper.PageURLDashboard())
}

func GetRegister(c echo.Context) error {
	return RenderTmpl(c, tmplRegisterUser, registerUserTmplData{})
}

func PostRegister(c echo.Context) error {
	// bind request
	tmplData := &registerUserTmplData{}
	if err := c.Bind(&tmplData.Request); err != nil {
		return httpError(*apperror.BadRequest(err, "failed to bind register user request"))
	}

	// validate request
	tmplData.InvalidFields = tmplData.Request.Validate()
	if len(tmplData.InvalidFields) > 0 {
		return RenderTmpl(c, tmplRegisterUser, tmplData)
	}

	// create a new user
	user, appErr := service.UserRegister(tmplData.Request.Sanitized())
	if appErr != nil {
		if appErr.IsValidation() {
			tmplData.Error = appErr

			return RenderTmpl(c, tmplRegisterUser, tmplData)
		}

		return httpError(*appErr)
	}

	// update session cookie
	setSessionCookie(c, *user.AccessToken)

	return redirect(c, helper.PageURLDashboard())
}

func PostLogout(c echo.Context) error {
	ctx := c.(AppContext)

	// remove user's access token
	if appErr := service.UserLogout(ctx.User); appErr != nil {
		return httpError(*appErr)
	}

	return redirect(c, helper.PageURLLogin())
}

func GetSettings(c echo.Context) error {
	ctx := c.(AppContext)

	tmplData := userSettingsTmplData{
		tmplData: tmplData{
			User: *ctx.User,
			NavbarData: NavbarData{
				User:   *ctx.User,
				Search: "",
			},
			SidebarData: SidebarData{
				User:       *ctx.User,
				ActivePage: "",
			},
		},
		Request: request.EditUser{
			Email:    ctx.User.Email,
			Language: ctx.User.Language,
		},
	}

	return RenderTmpl(c, tmplUserSettings, tmplData)
}

func PostSettings(c echo.Context) error {
	ctx := c.(AppContext)

	tmplData := userSettingsTmplData{
		tmplData: tmplData{
			User: *ctx.User,
			NavbarData: NavbarData{
				User:   *ctx.User,
				Search: "",
			},
			SidebarData: SidebarData{
				User:       *ctx.User,
				ActivePage: "",
			},
		},
	}

	// bind request
	if err := c.Bind(&tmplData.Request); err != nil {
		return httpError(*apperror.BadRequest(err, "failed to bind edit user request"))
	}

	// validate request
	tmplData.InvalidFields = tmplData.Request.Validate()
	if len(tmplData.InvalidFields) > 0 {
		return RenderTmpl(c, tmplCreateCustomer, tmplData)
	}

	// update user
	if appErr := service.UserUpdate(ctx.User, tmplData.Request.Sanitized()); appErr != nil {
		if appErr.IsValidation() {
			tmplData.Error = appErr

			return RenderTmpl(c, tmplUserSettings, tmplData)
		}
	}

	// reset template data for the changes to take effect
	tmplData.tmplData.User = *ctx.User
	tmplData.tmplData.NavbarData.User = *ctx.User
	tmplData.tmplData.SidebarData.User = *ctx.User

	// render template
	return RenderTmpl(c, tmplUserSettings, tmplData)
}
