package server

import (
	"captain/controllers"

	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
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

	captainCtrl := controllers.CaptainController{}

	captainGroup := e.Group("/captain")

	{
		captainGroup.POST("/register", captainCtrl.Register())
		captainGroup.POST("/login", captainCtrl.Login())
		captainGroup.GET("/logout", captainCtrl.Logout())
		captainGroup.GET("/profile", captainCtrl.GetCaptainProfile(), Auth())
		captainGroup.PUT("/toggleActive/:id", captainCtrl.ToggleActive(), Auth())
	}

	return e

}
