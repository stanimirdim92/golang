package middlewares

import (
	"github.com/labstack/echo/v4"
	"net"
	"net/http"
)

// TODO validate and test nets
func CheckIP(next echo.HandlerFunc) echo.HandlerFunc {
	// The URIs that should have access to APP
	clientIPs := []string{}

	return func(c echo.Context) error {
		clientIP := net.ParseIP(c.RealIP()).String()
		for _, clientip := range clientIPs {
			ip := net.ParseIP(clientip)
			if ip.String() == clientIP {
				return next(c)
			}
		}

		return echo.NewHTTPError(http.StatusUnauthorized, "Not Allowed"+clientIP)
	}
}
