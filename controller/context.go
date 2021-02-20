package controller

import (
	"errors"
	"net/http"

	"github.com/dmitriivoitovich/heameelega/dao"
	"github.com/dmitriivoitovich/heameelega/dao/db"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const (
	sessionCookieName = "session"
)

type AppContext struct {
	echo.Context

	User *db.User
}

func GenerateAppContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := AppContext{Context: c}

		cookie, err := c.Cookie(sessionCookieName)
		if err != nil {
			return next(ctx)
		}

		tokenUUID, err := uuid.Parse(cookie.Value)
		if err != nil {
			return next(ctx)
		}

		if cookie.Value != "" {
			user, err := dao.UserByAccessToken(tokenUUID)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to load user")
			}

			ctx.User = user
		}

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
