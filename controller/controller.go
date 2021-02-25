package controller

import (
	"bytes"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/dmitriivoitovich/heameelega/dao/db"
	"github.com/dmitriivoitovich/heameelega/helper"
	"github.com/dmitriivoitovich/heameelega/service"
	"github.com/dmitriivoitovich/heameelega/util/apperror"
	"github.com/dmitriivoitovich/heameelega/util/i18n"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	TmplError400        = "error400"
	TmplError401        = "error401"
	TmplError404        = "error404"
	TmplError500        = "error500"
	tmplLoginUser       = "loginUser"
	tmplRegisterUser    = "registerUser"
	tmplDashboard       = "dashboard"
	tmplCreateCustomer  = "createCustomer"
	tmplViewCustomer    = "viewCustomer"
	tmplEditCustomer    = "editCustomer"
	tmplSearchCustomers = "searchCustomers"

	sessionCookieName = "sess"
	sessionCookieTTL  = time.Hour * 24 * 31 * 12
)

type tmplData struct {
	User        db.User
	NavbarData  NavbarData
	SidebarData SidebarData
}

type NavbarData struct {
	User   db.User
	Search string
}

type SidebarData struct {
	User       db.User
	ActivePage string
}

var (
	funcMap = template.FuncMap{
		"linkHome":            helper.PageURLHome,
		"linkLogin":           helper.PageURLLogin,
		"linkRegister":        helper.PageURLRegister,
		"linkLogout":          helper.PageURLLogout,
		"linkDashboard":       helper.PageURLDashboard,
		"linkCustomers":       helper.PageURLCustomers,
		"linkSearchCustomers": helper.PageURLSearchCustomers,
		"linkNewCustomer":     helper.PageURLNewCustomer,
		"linkViewCustomer":    helper.PageURLViewCustomer,
		"linkEditCustomer":    helper.PageURLEditCustomer,
		"inSlice": func(slice []string, key string) bool {
			for i := range slice {
				if slice[i] == key {
					return true
				}
			}

			return false
		},
		"seq": func(items uint32) []uint32 {
			res := make([]uint32, 0, items)

			for i := uint32(0); i < items; i++ {
				res = append(res, i)
			}

			return res
		},
		"inc": func(item uint32) uint32 {
			return item + 1
		},
		"dec": func(item uint32) uint32 {
			return item - 1
		},
		"i18n": i18n.Translate,
	}

	layout = []string{
		"template/layout/header.gohtml",
		"template/layout/footer.gohtml",
		"template/layout/navbar.gohtml",
		"template/layout/sidebar.gohtml",
	}

	templates = map[string]*template.Template{
		TmplError400:        tmplWithLayout("template/error/400.gohtml"),
		TmplError401:        tmplWithLayout("template/error/401.gohtml"),
		TmplError404:        tmplWithLayout("template/error/404.gohtml"),
		TmplError500:        tmplWithLayout("template/error/500.gohtml"),
		tmplLoginUser:       tmplWithLayout("template/user/login.gohtml"),
		tmplRegisterUser:    tmplWithLayout("template/user/register.gohtml"),
		tmplDashboard:       tmplWithLayout("template/dashboard/dashboard.gohtml"),
		tmplCreateCustomer:  tmplWithLayout("template/customer/create.gohtml"),
		tmplViewCustomer:    tmplWithLayout("template/customer/view.gohtml"),
		tmplEditCustomer:    tmplWithLayout("template/customer/edit.gohtml"),
		tmplSearchCustomers: tmplWithLayout("template/customer/search.gohtml"),
	}
)

type AppContext struct {
	echo.Context

	User *db.User
}

func GenerateAppContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := AppContext{Context: c}

		cookie := getSessionCookie(c)

		user, appErr := service.UserAuth(cookie)
		if appErr != nil {
			return httpError(*appErr)
		}

		ctx.User = user

		return next(ctx)
	}
}

func CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(AppContext)

		if ctx.User != nil {
			return next(ctx)
		}

		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized to perform request")
	}
}

func httpError(e apperror.Error) *echo.HTTPError {
	return echo.NewHTTPError(e.HTTPCode, e.Message())
}

func redirect(c echo.Context, url string) error {
	if err := c.Redirect(http.StatusSeeOther, url); err != nil {
		return httpError(*apperror.Internal(err, "failed to perform redirect"))
	}

	return nil
}

func getSessionCookie(c echo.Context) string {
	cookie, err := c.Cookie(sessionCookieName)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(cookie.Value)
}

func setSessionCookie(c echo.Context, token uuid.UUID) {
	cookie := &http.Cookie{
		Name:    sessionCookieName,
		Value:   token.String(),
		Expires: time.Now().Add(sessionCookieTTL),
	}

	c.SetCookie(cookie)
}

func RenderTmpl(c echo.Context, name string, data interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		return httpError(*apperror.Internal(nil, "no template for key "+name))
	}

	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	if err := tmpl.Execute(buf, data); err != nil {
		return httpError(*apperror.Internal(err, "failed to render template"))
	}

	if _, err := c.Response().Write(buf.Bytes()); err != nil {
		return httpError(*apperror.Internal(err, "failed to write response body"))
	}

	return nil
}

func tmplWithLayout(name string) *template.Template {
	return template.Must(
		template.New(
			filepath.Base(name)).
			Funcs(funcMap).
			ParseFiles(append(layout, name)...),
	)
}
