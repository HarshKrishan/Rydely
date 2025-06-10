package server

import (
	"ride/controllers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter() *echo.Echo {
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			if c.Request().RequestURI == "/health" {
				return true
			}
			return false
		},
	}))
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Use(middleware.RequestID())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	// e.Use(Auth())

	health := new(controllers.HealthController)

	e.GET("/health", health.Status())

	rideCtrl := controllers.RideController{}

	rideGroup := e.Group("/ride")

	{
		rideGroup.POST("/create", rideCtrl.CreateRide(), UserAuth())
		rideGroup.PUT("/accept/:id", rideCtrl.AcceptRide(), CaptainAuth())
	}

	return e

}
