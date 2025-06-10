package server

import (
	"log"
	"net/http"
	"os"
	"ride/utils"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
)

var restClient = resty.New()

func UserAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			tokenString, err := utils.GetTokenFromHeaderOrCookie(c)
			if err != nil || tokenString == "" {
				log.Println("err: ", err)
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing or invalid token"})
			}

			token, err := utils.ParseJWT(tokenString)
			if err != nil || token == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			url := os.Getenv("GET_USER_PROFILE")

			response := make(map[string]interface{})

			resp, _ := restClient.R().
				SetHeader("Authorization", "Bearer "+tokenString).
				SetResult(&response).
				Get(url)
			log.Println("response: ", response)
			if resp.StatusCode() == 200 {
				c.Set("userID", response["id"])
				return next(c)
			}

			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized!"})
		}
	}
}


func CaptainAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			tokenString, err := utils.GetTokenFromHeaderOrCookie(c)
			if err != nil || tokenString == "" {
				log.Println("err: ", err)
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing or invalid token"})
			}

			token, err := utils.ParseJWT(tokenString)
			if err != nil || token == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			url := os.Getenv("GET_CAPTAIN_PROFILE")

			response := make(map[string]interface{})

			resp, _ := restClient.R().
				SetHeader("Authorization", "Bearer "+tokenString).
				SetResult(&response).
				Get(url)
			log.Println("response: ", response)
			if resp.StatusCode() == 200 {
				c.Set("captainID", response["id"])
				return next(c)
			}

			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized!"})
		}
	}
}
