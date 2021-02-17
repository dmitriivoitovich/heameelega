package main

import (
	"context"
	"github.com/dmitriivoitovich/wallester-test-assignment/config"
	"github.com/dmitriivoitovich/wallester-test-assignment/controller"
	"github.com/dmitriivoitovich/wallester-test-assignment/dao/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"os/signal"
	"time"
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
	// echo instance
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = httpErrorHandler

	// DB connection
	dbConf := config.DBConfig()
	go db.InitConn(dbConf, e.Logger)

	// middlewares
	e.Use(middleware.Logger())

	// routes
	e.Static("/public", "public")
	e.GET("/", controller.GetHome)
	e.GET("/dashboard", controller.GetDashboard)
	e.GET("/customers", controller.GetSearchCustomers)
	e.GET("/customers/new", controller.GetCreateCustomer)
	e.POST("/customers/new", controller.PostCreateCustomer)
	e.GET("/customers/:id", controller.GetViewCustomer)
	//e.GET("/customers/:id/edit", controller.GetEditCustomer)
	//e.POST("/customers/:id/edit", controller.PostEditCustomer)

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
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	t := controller.Error500Tmpl
	if code == http.StatusNotFound {
		t = controller.Error400Tmpl
	}

	if err := t.Execute(c.Response(), nil); err != nil {
		c.Logger().Error(err)
	}

	c.Logger().Error(err)
}
