package server

import (
	"user/controllers"

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

	userCtrl := controllers.UserController{}

	userGroup := e.Group("/user")

	{
		userGroup.POST("/register", userCtrl.Register())
		userGroup.POST("/login", userCtrl.Login())
		userGroup.GET("/logout", userCtrl.Logout())
		userGroup.GET("/profile", userCtrl.GetUserProfile(), Auth())
	}

	return e

}
