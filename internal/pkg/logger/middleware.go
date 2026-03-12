package logger

import (
	"time"

	"github.com/labstack/echo/v4"
)

func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {

			start := time.Now()

			err := next(c)

			req := c.Request()
			res := c.Response()

			duration := time.Since(start)

			Log.Info("http request",
				"method", req.Method,
				"path", req.URL.Path,
				"status", res.Status,
				"latency", duration.String(),
				"ip", c.RealIP(),
			)

			return err
		}
	}
}