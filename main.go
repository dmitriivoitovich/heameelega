package main

import (
	"context"
	"errors"
	"html/template"
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
	httpPort = "80"

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

	// routes
	e.Static("/public", "public")
	e.GET("/", controller.GetHome)
	e.GET("/login", controller.GetLogin)
	// e.POST("/login", controller.PostLogin)
	e.GET("/register", controller.GetRegister)
	// e.POST("/register", controller.PostRegister)
	e.GET("/dashboard", controller.GetDashboard)
	e.GET("/customers", controller.GetSearchCustomers)
	e.GET("/customers/new", controller.GetCreateCustomer)
	e.POST("/customers/new", controller.PostCreateCustomer)
	e.GET("/customers/:id", controller.GetViewCustomer)
	e.GET("/customers/:id/edit", controller.GetEditCustomer)
	e.POST("/customers/:id/edit", controller.PostEditCustomer)

	// http server instance
	go startWebServer(e)

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
	s := &http.Server{
		Addr:           ":" + httpPort,
		ReadTimeout:    httpReadTimeout,
		WriteTimeout:   httpWriteTimeout,
		IdleTimeout:    httpIdleTimeout,
		MaxHeaderBytes: maxHeaderSize,
	}

	if err := e.StartServer(s); err != nil {
		e.Logger.Info("shutting down the server")
	}
}

func httpErrorHandler(err error, c echo.Context) {
	var code int

	var httpError *echo.HTTPError
	if errors.As(err, &httpError) {
		code = httpError.Code
	}

	var tmpl *template.Template

	switch code {
	case http.StatusBadRequest, http.StatusMethodNotAllowed:
		tmpl = controller.Error400Tmpl
	case http.StatusNotFound:
		tmpl = controller.Error404Tmpl
	default:
		tmpl = controller.Error500Tmpl
	}

	if err := tmpl.Execute(c.Response(), nil); err != nil {
		c.Logger().Error(err)
	}

	c.Logger().Error(err)
}
