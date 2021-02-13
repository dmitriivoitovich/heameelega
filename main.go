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

	// DB connection
	dbConf := config.DBConfig()
	go db.InitConn(dbConf, e.Logger)

	// middlewares
	e.Use(middleware.Logger())

	// routes
	//e.GET("/", controller.GetHome)
	//e.GET("/customers", controller.GetListCustomers)
	e.GET("/customers/new", controller.GetCreateCustomer)
	e.POST("/customers/new", controller.PostCreateCustomer)
	//e.GET("/customers/:id", controller.GetShowCustomer)
	//e.GET("/customers/:id/edit", controller.GetEditCustomer)
	//e.POST("/customers/:id/edit", controller.PostEditCustomer)

	// http server instance
	s := &http.Server{
		Addr:           ":" + httpPort,
		ReadTimeout:    httpReadTimeout,
		WriteTimeout:   httpWriteTimeout,
		IdleTimeout:    httpIdleTimeout,
		MaxHeaderBytes: maxHeaderSize,
	}

	go func() {
		if err := e.StartServer(s); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), httpServerShutdownTimeout)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
