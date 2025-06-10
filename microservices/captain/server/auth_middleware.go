package server

import (
	"captain/repository"
	"captain/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

var CaptainRepository = new(repository.CaptainRepository)

func Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			tokenString, _ := utils.GetTokenFromHeaderOrCookie(c)

			token, err := utils.ParseJWT(tokenString)

			// log.Println("token: ", token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			captainID, err := utils.GetUserIDFromToken(token)

			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}
			captain, err := CaptainRepository.GetCaptainByID(captainID)

			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}

			c.Set("captainID", captain.ID)
			return next(c)
		}
	}
}
