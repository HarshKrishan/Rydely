package controllers

import (
	"captain/models"
	"captain/repository"
	"captain/utils"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

type CaptainController struct{}

var CaptainRepository = new(repository.CaptainRepository)

func (uc *CaptainController) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Registration logic goes here
		captain := models.Captain{}

		err := c.Bind(&captain)

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}

		if captain.Name == "" || captain.Email == "" || captain.Password == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "All fields are required"})
		}

		// check if captain.Email already exists in the database
		captainExists, err := CaptainRepository.CheckCaptainExists(captain.Email)

		if err != nil {
			log.Infoln("Error checking captain existence:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}

		if captainExists {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Captain already exists"})
		}

		hashedPassword, err := utils.HashPassword(captain.Password)
		if err != nil {
			log.Infoln("Error hashing password:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}
		captain.Password = hashedPassword

		err = CaptainRepository.CreateCaptain(captain)
		if err != nil {
			log.Infoln("Error creating captain:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}

		return c.String(http.StatusOK, "Captain registered successfully!")
	}
}
func (uc *CaptainController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {

		request := models.LoginRequest{}

		err := c.Bind(&request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}

		if request.Email == "" || request.Password == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email and password are required"})
		}

		// compare the email and password with the database
		captainExists, err := CaptainRepository.CheckCaptainExists(request.Email)
		if err != nil {
			log.Infoln("Error checking captain existence:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}
		if !captainExists {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
		}

		captain, err := CaptainRepository.GetCaptainByEmail(request.Email)

		if err != nil {
			log.Infoln("Error retrieving captain:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}

		// Verify the password
		isValidPassword := utils.CheckPasswordHash(request.Password, captain.Password)
		if !isValidPassword {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
		}

		token, err := utils.GenerateJWT(captain.ID)
		if err != nil {
			log.Infoln("Error generating JWT:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}

		// Set the token in the response header and cookie

		c.Response().Header().Set("Authorization", "Bearer "+token)
		c.SetCookie(&http.Cookie{
			Name:  "token",
			Value: token,
		})

		return c.String(http.StatusOK, "Login successful!")
	}
}

func (uc *CaptainController) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Clear the token from the response header and cookie
		c.Response().Header().Del("Authorization")
		c.SetCookie(&http.Cookie{
			Name:   "token",
			Value:  "",
			MaxAge: -1, // Set MaxAge to -1 to delete the cookie
		})
		return c.String(http.StatusOK, "Logout successful!")
	}
}

func (uc *CaptainController) ToggleActive() echo.HandlerFunc {
	return func(c echo.Context) error {
		captainID := c.Param("id")
		if captainID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Captain ID is required"})
		}

		captain, err := CaptainRepository.GetCaptainByID(captainID)
		if err != nil {
			log.Infoln("Error retrieving captain:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}
		log.Println("Captain retrieved:", captain)
		captain.IsActive = !captain.IsActive // Toggle the IsActive status

		err = CaptainRepository.UpdateCaptain(captain)
		if err != nil {
			log.Infoln("Error updating captain:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}

		status := "deactivated"
		if captain.IsActive {
			status = "activated"
		}

		return c.String(http.StatusOK, "Captain "+status+" successfully!")
	}
}

func (uc *CaptainController) GetCaptainProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the user ID from the JWT token

		tokenString, _ := utils.GetTokenFromHeaderOrCookie(c)

		token, err := utils.ParseJWT(tokenString)
		captainID, err := utils.GetUserIDFromToken(token)
		if err != nil {
			log.Infoln("Error getting user ID from token:", err)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		log.Infoln("Captain ID from token:", captainID)
		// Retrieve the user profile from the database
		user, err := CaptainRepository.GetCaptainByID(captainID)
		if err != nil {
			log.Infoln("Error retrieving user profile:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}

		return c.JSON(http.StatusOK, user)
	}
}
