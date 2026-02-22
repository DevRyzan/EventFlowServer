package middleware

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

// RequestLogger logs HTTP requests with method, path, status, duration, and client IP.
func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			path := req.URL.Path
			method := req.Method
			clientIP := c.RealIP()

			err := next(c)

			status := c.Response().Status
			duration := time.Since(start)

			level := slog.LevelInfo
			if status >= 500 {
				level = slog.LevelError
			} else if status >= 400 {
				level = slog.LevelWarn
			}

			slog.LogAttrs(c.Request().Context(), level, "request",
				slog.String("method", method),
				slog.String("path", path),
				slog.Int("status", status),
				slog.Duration("duration", duration),
				slog.String("client_ip", clientIP),
			)

			return err
		}
	}
}
