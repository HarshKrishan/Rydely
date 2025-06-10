package server

import (
	"net/http"
	"user/repository"
	"user/utils"

	"github.com/labstack/echo/v4"
)

var UserRepository = new(repository.UserRepository)

func Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			tokenString, _ := utils.GetTokenFromHeaderOrCookie(c)

			token, err := utils.ParseJWT(tokenString)

			// log.Println("token: ", token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			userID, err := utils.GetUserIDFromToken(token)

			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}
			user, err := UserRepository.GetUserByID(userID)

			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}

			c.Set("userID", user.ID)
			return next(c)
		}
	}
}
