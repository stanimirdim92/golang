package routes

import (
	"os"
	"postcodes/controllers"
	"postcodes/middlewares"
	"strconv"

	// 	"strconv"

	"github.com/labstack/echo/v4"
	echoMW "github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"
)

func LoadRoutes() *echo.Echo {
	router := echo.New()
	router.Debug = true
	router.AutoTLSManager.HostPolicy = autocert.HostWhitelist(os.Getenv("HTTP_HOST"))
	router.Debug, _ = strconv.ParseBool(os.Getenv("APP_DEBUG"))

	if s, _ := strconv.ParseBool(os.Getenv("USE_SSL")); s == false {
		router.Pre(echoMW.HTTPSRedirect())
		router.Use(echoMW.Secure())
	}

	router.Pre(echoMW.AddTrailingSlash())
	router.Use(middlewares.CheckIP)
	router.Use(echoMW.Logger())
	router.Use(echoMW.Recover())
	router.Use(echoMW.CORS())
	router.Use(echoMW.Gzip())

	v1 := router.Group("/api/v1/postcodes")
	{
		v1.GET("/", controllers.Index)
		v1.GET("/:postcode", controllers.FindPostcode)
	}

	return router
}
