package controllers

import (
	"net/http"
	"user/models"
	"user/repository"
	"user/utils"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

type UserController struct{}

var UserRepository = new(repository.UserRepository)

func (uc *UserController) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Registration logic goes here
		user := models.User{}

		err := c.Bind(&user)

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}

		if user.Name == "" || user.Email == "" || user.Password == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "All fields are required"})
		}

		// check if user.Email already exists in the database
		userExists, err := UserRepository.CheckUserExists(user.Email)

		if err != nil {
			log.Infoln("Error checking user existence:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}

		if userExists {
			return c.JSON(http.StatusConflict, map[string]string{"error": "User already exists"})
		}

		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			log.Infoln("Error hashing password:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}
		user.Password = hashedPassword

		err = UserRepository.CreateUser(user)
		if err != nil {
			log.Infoln("Error creating user:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}

		return c.String(http.StatusOK, "User registered successfully!")
	}
}
func (uc *UserController) Login() echo.HandlerFunc {
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
		userExists, err := UserRepository.CheckUserExists(request.Email)
		if err != nil {
			log.Infoln("Error checking user existence:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}
		if !userExists {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
		}

		user, err := UserRepository.GetUserByEmail(request.Email)

		if err != nil {
			log.Infoln("Error retrieving user:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}

		// Verify the password
		isValidPassword := utils.CheckPasswordHash(request.Password, user.Password)
		if !isValidPassword {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
		}

		token, err := utils.GenerateJWT(user.ID)
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

func (uc *UserController) GetUserProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the user ID from the JWT token

		tokenString, _ := utils.GetTokenFromHeaderOrCookie(c)

		token, err := utils.ParseJWT(tokenString)
		userID, err := utils.GetUserIDFromToken(token)
		if err != nil {
			log.Infoln("Error getting user ID from token:", err)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		log.Infoln("User ID from token:", userID)
		// Retrieve the user profile from the database
		user, err := UserRepository.GetUserByID(userID)
		if err != nil {
			log.Infoln("Error retrieving user profile:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}

		return c.JSON(http.StatusOK, user)
	}
}

func (uc *UserController) Logout() echo.HandlerFunc {
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
