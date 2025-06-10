package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthController struct{}

func (h *HealthController) Status() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		return c.String(http.StatusOK, "Working!")
	}

}
