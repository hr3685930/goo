package handler

import (
	"goo/internal/errors"
	"bytes"
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"io/ioutil"
	"strings"
)

func Route(e *echo.Echo, mux *runtime.ServeMux) {
	if !viper.GetBool("app.debug") {
		e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
			Skipper:           middleware.DefaultSkipper,
			StackSize:         4 << 10, // 4 KB
			DisableStackAll:   true,
			DisablePrintStack: true,
		}))
	}
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, "api")
	})

	e.POST("/auth/token", MuxServer(mux))
	e.GET("/test", func(i echo.Context) error {
		return errors.ResourceNotFound("错误")
	})

	g := e.Group("/api")
	g.GET("/me/profile", MuxServer(mux))
	g.PUT("/me/profile", MuxServer(mux))
}

func MuxServer(mux *runtime.ServeMux) func(c echo.Context) error {
	return func(c echo.Context) error {
		req := c.Request()
		cType := req.Header.Get(echo.HeaderContentType)
		if !strings.HasPrefix(cType, echo.MIMEMultipartForm) {
			bodyBytes, _ := ioutil.ReadAll(c.Request().Body)
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			ctx := context.WithValue(c.Request().Context(), "request_body", string(bodyBytes))
			req = c.Request().WithContext(ctx)
		}
		mux.ServeHTTP(c.Response(), req)

		return nil
	}
}
