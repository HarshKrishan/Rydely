package controllers

import (
	"net/http"
	"ride/models"
	"ride/rabbitmq"
	"ride/repository"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

type RideController struct{}

var RideRepository = new(repository.RideRepository)

func (rc RideController) CreateRide() echo.HandlerFunc {
	return func(c echo.Context) error {
		var rideRequest models.RideRequest
		if err := c.Bind(&rideRequest); err != nil {
			log.Errorln("Error binding ride:", err)
			return c.JSON(http.StatusBadRequest, models.RideResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid request data",
			})
		}

		ride := models.Ride{
			UserID:      c.Get("userID").(string),
			Pickup:      rideRequest.Pickup,
			Destination: rideRequest.Destination,
			Status:      models.RideStatusRequested,
			CaptainID:   "",
		}
		// Create the ride
		if err := RideRepository.CreateRide(&ride); err != nil {
			log.Errorln("Error creating ride:", err)
			return c.JSON(http.StatusInternalServerError, models.RideResponse{
				Status:  http.StatusInternalServerError,
				Message: "Failed to create ride",
			})
		}

		rabbitmq.PublishRide(ride,"ride_queue")

		return c.JSON(http.StatusCreated, models.RideResponse{
			Status:  http.StatusCreated,
			Message: "Ride created successfully",
		})
	}
}

func (rc RideController) AcceptRide() echo.HandlerFunc {
	return func(c echo.Context) error {

		rideID := c.Param("id")
		if rideID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Ride ID is required"})
		}

		ride, err := RideRepository.GetRideByID(rideID)
		if err != nil {
			log.Infoln("Error retrieving ride:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}
		log.Println("Ride retrieved:", ride)
		ride.Status = models.RideStatusAccepted
		ride.CaptainID = c.Get("captainID").(string)

		err = RideRepository.UpdateRide(ride)
		if err != nil {
			log.Infoln("Error updating ride:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}

		rabbitmq.PublishRide(ride, "accept-ride")
		return c.String(http.StatusOK, "Ride accepted successfully!")

	}
}
