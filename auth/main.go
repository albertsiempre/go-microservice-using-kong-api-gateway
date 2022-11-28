package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// @title switcher API
// @version 1.0
// @description This is an API documentation Partner Only.
func main() {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		e.DefaultHTTPErrorHandler(err, c)
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	// Handler for putting app request and response timestamp. This used for get elapsed time
	e.Use(ServiceRequestTime)

	e.GET("/auth", func(c echo.Context) error {
		jmH, _ := json.Marshal(c.Request().Header)
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, jmH, "", "\t"); err != nil {
			fmt.Printf("JSON parse error: %v", err)
			return c.String(http.StatusInternalServerError, "")
		}

		strHeader := prettyJSON.String()
		fmt.Println(strHeader)

		return c.String(http.StatusOK, strHeader)
	})

	// Start server
	go func() {
		if err := e.Start(":8888"); err != nil {
			e.Logger.Info("Shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}

// ServiceRequestTime middleware adds a `Server` header to the response.
func ServiceRequestTime(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Request().Header.Set("X-Service", "auth-external")
		c.Request().Header.Set("X-Auth-RequestTime", time.Now().Format(time.RFC3339))
		return next(c)
	}
}
