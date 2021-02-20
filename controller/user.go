package controller

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/dmitriivoitovich/heameelega/controller/request"
	"github.com/dmitriivoitovich/heameelega/dao"
	"github.com/dmitriivoitovich/heameelega/dao/db"
	"github.com/dmitriivoitovich/heameelega/helper"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type registerUserTmplData struct {
	Request       request.RegisterUser
	InvalidFields []string
	Error         string
}

type loginUserTmplData struct {
	Request       request.LoginUser
	InvalidFields []string
	Error         string
}

var (
	loginUserTmpl = template.Must(
		template.New("login.gohtml").
			Funcs(funcMap).
			ParseFiles(
				"template/user/login.gohtml",
				"template/layout/header.gohtml",
				"template/layout/footer.gohtml",
			),
	)
	registerUserTmpl = template.Must(
		template.New("register.gohtml").
			Funcs(funcMap).
			ParseFiles(
				"template/user/register.gohtml",
				"template/layout/header.gohtml",
				"template/layout/footer.gohtml",
			),
	)
)

func GetLogin(c echo.Context) error {
	return loginUserTmpl.Execute(c.Response(), nil)
}

func PostLogin(c echo.Context) error {
	tmplData := &loginUserTmplData{}

	if err := c.Bind(&tmplData.Request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to bind request")
	}

	tmplData.InvalidFields = tmplData.Request.Validate()

	if len(tmplData.InvalidFields) > 0 {
		return loginUserTmpl.Execute(c.Response(), tmplData)
	}

	user, err := dao.UserByEmail(tmplData.Request.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tmplData.Error = "Credentials invalid"

			return loginUserTmpl.Execute(c.Response(), tmplData)
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tmplData.Request.Password)); err != nil {
		tmplData.Error = "Credentials invalid"

		return loginUserTmpl.Execute(c.Response(), tmplData)
	}

	return c.Redirect(http.StatusSeeOther, helper.PageURLDashboard())
}

func GetRegister(c echo.Context) error {
	return registerUserTmpl.Execute(c.Response(), nil)
}

func PostRegister(c echo.Context) error {
	tmplData := &registerUserTmplData{}

	if err := c.Bind(&tmplData.Request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to bind request")
	}

	tmplData.InvalidFields = tmplData.Request.Validate()

	if len(tmplData.InvalidFields) > 0 {
		return registerUserTmpl.Execute(c.Response(), tmplData)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(tmplData.Request.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate password hash")
	}

	user := &db.User{
		ID:       uuid.New(),
		Email:    tmplData.Request.Email,
		Password: string(hash),
	}

	if err := dao.CreateUser(user); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			tmplData.InvalidFields = append(tmplData.InvalidFields, "email")
			tmplData.Error = "Email taken"

			return registerUserTmpl.Execute(c.Response(), tmplData)
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create new user")
	}

	return c.Redirect(http.StatusSeeOther, helper.PageURLDashboard())
}
