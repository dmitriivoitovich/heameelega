package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dmitriivoitovich/heameelega/config"
	"github.com/dmitriivoitovich/heameelega/controller"
	"github.com/dmitriivoitovich/heameelega/dao/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	httpPort  = "80"
	httpsPort = "443"

	httpServerShutdownTimeout = time.Second * 30

	httpReadTimeout  = time.Second * 10
	httpWriteTimeout = time.Second * 10
	httpIdleTimeout  = time.Second * 5
	maxHeaderSize    = 1024 * 4
)

func main() {
	// set app timezone
	time.Local = time.UTC

	// echo instance
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = httpErrorHandler

	// read configuration file
	if err := config.Read(); err != nil {
		e.Logger.Fatal(err, "failed to read config file")
	}

	// DB connection
	dbConf := config.DBConfig()
	go db.InitConn(dbConf, e.Logger)

	// middlewares
	e.Use(middleware.Logger())
	e.Use(controller.GenerateAppContext)

	// routes without authorisation
	e.Static("/public", "public")
	e.GET("/", controller.GetHome)
	e.GET("/login", controller.GetLogin)
	e.POST("/login", controller.PostLogin)
	e.GET("/register", controller.GetRegister)
	e.POST("/register", controller.PostRegister)

	// routes require authorisation
	e.POST("/logout", controller.PostLogout, controller.CheckAuth)
	e.GET("/settings", controller.GetSettings, controller.CheckAuth)
	e.POST("/settings", controller.PostSettings, controller.CheckAuth)
	e.GET("/dashboard", controller.GetDashboard, controller.CheckAuth)
	e.GET("/customers", controller.GetSearchCustomers, controller.CheckAuth)
	e.GET("/customers/new", controller.GetCreateCustomer, controller.CheckAuth)
	e.GET("/customers/fake", controller.GetGenerateCustomers, controller.CheckAuth)
	e.POST("/customers/new", controller.PostCreateCustomer, controller.CheckAuth)
	e.GET("/customers/:id", controller.GetViewCustomer, controller.CheckAuth)
	e.GET("/customers/:id/edit", controller.GetEditCustomer, controller.CheckAuth)
	e.POST("/customers/:id/edit", controller.PostEditCustomer, controller.CheckAuth)
	e.POST("/customers/:id/delete", controller.DeleteCustomer, controller.CheckAuth)

	// http server instance
	startWebServer(e)

	// wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), httpServerShutdownTimeout)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func startWebServer(e *echo.Echo) {
	e.Server = &http.Server{
		ReadTimeout:    httpReadTimeout,
		WriteTimeout:   httpWriteTimeout,
		IdleTimeout:    httpIdleTimeout,
		MaxHeaderBytes: maxHeaderSize,
	}

	e.TLSServer = &http.Server{
		ReadTimeout:    httpReadTimeout,
		WriteTimeout:   httpWriteTimeout,
		IdleTimeout:    httpIdleTimeout,
		MaxHeaderBytes: maxHeaderSize,
	}

	go func() {
		if err := e.Start(":" + httpPort); err != nil {
			e.Logger.Info("shutting down the http server")
		}
	}()

	go func() {
		tlsConfig := config.AppTLS()

		if tlsConfig.Enabled {
			if err := e.StartTLS(":"+httpsPort, tlsConfig.Cert, tlsConfig.Key); err != nil {
				e.Logger.Info("shutting down the https server")
			}
		}
	}()
}

func httpErrorHandler(err error, c echo.Context) {
	var httpErr *echo.HTTPError
	if !errors.As(err, &httpErr) {
		httpErr = echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("unexpected error type: %v", err))
	}

	var tmplName string

	switch httpErr.Code {
	case http.StatusBadRequest, http.StatusMethodNotAllowed:
		tmplName = controller.TmplError400
	case http.StatusUnauthorized:
		tmplName = controller.TmplError401
	case http.StatusNotFound:
		tmplName = controller.TmplError404
	default:
		tmplName = controller.TmplError500
	}

	c.Response().WriteHeader(httpErr.Code)

	if err := controller.RenderTmpl(c, tmplName, nil); err != nil {
		c.Logger().Error(err)
	}
}
